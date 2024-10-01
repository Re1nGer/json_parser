package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
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
	EOF
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
func (r *Lexer) Tokenize() ([]Token, error) {
	var cur byte

	var err error

	rd := r.Reader

	//introduce counter to validate comma positions
	//wouldn't work in nested structures
	//stack is neccessary to keep structures validated

	for err == nil {
		cur, err = rd.ReadByte()

		if err != nil {
			if err == io.EOF {
				r.Tokens = append(r.Tokens, Token{TokenType: EOF})
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
			rd.UnreadByte()
			t, err := r.tokenizeBool(rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case 't':
			rd.UnreadByte()
			t, err := r.tokenizeBool(rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)

		case 'n':
			rd.UnreadByte()
			t, err := r.tokenizeNull(rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case '"':
			//handle string
			t, err := r.tokenizeString(rd)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case ':':
			r.Tokens = append(r.Tokens, Token{TokenType: COLON, Value: ":"})
		case ',':
			r.Tokens = append(r.Tokens, Token{TokenType: COMMA, Value: ","})
		case ' ': //empty space skip
			continue
		default:
			if unicode.IsDigit(rune(cur)) {
				r.Reader.UnreadByte()
				t, err := r.tokenizeNumber()
				if err != nil {
					return nil, err
				}
				r.Tokens = append(r.Tokens, *t)
			} else {
				//erronous state
				return nil, fmt.Errorf("incorrect json structure")
			}
		}
	}

	return r.Tokens, nil
}

func (r *Lexer) tokenizeBool(rd *bufio.Reader) (*Token, error) {
	//look further 5 bytes false
	n, err := rd.Peek(5)
	if err != nil {
		return nil, err
	}
	ok := bytes.Equal(n, []byte("false"))
	if !ok {
		n, err = rd.Peek(4)
		if err != nil {
			return nil, err
		}
	}
	ok1 := bytes.Equal(n, []byte("true"))
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
		rd.Discard(5)
	} else {
		rd.Discard(4)
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

func (r *Lexer) tokenizeNumber() (*Token, error) {
	var buf bytes.Buffer
	for {
		d, _, err := r.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if unicode.IsDigit(d) || d == '-' || d == '+' || d == '.' || d == 'e' || d == 'E' {
			buf.WriteRune(d)
		} else {
			r.Reader.UnreadRune()
			break
		}
	}
	return &Token{TokenType: NUMBER, Value: buf.String()}, nil
}

func (r *Lexer) tokenizeNull(rd *bufio.Reader) (*Token, error) {
	n, err := rd.Peek(4)
	if err != nil {
		return nil, err
	}
	ok := bytes.Equal(n, []byte("null"))
	if !ok {
		return nil, fmt.Errorf("error parsing null")
	}

	rd.Discard(4)

	return &Token{TokenType: NULL, Value: "null"}, nil
}
