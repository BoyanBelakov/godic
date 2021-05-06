package data

import (
	"testing"
)

func TestKeysWithPrefix(t *testing.T) {
	a := NewAlphabet("abcdefghijklmnopqrstuvwxyz")
	sm := NewStringMap(a)
	sm.Put("she", 0)
	sm.Put("shells", 1)
	sm.Put("sea", 3)
	sm.Put("sells", 4)
	sm.Put("shore", 5)
	sm.Put("by", 2)
	sm.Put("the", 6)

	want := []string{"she", "shells", "shore"}
	ll := sm.KeysWithPrefix("sh", -1)

	if len(want) != ll.Len() {
		t.Fatalf("want: %d", len(want))
	}

	i := 0
	for e := ll.Front(); e != nil; e = e.Next() {
		v := e.Value.(string)
		if want[i] != v {
			t.Fatalf("want: %s, got: %s", want[i], v)
		}
		i++
	}
}

func TestPut(t *testing.T) {
	a := NewAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ '-АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ")
	sm := NewStringMap(a)
	sm.Put("SHE", "тя")
	sm.Put("SHE", "тя")
	sm.Put("SHELL", "черупка")
	sm.Put("SHELL", "черупка")
	sm.Put("NIL", nil)
	sm.Put("NIL", nil)

	tests := []struct {
		key  string
		want interface{}
	}{
		{key: "SHE", want: "тя"},
		{key: "SHELL", want: "черупка"},
		{key: "NIL", want: nil},
		{key: "SH", want: nil},
		{key: "ELL", want: nil},
		{key: "HEL", want: nil},
		{key: "she", want: nil},
		{key: "черупка", want: nil},
		{key: "", want: nil},
	}

	if s := sm.Size(); s != 3 {
		t.Fatalf("want: 3, got: %d", s)
	}

	for _, test := range tests {
		if got := sm.Get(test.key); got != test.want {
			t.Fatalf("want: %s, got: %s", test.want, got)
		}
	}
}

func TestLongestPrefixOf(t *testing.T) {
	a := NewAlphabet("ABCDEFGHIJKLMNOPQRSTUVWXYZ '-АБВГДЕЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЬЮЯ")
	sm := NewStringMap(a)
	sm.Put("SHE", true)
	sm.Put("SHELL", true)
	sm.Put("ДОК", true)
	sm.Put("ДОКТОР", true)

	tests := []struct {
		pre  string
		want string
	}{
		{pre: "SHE", want: "SHE"},
		{pre: "SHELLS", want: "SHELL"},
		{pre: "ДОК", want: "ДОК"},
		{pre: "ДОКТ", want: "ДОК"},
		{pre: "ДО", want: ""},
		{pre: "ОК", want: ""},
		{pre: "ТОР", want: ""},
		{pre: "", want: ""},
	}

	for _, test := range tests {
		if got := sm.LongestPrefixOf(test.pre); got != test.want {
			t.Fatalf("want: %s, got: %s", test.want, got)
		}
	}
}
