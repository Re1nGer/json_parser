package parser

import (
	"testing"
)

func TestEmptyObject(t *testing.T) {
	sample := []byte("{}")

	s := NewParser()

	res := s.Parse(sample, 0)

	t.Log(res)

	if !res {
		t.Error("incorrect parsing")
	}
}
