package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

type Parser struct {
	tokens []Token
	curIdx int
	stack  []TokenType
}

func (r *Parser) Parse() (bool, error) {

	for r.tokens[r.curIdx].TokenType != EOF {
		err := r.parseValue()
		if err != nil {
			return false, err
		}
		if len(r.tokens) > r.curIdx {
			return false, fmt.Errorf("sequence is never finished")
		}
		r.curIdx++
	}

	return true, nil
}

// gotta resolve comma problem
func (r *Parser) parseValue() error {
	cur := r.tokens[r.curIdx]

	switch cur.TokenType {
	case LEFT_BRACE:
		return r.parseObj()
	case LEFT_BRACKET:
		return r.parseArray()
	case STRING: //just skip ahead
		return r.parseString()
	case TRUE:
		return r.parseTrue()
	case FALSE:
		return r.parseFalse()
	case NUMBER:
		return r.parseNumber()
	case NULL:
		return r.parseNull()
	default:
		return fmt.Errorf("unknown entity or incorrect structure %s", cur.Value)
	}
}

func NewParser(input []byte) (*Parser, error) {

	lexer := NewLexer(bufio.NewReader(bytes.NewReader(input)))

	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("unable to tokenize %v", err)
	}

	return &Parser{
		tokens: tokens,
		curIdx: 0,
		stack:  make([]TokenType, 0),
	}, nil
}

func (r *Parser) parseString() error {
	if r.tokens[r.curIdx].TokenType != STRING {
		return fmt.Errorf("incorrect json structure (string)")
	}
	return nil
}

func (r *Parser) parseTrue() error {
	if r.tokens[r.curIdx].Value != "true" {
		return fmt.Errorf("incorrect json structure (boolean) true")
	}
	return nil
}

func (r *Parser) parseFalse() error {
	if r.tokens[r.curIdx].Value != "false" {
		return fmt.Errorf("incorrect json structure (boolean) false")
	}
	return nil
}

func (r *Parser) parseNumber() error {
	if !r.isValidNumber(r.tokens[r.curIdx].Value) {
		return fmt.Errorf("incorrect json structure number")
	}
	return nil
}

func (r *Parser) parseNull() error {
	if r.tokens[r.curIdx].Value != "null" {
		return fmt.Errorf("incorrect json structure null")
	}
	return nil
}

func (r *Parser) parseArray() error {
	r.curIdx++ //skip opening bracket
	r.stack = append(r.stack, LEFT_BRACKET)

	cur := r.tokens[r.curIdx]

	switch cur.TokenType {
	case LEFT_BRACKET:
		return r.parseObj()
	case LEFT_BRACE:
		//handle deeply nested structures
		return r.parseArray() //recursion
	default:
		//array of primitives
		comma := false
		for r.tokens[r.curIdx].TokenType != RIGHT_BRACKET {
			//value and comma
			if comma && r.tokens[r.curIdx].TokenType != COMMA {
				return fmt.Errorf("incorrect json structure")
			}
			comma = !comma
			r.curIdx++
			//[1,2]
		}

		for r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
			if r.stack[len(r.stack)-1] != LEFT_BRACKET {
				return fmt.Errorf("incorrect json structure")
			}
			r.popStack()
		}
		return nil
	}
}

func (r *Parser) parseObj() error {

	r.curIdx++ //skip opening bracket
	r.stack = append(r.stack, LEFT_BRACE)

	cur := r.tokens[r.curIdx]

	if cur.TokenType != STRING {
		return fmt.Errorf("incorrect json structure (object) 1")
	}

	//that's where we can consume key part
	//but right now we just want to check integrity
	r.curIdx++

	if r.tokens[r.curIdx].TokenType != COLON {
		return fmt.Errorf("incorrect json structure (object) 2")
	}

	r.curIdx++

	err := r.parseValue()
	if err != nil {
		return err
	}

	r.curIdx++

	if r.tokens[r.curIdx].TokenType != RIGHT_BRACE {
		fmt.Println(r.curIdx, r.tokens, r.stack)
		return fmt.Errorf("incorrect json structure (object) 3")
	}

	if r.stack[len(r.stack)-1] != LEFT_BRACE {
		return fmt.Errorf("incorrect json structure (object) 4")
	}

	r.popStack()

	return nil
}

func (r *Parser) popStack() {
	r.stack = r.stack[:len(r.stack)-1]
}

func (r *Parser) isValidNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
