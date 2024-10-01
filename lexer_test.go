package parser

import (
	"bufio"
	"bytes"
	"testing"
)

func TestSimpleJson(t *testing.T) {
	sample := []byte("{\"key\":\"value\"}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize(*rd)
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestBoolJson(t *testing.T) {
	sample := []byte("{\"key\":true}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize(*rd)
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestSimpleErrorJson(t *testing.T) {
	sample := []byte("{\"key\":\"value}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize(*rd)
	if err == nil {
		t.Error("error didn't trigger")
	}

	t.Log(p)
}
