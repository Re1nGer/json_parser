package parser

import (
	"testing"
)

func TestSimpleJsonParser(t *testing.T) {
	sample := []byte("{\"key\": \"key\"}")

	parser, err := NewParser(sample)
	if err != nil {
		t.Errorf("error instantiating parser %v", err)
	}

	_, err = parser.Parse()
	t.Log("tokens", parser.tokens)
	if err != nil {
		t.Errorf("error parsing %v", err)
	}
}

func TestSimplArray(t *testing.T) {
	sample := []byte("[[]]")

	parser, err := NewParser(sample)
	if err != nil {
		t.Errorf("error instantiating parser %v", err)
	}

	t.Log("tokens", parser.tokens)
	_, err = parser.Parse()
	if err != nil {
		t.Errorf("error parsing %v", err)
	}
}
