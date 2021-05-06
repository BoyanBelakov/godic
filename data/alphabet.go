package data

import (
	"unicode/utf8"
)

type Alphabet interface {
	Radix() int
	Contains(rune) bool
	ToIndex(rune) int
	ToChar(int) rune
	Valid(string) bool
}

type alphabet struct {
	chars   []rune
	indexes map[rune]int
}

func NewAlphabet(alphabetChars string) Alphabet {
	count := utf8.RuneCountInString(alphabetChars)
	chars := make([]rune, count)
	indexes := make(map[rune]int, count)

	i := 0
	for _, c := range alphabetChars {
		chars[i] = c
		indexes[c] = i
		i++
	}

	return &alphabet{
		chars,
		indexes,
	}
}

func (a *alphabet) Radix() int {
	return len(a.chars)
}

func (a *alphabet) Contains(r rune) bool {
	_, ok := a.indexes[r]
	return ok
}

func (a *alphabet) ToIndex(r rune) int {
	return a.indexes[r]
}

func (a *alphabet) ToChar(i int) rune {
	return a.chars[i]
}

func (a *alphabet) Valid(word string) bool {
	for _, c := range word {
		if !a.Contains(c) {
			return false
		}
	}
	return true
}
