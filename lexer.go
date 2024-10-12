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
	EOF
	SPACE
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
	for {
		cur, err := r.Reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				r.Tokens = append(r.Tokens, Token{TokenType: EOF})
				return r.Tokens, nil
			}
			return nil, err
		}

		switch cur {
		case ' ', '\n', '\r':
			// Skip whitespace
			continue
		case '{':
			r.Tokens = append(r.Tokens, Token{TokenType: LEFT_BRACE, Value: "{"})
		case '}':
			r.Tokens = append(r.Tokens, Token{TokenType: RIGHT_BRACE, Value: "}"})
		case '[':
			r.Tokens = append(r.Tokens, Token{TokenType: LEFT_BRACKET, Value: "["})
		case ']':
			r.Tokens = append(r.Tokens, Token{TokenType: RIGHT_BRACKET, Value: "]"})
		case ':':
			r.Tokens = append(r.Tokens, Token{TokenType: COLON, Value: ":"})
		case ',':
			r.Tokens = append(r.Tokens, Token{TokenType: COMMA, Value: ","})
		case '"':
			r.Reader.UnreadByte()
			t, err := r.tokenizeString()
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case 'f':
			r.Reader.UnreadByte()
			t, err := r.tokenizeBool(r.Reader)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case 't':
			r.Reader.UnreadByte()
			t, err := r.tokenizeBool(r.Reader)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)

		case 'n':
			r.Reader.UnreadByte()
			t, err := r.tokenizeNull(r.Reader)
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			r.Reader.UnreadByte()
			t, err := r.tokenizeNumber()
			if err != nil {
				return nil, err
			}
			r.Tokens = append(r.Tokens, *t)
		case '\t':
			return nil, fmt.Errorf("incorrect json structure: tab charater")
		default:
			return nil, fmt.Errorf("unexpected character: %c", cur)
		}
	}
}

// escape character only allowed for ", \, /, b, f, r, t, u
func (r *Lexer) handleEscapeCharacters() error {
	_, err := r.Reader.ReadByte()
	if err != nil {
		return fmt.Errorf("reading escape character")
	}

	esc, err := r.Reader.ReadByte()

	if esc == '\\' || esc == '"' || esc == '/' || esc == '\b' || esc == '\f' || esc == '\n' || esc == '\r' || esc == '\t' {
		return nil
	}

	return fmt.Errorf("invalid escape sequence")
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

func (r *Lexer) tokenizeString() (*Token, error) {

	rd := r.Reader

	var val []byte

	two, err := r.Reader.Peek(2)
	if err != nil {
		return nil, err
	}

	if bytes.Equal(two, []byte{'"', '"'}) {
		rd.Discard(2)
		val = append(val, '"', '"')
		return &Token{TokenType: STRING, Value: string(val)}, nil
	}

	rd.ReadByte() // Consume opening quote

	for {
		cur, err := rd.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("error while parsing string token: %v", err)
		}

		if cur == '"' {
			break
		}

		if cur == '\\' {
			next, err := rd.ReadByte()
			if err != nil {
				return nil, fmt.Errorf("error while parsing escape sequence: %v", err)
			}
			switch next {
			case '\\', '/', 'b', 'f', 'n', 'r':
				val = append(val, '\\', next)
			case 'u':
				unicodeSeq := make([]byte, 4)
				_, err := io.ReadFull(rd, unicodeSeq)
				if err != nil {
					return nil, fmt.Errorf("error while parsing unicode sequence: %v", err)
				}
				val = append(val, '\\', 'u')
				val = append(val, unicodeSeq...)
			default:
				return nil, fmt.Errorf("invalid escape sequence: \\%c", next)
			}
		} else if cur == '\t' {
			return nil, fmt.Errorf("tab character")
		} else if cur < 32 || cur > 126 {
			return nil, fmt.Errorf("non-printable character")
		} else {
			val = append(val, cur)
		}
	}

	return &Token{TokenType: STRING, Value: string(val)}, nil
}

func (r *Lexer) tokenizeNumber() (*Token, error) {
	var buf bytes.Buffer
	var hasDecimal, hasExponent bool

	// Read the first character (- or digit)
	first, _ := r.Reader.ReadByte()

	n, _ := r.Reader.Peek(1)

	if first == '0' && isValidNumberByte(n[0]) {
		return nil, fmt.Errorf("cannot have leading zeros")
	}

	buf.WriteByte(first)

	for {
		next, err := r.Reader.ReadByte()
		if err == io.EOF {
			break
		}

		switch {
		case next >= '0' && next <= '9':
			buf.WriteByte(next)
		case next == '.' && !hasDecimal && !hasExponent:
			hasDecimal = true
			buf.WriteByte(next)
		case (next == 'e' || next == 'E') && !hasExponent:
			hasExponent = true
			buf.WriteByte(next)
			// Check for + or - after E
			expSign, err := r.Reader.ReadByte()
			if err == nil && (expSign == '+' || expSign == '-') {
				buf.WriteByte(expSign)
			} else if err == nil {
				r.Reader.UnreadByte()
			}
		default:
			r.Reader.UnreadByte()
			return &Token{TokenType: NUMBER, Value: buf.String()}, nil
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

func isValidNumberByte(b byte) bool {
	return b >= '0' && b <= '9'
}
