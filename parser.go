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

	for r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType != EOF {
		err := r.parseValue()
		if err != nil {
			return false, err
		}
		if len(r.tokens) < r.curIdx {
			return false, fmt.Errorf("sequence is never finished")
		}
		r.curIdx++
	}
	/* 	if len(r.stack) != 0 {
		return false, fmt.Errorf("braces or brackets are inbalanced")
	} */

	return true, nil

}

func (r *Parser) GetTokens() []string {
	var s []string = make([]string, 0)

	for i := 0; i < len(r.tokens); i++ {
		s = append(s, r.tokens[i].Value)
	}

	return s
}

func (r *Parser) GetStack() []TokenType {
	return r.stack
}

// gotta resolve comma problem
func (r *Parser) parseValue() error {
	cur := r.tokens[r.curIdx]

	switch cur.TokenType {
	case LEFT_BRACE:
		return r.parseObj()
	case LEFT_BRACKET:
		return r.parseArray()
	case STRING, TRUE, FALSE, NULL, NUMBER:
		return nil
	case RIGHT_BRACE:
		if len(r.stack) == 0 || r.stack[len(r.stack)-1] != LEFT_BRACE {
			return fmt.Errorf("unexpected }")
		}
		r.stack = r.stack[:len(r.stack)-1]
		return nil
	case RIGHT_BRACKET:
		if len(r.stack) == 0 || r.stack[len(r.stack)-1] != LEFT_BRACKET {
			return fmt.Errorf("unexpected ]")
		}
		r.stack = r.stack[:len(r.stack)-1]
		return nil
	default:
		return fmt.Errorf("expected value but got: %v", cur)
	}
}

func NewParser(input []byte) (*Parser, error) {

	lexer := NewLexer(bufio.NewReader(bytes.NewReader(input)))

	if lexer.Tokens == nil {
		return nil, fmt.Errorf("empty tokens")
	}

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

func (r *Parser) parseRightBracket() error {
	fmt.Println(r.stack, r.curIdx)
	if r.last() != LEFT_BRACKET {
		return fmt.Errorf("incorrect json structure (right bracket)")
	}
	r.popStack()
	return nil
}

func (r *Parser) parseRightBrace() error {
	fmt.Println("parse right brace")
	if r.last() != LEFT_BRACE {
		return fmt.Errorf("incorrect json structure (right brace)")
	}
	r.popStack()

	return nil
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

func (r *Parser) parseComma() error {
	if r.curIdx+1 >= len(r.tokens) || r.curIdx-1 < 0 {
		return fmt.Errorf("incorrect json structure (comma)")
	}
	//next := r.tokens[r.curIdx+1].TokenType
	//cur := r.tokens[r.curIdx]
	//prev := r.tokens[r.curIdx-1].TokenType

	//in-between cases
	/* 	a := prev == RIGHT_BRACKET && next == LEFT_BRACKET
	   	o := prev == RIGHT_BRACE && next == LEFT_BRACE

	   	b := betweenValues(prev, next)

	   	if !a && !o && !b {
	   		return fmt.Errorf("incorrect json structure (comma)")
	   	} */

	return nil

}

func betweenValues(l, r TokenType) bool {
	lt := l == TRUE || l == FALSE || l == NULL || l == NUMBER || l == STRING

	rt := r == TRUE || r == FALSE || r == NULL || r == NUMBER || r == STRING

	return lt && rt
}

// { key : value, key1 : value1, obj : [ { val : arr1 }, { val2: arr2 } ] }
//
//
//
//
//

func (r *Parser) parseArray() error {
	// Move past the left bracket
	r.curIdx++

	r.stack = append(r.stack, LEFT_BRACKET)

	// Check if the array is empty
	if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
		return r.parseRightBracket()
	}

	for r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType != RIGHT_BRACKET {

		if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
			break
		}

		// Parse the value (this will handle nested arrays)
		err := r.parseValue()
		if err != nil {
			return err
		}

		fmt.Println("after", r.curIdx, "curr", r.tokens[r.curIdx])

		r.curIdx++

		if r.curIdx+1 < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType != COMMA {
			return fmt.Errorf("expected comma or closing bracket in array, but got %s %v", r.tokens[r.curIdx].Value, r.stack)
		}

		if r.curIdx+1 < len(r.tokens)-1 && r.tokens[r.curIdx+1].TokenType != RIGHT_BRACKET {
			r.curIdx++
		}
	}

	if r.tokens[r.curIdx].TokenType != RIGHT_BRACKET {
		return fmt.Errorf("error incorrect json structure (right bracket), got %v", r.tokens[r.curIdx])
	}

	if r.last() != LEFT_BRACKET {
		return fmt.Errorf("error incorrect json structure unbalanced brackets")
	}

	if len(r.stack) > 0 {
		r.popStack()
	}

	//fmt.Println("Got out", r.curIdx, r.stack, len(r.stack), r.tokens[r.curIdx])
	return nil
}

