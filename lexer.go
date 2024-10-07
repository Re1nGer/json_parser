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
		case '\\':
			continue
		case '-':
			// Handle negative numbers
			nextByte, err := rd.ReadByte()
			if err != nil {
				return nil, err
			}
			if unicode.IsDigit(rune(nextByte)) {
				rd.UnreadByte() // Put back the digit
				rd.UnreadByte() // Put back the minus sign
				t, err := r.tokenizeNumber()
				if err != nil {
					return nil, err
				}
				r.Tokens = append(r.Tokens, *t)
			} else {
				return nil, fmt.Errorf("invalid character after minus sign: %c", nextByte)
			}
		default:
			if unicode.IsSpace(rune(cur)) {
				continue
			}
			fmt.Println("Unknown elements", string(cur), cur)
			if unicode.IsDigit(rune(cur)) {
				r.Reader.UnreadByte()
				t, err := r.tokenizeNumber()
				if err != nil {
					return nil, err
				}
				r.Tokens = append(r.Tokens, *t)
			} else {
				//erronous state
				fmt.Println("error element", cur, string(cur), r.Tokens)
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

	//handle "" case
	var val []byte
	var cur_val byte

	n, err := rd.Peek(1)

	if n[0] == '"' {
		rd.Discard(1)
		return &Token{TokenType: STRING, Value: ""}, nil
	}

	cur, err := rd.ReadByte()
	if err != nil {
		return nil, err
	}

	//fmt.Println("cur value", string(cur))

	if cur == '\\' || cur == '"' {
		//skip
		cur, err = rd.ReadByte()
	} else {
		val = append(val, cur)
	}

	//has to handle \" escape double string case
	t := &Token{}

	for cur_val != '"' && err == nil {
		next, _ := rd.Peek(1)
		if cur_val == '\t' || cur_val == '\r' || cur_val == '\b' || cur_val == '\f' || cur_val == '\n' || cur_val == '\u0022' {
			rd.Discard(1)
			continue
		}
		if cur_val == '\\' && next[0] == '"' {
			//skip escape string
			rd.Discard(1)
		}
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
	isFirstDigit := true
	hasLeadingZero := false
	isFloat := false
	isExponent := false

	for {
		d, _, err := r.Reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if unicode.IsDigit(d) {
			if isFirstDigit {
				if d == '0' {
					hasLeadingZero = true
				}
				isFirstDigit = false
			} else if hasLeadingZero && !isFloat && !isExponent {
				// Check for hexadecimal numbers
				if buf.Len() == 1 { // We've only seen the leading '0' so far
					nextChar, _, err := r.Reader.ReadRune()
					if err == nil {
						if nextChar == 'x' || nextChar == 'X' {
							return nil, fmt.Errorf("hexadecimal numbers are not allowed")
						}
						r.Reader.UnreadRune() // Put back the character we just read
					}
				}
				return nil, fmt.Errorf("invalid number format: leading zero")
			}
			buf.WriteRune(d)
		} else if d == '-' || d == '+' {
			if buf.Len() > 0 && !isExponent {
				r.Reader.UnreadRune()
				break
			}
			buf.WriteRune(d)
		} else if d == '.' {
			if isFloat {
				return nil, fmt.Errorf("invalid number format: multiple decimal points")
			}
			isFloat = true
			buf.WriteRune(d)
		} else if d == 'e' || d == 'E' {
			if isExponent {
				return nil, fmt.Errorf("invalid number format: multiple exponents")
			}
			isExponent = true
			isFloat = true // Treat numbers with exponents as floats
			buf.WriteRune(d)
		} else {
			r.Reader.UnreadRune()
			break
		}
	}

	// Final check for hex numbers at the end of parsing
	if hasLeadingZero && buf.Len() > 1 && (buf.Bytes()[1] == 'x' || buf.Bytes()[1] == 'X') {
		return nil, fmt.Errorf("hexadecimal numbers are not allowed")
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
