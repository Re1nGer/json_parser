package parser

import (
	"bytes"
	"fmt"
	"unicode"
)

type parseStr struct {
	stack        []byte
	specialChars map[byte]bool
}

func (r parseStr) Parse(str []byte, idx int) bool {

	//"{}"; stack = []
	//return len(stack) == 0
	// {"key": "value"}
	// {"": }
	// if c is 't' or 'f': handle bool
	// if c is '"': handle string
	// if isnumeric(c): handle num
	// if c is "{": handle object
	// if c is "[": handle array

	for {
		if idx > len(str) {
			return len(r.stack) != 0
		}

		cur := str[idx]
		fmt.Println("cur", string(cur))
		fmt.Println("stack", string(r.stack))

		if len(r.stack) > 0 && (r.stack[len(r.stack)-1] == cur || (r.stack[len(r.stack)-1] == '{' && cur == '}')) {
			fmt.Println("stack inside", string(r.stack))
			r.stack = r.stack[:len(r.stack)-1]
		}

		_, ok := r.specialChars[cur]
		if ok {
			r.stack = append(r.stack, cur)
		}

		switch cur {
		case '"':
			i := r.handlestr(str, idx+1)
			idx = i + 1
		case '{': //recursively keep parsing
			//i + 1 since we should move the cur
			ok := r.Parse(str, idx+1)
			if !ok {
				return false
			}
		case ':':
			idx++
			continue

		default:
			return false
		}
	}
}

func (r parseStr) handlestr(str []byte, i int) int {
	c := 0

	for j := i; i < len(str); j++ {
		if str[j] == '"' {
			c += 1
		}
		if c == 2 {
			return j
		}
	}
	return i
}

func removeWhitespace(data []byte) []byte {
	return bytes.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, data)
}

func NewParser() *parseStr {
	return &parseStr{
		stack:        []byte{},
		specialChars: map[byte]bool{'"': true, '{': true, '}': true, '[': true, ']': true},
	}
}

//{"key":"v"al,fs" " }
