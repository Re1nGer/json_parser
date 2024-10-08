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
		return fmt.Errorf("unexpected token: %v", cur)
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
		return nil
	}

	for {
		// Parse the value (this will handle nested arrays)
		err := r.parseValue()
		if err != nil {
			return err
		}
		if r.curIdx >= len(r.tokens)-1 {
			return nil
		}

		// Check if we've reached the end of the array
		if r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
			return nil
		}

		// Move to the next token
		r.curIdx++
		//fmt.Println("blew up here?", r.curIdx)
		if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == LEFT_BRACE {
			continue
		}

		//fmt.Println("cur value:", r.tokens[r.curIdx].Value, "cur idx", r.curIdx, "next:", r.tokens[r.curIdx:])

		/* 		if r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
			if r.last() == LEFT_BRACKET {
				r.popStack()
			}
			return nil
		} */

		if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == COMMA && r.tokens[r.curIdx+1].TokenType == RIGHT_BRACKET {
			return fmt.Errorf("incorrect json structure: extra comma")
		}

		// If not, the next token should be a comma or not if single element

		if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
			break
		}

		//fmt.Println("cur", r.tokens[r.curIdx], "before", r.tokens[r.curIdx-1], "after", r.tokens[r.curIdx+1])
		if r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType != COMMA {
			return fmt.Errorf("expected comma or closing bracket in array, but got %s %v", r.tokens[r.curIdx].Value, r.stack)
		}

		// Move past the comma
		r.curIdx++
	}

	for r.curIdx < len(r.tokens)-1 && r.tokens[r.curIdx].TokenType == RIGHT_BRACKET {
		if len(r.stack) > 0 && r.last() == LEFT_BRACKET {
			r.popStack()
		} else {
			return fmt.Errorf("incorrect json: unbalanced brackets")
		}
		r.curIdx++
	}
	//fmt.Println("Got out", r.curIdx, r.stack, len(r.stack), r.tokens[r.curIdx])
	return nil
}

func (r *Parser) isValue(el TokenType) bool {
	return el == TRUE || el == FALSE || el == NULL || el == NUMBER || el == STRING
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

	//r.curIdx++ // could be } or ,

	for r.tokens[r.curIdx].TokenType != RIGHT_BRACE {
		err := r.parseKeyvalue()
		if err != nil {
			return err
		}
		r.curIdx++
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

	/* 	r.curIdx++ //either could be } or , or another key

	   	if r.curIdx < len(r.tokens)-2 && r.tokens[r.curIdx].TokenType == COMMA && r.tokens[r.curIdx+1].TokenType == RIGHT_BRACE {
	   		return fmt.Errorf("incorrect json structure: extra comma at obj")
	   	}

	   	if (r.tokens[r.curIdx].TokenType == COMMA && r.tokens[r.curIdx+1].TokenType == STRING) || r.tokens[r.curIdx].TokenType == RIGHT_BRACE {
	   		if r.tokens[r.curIdx].TokenType == RIGHT_BRACE && r.last() == LEFT_BRACE {
	   			r.popStack()
	   		}
	   		r.curIdx++
	   	} */
	//r.curIdx++

	//fmt.Println("cur", r.tokens[r.curIdx], "rest:", r.tokens[r.curIdx:], "next:", r.tokens[r.curIdx+1])
	/* 	if r.tokens[r.curIdx].TokenType == COMMA && r.tokens[r.curIdx+1].TokenType != STRING {
		return fmt.Errorf("incorrect json structure: extra comma at obj")
	} */

	/* 	for r.tokens[r.curIdx].TokenType == STRING {
		err := r.parseKeyvalue()
		if err != nil {
			return err
		}
	} */

	fmt.Println(r.stack, "went after", r.curIdx, len(r.tokens), "cur value:", r.tokens[r.curIdx].Value)
	//need to handle deeply nested structures

	return nil
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
