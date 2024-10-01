package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
)

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	COLON
	LEFT_BRACKET
	RIGHT_BRACKET
	NULL
	FALSE
	TRUE
	STRING
	NUMBER
	COMMA
)

// how can you add row and col
type Lexer struct {
	Tokens []Token
	Reader *bufio.Reader
}

type Token struct {
	TokenType TokenType
	Value     string
}

func NewLexer(rd *bufio.Reader) *Lexer {
	return &Lexer{
		Reader: rd,
		Tokens: []Token{},
	}
}

// let's just assume it's an array of bytes
func (r *Lexer) Tokenize(rd bufio.Reader) ([]Token, error) {
	var cur byte

	var err error

	for err == nil {
		cur, err = rd.ReadByte()

		if err != nil {
			if err == io.EOF {
				return r.Tokens, nil
			}
			fmt.Printf("reading error %v", err)
			return nil, err
		}

		switch cur {
		case '{':
			r.Tokens = append(r.Tokens, Token{TokenType: LEFT_BRACE, Value: "{"})
		case '}':
			r.Tokens = append(r.Tokens, Token{TokenType: RIGHT_BRACE, Value: "}"})
		case '[':
			r.Tokens = append(r.Tokens, Token{TokenType: LEFT_BRACKET, Value: "["})
		case ']':
			r.Tokens = append(r.Tokens, Token{TokenType: RIGHT_BRACKET, Value: "]"})
		case 'f':
			t, err := r.tokenizeBool(&rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case 't':
			t, err := r.tokenizeBool(&rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)

		case '"':
			//handle string
			t, err := r.tokenizeString(&rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)

		case ':':
			r.Tokens = append(r.Tokens, Token{TokenType: COLON, Value: ":"})
		default:

		}
	}

	return r.Tokens, nil
}

func (r *Lexer) tokenizeBool(rd *bufio.Reader) (*Token, error) {
	//look further 5 bytes false
	n, err := rd.Peek(4)
	if err != nil {
		return nil, err
	}
	ok := bytes.Equal(n, []byte("alse"))
	if !ok {
		n, err = rd.Peek(3)
	}
	ok1 := bytes.Equal(n, []byte("rue"))
	if !ok && !ok1 {
		return nil, fmt.Errorf("error parsing bool")
	}
	re := &Token{}
	if ok {
		re.TokenType = FALSE
		re.Value = "false"
	} else {
		re.TokenType = TRUE
		re.Value = "true"
	}

	if ok {
		rd.Discard(4)
	} else {
		rd.Discard(3)
	}

	return re, nil
}

func (r *Lexer) tokenizeString(rd *bufio.Reader) (*Token, error) {

	cur, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}

	var val []byte
	var cur_val byte

	if cur == '\\' || cur == '"' {
		//skip
		cur, err = rd.ReadByte()
	} else {
		val = append(val, cur)
	}

	t := &Token{}

	for cur_val != '"' && err == nil {
		cur_val, err = rd.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error while parsing string token %v", err)
		}

		if cur_val != '"' {
			val = append(val, cur_val)
		}
	}

	t.TokenType = STRING
	t.Value = string(val)

	return t, nil
}

func tokenizeArray() Token {
	return Token{}
}

func tokenizeObj() Token {
	return Token{}
}

func tokenizeNumber() Token {
	return Token{}
}

func tokenizeNull() Token {
	return Token{}
}
