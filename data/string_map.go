package data

import "container/list"

type StringMap interface {
	Size() int
	Put(key string, value interface{})
	Get(key string) interface{}
	KeysWithPrefix(pre string, maxResults int) *list.List
	LongestPrefixOf(pre string) string
}

type node struct {
	children []*node
	value    interface{}
	hasValue bool
}

func newNode(radix int) *node {
	return &node{
		children: make([]*node, radix),
	}
}

type trie struct {
	root     *node
	alphabet Alphabet
	size     int
}

func NewStringMap(a Alphabet) StringMap {
	return &trie{
		alphabet: a,
	}
}

func (t *trie) Size() int {
	return t.size
}

func (t *trie) Put(key string, value interface{}) {
	t.root = t.put(t.root, []rune(key), value, 0)
}

func (t *trie) put(node *node, key []rune, value interface{}, depth int) *node {
	if node == nil {
		node = newNode(t.alphabet.Radix())
	}
	if depth == len(key) {
		if !node.hasValue {
			t.size++
		}
		node.value = value
		node.hasValue = true
		return node
	}
	c := key[depth]
	i := t.alphabet.ToIndex(c)
	node.children[i] = t.put(node.children[i], key, value, depth+1)
	return node
}

func (t *trie) Get(key string) interface{} {
	node := t.get(t.root, []rune(key), 0)
	if node != nil {
		return node.value
	}
	return nil
}

func (t *trie) get(node *node, key []rune, depth int) *node {
	if node == nil {
		return nil
	}
	if depth == len(key) {
		return node
	}
	c := key[depth]
	i := t.alphabet.ToIndex(c)
	return t.get(node.children[i], key, depth+1)
}

func (t *trie) KeysWithPrefix(pre string, maxResults int) *list.List {
	results := list.New()
	node := t.get(t.root, []rune(pre), 0)
	t.travel(node, pre, results, maxResults)
	return results
}

func (t *trie) travel(node *node, pre string, results *list.List, maxResults int) {
	if node == nil {
		return
	}
	if results.Len() == maxResults {
		return
	}
	if node.hasValue {
		results.PushBack(pre)
	}

	count := t.alphabet.Radix()
	for i := 0; i < count; i++ {
		c := t.alphabet.ToChar(i)
		t.travel(node.children[i], pre+string(c), results, maxResults)
	}
}

func (t *trie) LongestPrefixOf(pre string) string {
	s := []rune(pre)
	length := t.longestPrefixOf(t.root, s, 0, 0)
	return string(s[:length])
}

func (t *trie) longestPrefixOf(node *node, pre []rune, depth int, length int) int {
	if node == nil {
		return length
	}
	if node.hasValue {
		length = depth
	}
	if len(pre) == depth {
		return length
	}
	c := pre[depth]
	i := t.alphabet.ToIndex(c)
	return t.longestPrefixOf(node.children[i], pre, depth+1, length)
}
