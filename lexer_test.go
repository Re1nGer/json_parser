package parser

import (
	"bufio"
	"bytes"
	"testing"
)

func TestSimpleJson(t *testing.T) {
	sample := []byte("{\"key\":\"value\"}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestBoolJson(t *testing.T) {
	sample := []byte("{\"key\":true}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestJohnSimpleTest(t *testing.T) {
	sample := []byte("{\"key\": \"value\", \"key2\": \"value\"}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestJohnErrorTest(t *testing.T) {
	sample := []byte("{\"key\": \"value\", key2: \"value\" } ")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err == nil {
		t.Errorf("error didn't trigger")
	}

	t.Log(p)
}

func TestJohnTest2(t *testing.T) {
	sample := []byte("{\"key1\": true, \"key2\": false, \"key3\": null, \"key4\": \"value\"}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestJohnTest3(t *testing.T) {
	sample := []byte("{\"key\": \"value\", \"key-n\": 234, \"key-o\": {}, \"key-l\": [] } ")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}

func TestSimpleErrorJson(t *testing.T) {
	sample := []byte("{\"key\":\"value}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err == nil {
		t.Error("error didn't trigger")
	}

	t.Log(p)
}

func TestJohn4(t *testing.T) {
	sample := []byte("{ \"key\": \"value\", \"key-n\": 101, \"key-o\": { \"inner key\": \"inner value\" }, \"key-l\": [\"list value\"]}")

	rd := bufio.NewReader(bytes.NewReader(sample))

	p, err := NewLexer(rd).Tokenize()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}

	t.Log(p)
}
