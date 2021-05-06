package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BoyanBelakov/godic/data"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	splitSymbol = "@"
)

var (
	red     = lipgloss.Color("9")
	grey    = lipgloss.Color("247")
	hotPink = lipgloss.Color("205")

	errStyle  = lipgloss.NewStyle().Foreground(red)
	helpStyle = lipgloss.NewStyle().Foreground(grey)
	loading   = lipgloss.NewStyle().Foreground(hotPink).SetString("Loading...").String()
	header    = lipgloss.NewStyle().Foreground(hotPink).SetString("Godic").String()
)

var dicFilePath string

type loadSuccessMsg struct{}
type trMsg string
type errMsg error

type model struct {
	alphabet   data.Alphabet
	dictionary data.StringMap
	viewport   viewport.Model
	textInput  textinput.Model
	output     string
	loading    bool
	err        error
}

func newModel() *model {
	a := data.NewAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ '-АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ")
	sm := data.NewStringMap(a)

	vp := viewport.Model{
		Width:  78,
		Height: 20,
	}

	ti := textinput.NewModel()
	ti.Placeholder = "Enter a word"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return &model{
		alphabet:   a,
		dictionary: sm,
		viewport:   vp,
		textInput:  ti,
		loading:    true,
	}
}

func load(sm data.StringMap) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Open(dicFilePath)
		if err != nil {
			return errMsg(err)
		}
		defer f.Close()

		var offset int64
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			offset += int64(len(line) + 1)
			i := strings.Index(line, splitSymbol)
			if i != -1 {
				key := line[i+1:]
				sm.Put(key, offset)
			}
		}
		if err := scanner.Err(); err != nil {
			return errMsg(err)
		}

		return loadSuccessMsg{}
	}
}

func tr(offset int64) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Open(dicFilePath)
		if err != nil {
			return errMsg(err)
		}
		defer f.Close()

		_, err = f.Seek(offset, 0)
		if err != nil {
			return errMsg(err)
		}

		sb := &strings.Builder{}
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			i := strings.Index(line, splitSymbol)
			if i != -1 {
				sb.WriteString(line[:i])
				break
			}
			sb.WriteString(line)
			sb.WriteRune('\n')
		}
		if err := scanner.Err(); err != nil {
			return errMsg(err)
		}

		return trMsg(sb.String())
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, load(m.dictionary))
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 3
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			return m, tea.Quit
		case "enter":
			word := m.getInput()
			if !m.alphabet.Valid(word) {
				m.output = "Invalid input"
				return m, nil
			}
			if value := m.dictionary.Get(word); value != nil {
				return m, tr(value.(int64))
			}
			m.output = m.dictionary.LongestPrefixOf(word)
			return m, nil
		case "tab":
			ll := m.dictionary.KeysWithPrefix(m.getInput(), 20)
			var sb strings.Builder
			for e := ll.Front(); e != nil; e = e.Next() {
				sb.WriteString(e.Value.(string))
				sb.WriteRune('\n')
			}
			m.output = sb.String()
			return m, nil
		}
	case loadSuccessMsg:
		m.loading = false
		return m, nil
	case trMsg:
		m.output = string(msg)
		return m, nil
	case errMsg:
		m.err = msg
		return m, nil
	}

	vp, cmd1 := m.viewport.Update(msg)
	m.viewport = vp

	ti, cmd2 := m.textInput.Update(msg)
	m.textInput = ti
	return m, tea.Batch(cmd1, cmd2)
}

func (m *model) View() string {
	if m.err != nil {
		return errStyle.Render(m.err.Error())
	}

	if m.loading {
		return loading
	}

	help := "esc: Quit"
	if m.viewport.Height < strings.Count(m.output, "\n") {
		help += " ↑/↓: Scroll"
	}
	m.viewport.SetContent(m.output)

	return fmt.Sprintf("%s\n%s\n%s\n%s", header, m.textInput.View(), m.viewport.View(), helpStyle.Render(help))
}

func (m *model) getInput() string {
	return strings.ToUpper(m.textInput.Value())
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(ex)
	dicFilePath = filepath.Join(dir, "dic.txt")
}

func main() {
	p := tea.NewProgram(newModel())
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
