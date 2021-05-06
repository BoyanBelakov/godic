package data

import (
	"testing"
)

func TestContains(t *testing.T) {
	a := NewAlphabet("QЯ")
	if !a.Contains('Q') {
		t.FailNow()
	}
	if !a.Contains('Я') {
		t.FailNow()
	}
	if a.Contains('z') {
		t.FailNow()
	}
}

func TestToIndex(t *testing.T) {
	a := NewAlphabet("abcdefghijklmnopqrstuvwxyz")
	if a.ToIndex('a') != 0 {
		t.FailNow()
	}
	if a.ToIndex('z') != 25 {
		t.FailNow()
	}
}

func TestToChar(t *testing.T) {
	a := NewAlphabet("abcdefghijklmnopqrstuvwxyz")
	if a.ToChar(0) != 'a' {
		t.FailNow()
	}
	if a.ToChar(25) != 'z' {
		t.FailNow()
	}

	defer func() {
		if err := recover(); err == nil {
			t.Fatal("expected index out of range")
		}
	}()
	a.ToChar(100)
}

func TestRadix(t *testing.T) {
	a := NewAlphabet("01")
	if a.Radix() != 2 {
		t.FailNow()
	}
}