func (r *Parser) ParseObj() error {
	return r.parseObj()
}

func (r *Parser) parseObj() error {

	r.stack = append(r.stack, LEFT_BRACE)

	r.curIdx++ //skip opening bracket {

	if r.tokens[r.curIdx].TokenType == RIGHT_BRACE {
		return r.parseRightBrace()
	}

	for r.tokens[r.curIdx].TokenType != RIGHT_BRACE {
		err := r.parseKeyvalue()

		if err != nil {
			return err
		}

		r.curIdx++

		if r.curIdx >= 0 && r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx-1].TokenType != COMMA && r.tokens[r.curIdx].TokenType != RIGHT_BRACE {
			return fmt.Errorf("expected comma, but got %v", r.tokens[r.curIdx])
		}

	}

	if r.tokens[r.curIdx].TokenType != RIGHT_BRACE {
		return fmt.Errorf("error incorrect json structure (right brace)")
	}

	if r.last() != LEFT_BRACE {
		return fmt.Errorf("error incorrect json structure unbalanced braces")
	}

	if len(r.stack) > 0 {
		r.popStack()
	}

	//
	// { key : { nested: val, nested1: val1 } }

	fmt.Println(r.stack, "went after", r.curIdx, len(r.tokens), "cur value:", r.tokens[r.curIdx].Value)

	return nil
}

func (r *Parser) isValue(el TokenType) bool {
	return el == TRUE || el == FALSE || el == NULL || el == NUMBER || el == STRING
}

func (r *Parser) parseKeyvalue() error {

	cur := r.tokens[r.curIdx]

	if cur.TokenType != STRING {
		return fmt.Errorf("incorrect json structure (object) 1, got: %s, prev: %v", cur.Value, r.tokens[r.curIdx-1])
	}

	r.curIdx++

	//then colon
	if r.tokens[r.curIdx].TokenType != COLON {
		return fmt.Errorf("incorrect json structure (object) 2")
	}

	//skip colon
	r.curIdx++

	err := r.parseValue()
	if err != nil {
		return err
	}

	if r.curIdx+1 < len(r.tokens)-1 && r.tokens[r.curIdx+1].TokenType == COMMA {
		r.curIdx++
	}

	fmt.Println(r.curIdx, "cur", r.tokens[r.curIdx].Value)
	//r.curIdx++
	/* 	if r.curIdx < len(r.tokens)-1 && (r.tokens[r.curIdx].TokenType == COMMA || r.tokens[r.curIdx].TokenType == RIGHT_BRACE) {
		if r.tokens[r.curIdx].TokenType == COMMA {
			err := r.parseComma()
			if err != nil {
				return err
			}
		}
		if r.tokens[r.curIdx].TokenType == RIGHT_BRACE && r.last() == LEFT_BRACE {
			r.popStack()
		}
		r.curIdx++
	} */

	//potential recursion in case value in an object

	return nil
}

func (r *Parser) popStack() {
	fmt.Print("before pop", r.stack)
	if len(r.stack) > 0 {
		r.stack = r.stack[:len(r.stack)-1]
	}
	fmt.Print("after pop", r.stack)
}

func (r *Parser) last() TokenType {
	return r.stack[len(r.stack)-1]
}

func (r *Parser) isValidNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
