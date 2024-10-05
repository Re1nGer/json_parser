package parser

import (
	"fmt"
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

func TestNestedJson(t *testing.T) {

	sample := []byte("{\"key\": { \"key1\":\"nested\" } }")

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

func TestNestedWithArrayJson(t *testing.T) {

	sample := []byte("{\"key\": \"value\", \"key-n\": 101, \"key-o\": { \"inner key\": \"inner value\", \"inner-key2\": \"inner value\" }, \"arr\":[{\"nested\": \"jee\"}]} ")

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

func TestFailStructure(t *testing.T) {
	sample := []byte("A JSON payload should be an object or array, not a string.")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailStructureArray(t *testing.T) {
	sample := []byte("[\"Unclosed array")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailIllegalInvocation(t *testing.T) {

	sample := []byte("{\"Illegal invocation\": alert()}")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailUnquotedkey(t *testing.T) {

	sample := []byte("{unquoted_key: \"keys must be quoted\"}")

	_, err := NewParser(sample)

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailExtraComma(t *testing.T) {

	sample := []byte("[\"extra comma\",]")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailDoubleExtraComma(t *testing.T) {
	sample := []byte("[\"double extra comma\",,]")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailMissingValue(t *testing.T) {
	sample := []byte("[   , \"<-- missing value\"]")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailCommaAfterClose(t *testing.T) {
	sample := []byte("[\"Comma after the close\"],")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}
