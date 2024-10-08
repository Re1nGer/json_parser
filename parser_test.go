package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestSimpleJsonParser(t *testing.T) {
	sample := []byte("{\"key\": \"value\"}")

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
		t.Errorf("error instantiating parser %v %v", err, parser.tokens)
	}

	_, err = parser.Parse()
	t.Log("tokens", parser.tokens, parser.stack)
	if err != nil {
		t.Errorf("error parsing %v", err)
	}
}

func TestNestedWithArrayJson(t *testing.T) {

	sample := []byte("{\"arr\":[{\"nested\":\"jee\"}, { \"another\" : \"val\" }]}")

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

func TestFailExtraClose(t *testing.T) {
	sample := []byte("[\"Extra close\"]]")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailExtraCommaTrue(t *testing.T) {
	sample := []byte("{\"Extra comma\": true, }")

	p, err := NewParser(sample)

	err = p.parseObj()

	fmt.Println("error raised", err, p.tokens)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailMisplacedQuoteValue(t *testing.T) {
	sample := []byte("{\"Extra value after close\": true} \"misplaced quoted value\"")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err, p.tokens)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailRandomStructure(t *testing.T) {

	sample := []byte("{\"Extra value after close\": true} \"misplaced quoted value\"")

	p, err := NewParser(sample)

	_, err = p.Parse()

	fmt.Println("error raised", err, p.tokens)

	if err == nil {
		t.Errorf("error should have been raised")
	}
}

func TestFailIllegalExpression(t *testing.T) {
	sample := []byte("{\"Illegal expression\":1+2}")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}

	//_, err = p.Parse()

	fmt.Println("error raised", err)
}

func TestFailINumbersCannotHaveLeadingZero(t *testing.T) {
	sample := []byte("{\"Numbers cannot have leading zeroes\": 013}")

	p, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised %v", p.tokens)
	}

	fmt.Println("error raised", err)
}

func TestFailNumbersCannotBeHex(t *testing.T) {
	sample := []byte("{\"Numbers cannot be hex\": 0x14}")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

// separate case to handle on lexer level
func TestFailIllegalBackslashEscape(t *testing.T) {
	sample := []byte("[\"Illegal backslash escape: \x15\"]")

	_, err := NewParser(sample)

	if err == nil {
		t.Fail()
	}

	fmt.Println("error raised", err)
}

func TestFailIllegalBackslashEscape2(t *testing.T) {
	sample := []byte("[\"Illegal backslash escape: \017\"]")

	p, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised: %v", p.tokens)
	}

	fmt.Println("error raised", err)
}

func TestFailIlTooDeep(t *testing.T) {
	sample := []byte("[[[[[[[[[[[[[[[[[[[[\"Not so Too deep\"]]]]]]]]]]]]]]]]]]]]")

	p, err := NewParser(sample)

	ok, err := p.Parse()

	if !ok {
		t.Fail()
	}

	if err != nil {
		t.Errorf("error  %v", err)
	}

	fmt.Println("error raised", err)
}

func TestFailMissingColon(t *testing.T) {
	sample := []byte("{\"Missing colon\" null}")

	p, err := NewParser(sample)

	_, err = p.Parse()

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

func TestFailDoubleColon(t *testing.T) {
	sample := []byte("{\"Double colon\":: null}")

	p, err := NewParser(sample)

	_, err = p.Parse()

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}
func TestFailCommaInsteadOfColon(t *testing.T) {
	sample := []byte("{\"Comma instead of colon\", null}")

	p, err := NewParser(sample)

	_, err = p.Parse()

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

func TestFailCommaInsteadOfColon2(t *testing.T) {
	sample := []byte("[\"Colon instead of comma\": false]")

	p, err := NewParser(sample)

	_, err = p.Parse()

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

func TestFailBadValue(t *testing.T) {
	sample := []byte("[\"Bad value\", truth]")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

func TestFailSingleQuote(t *testing.T) {
	sample := []byte("['single quote']")

	_, err := NewParser(sample)

	if err == nil {
		t.Errorf("error should have been raised")
	}

	fmt.Println("error raised", err)
}

// weird one shoud investigate
func TestFailTabCharacters(t *testing.T) {
	sample := []byte("[\"	tab	character	in	string	\"]")

	_, err := NewParser(sample)

	if err == nil {
		t.Fail()
	}

	fmt.Println("error raised", err)
}

func TestShouldPass(t *testing.T) {
	sample := []byte("{ \"JSON Test Pattern pass3\": { \"The outermost value\": \"must be an object or array.\", \"In this test\": \"It is an object.\" } } ")

	_, err := NewParser(sample)

	if err != nil {
		t.Errorf("error %v", err)
	}
}

func TestShouldPassNotDeep(t *testing.T) {
	sample := []byte("[[[[[[[[[[[[[[[[[[[[[[[[[[[[[\"not too deep\"]]]]]]]]]]]]]]]]]]]]]]]]]]]]]]")

	p, err := NewParser(sample)

	err = p.parseArray()

	if err != nil {
		t.Errorf("error %v", err)
		t.Fail()
	}

	if err != nil {
		t.Errorf("error %v", err)
	}

	t.Log(p.stack)
}

func TestParseObj(t *testing.T) {
	sample := []byte("{\"key\": { \"nested\": \"value\", \"nested2\": \"value1\" }}")

	p, _ := NewParser(sample)

	err := p.ParseObj()
	if err != nil {
		t.Errorf("error %v", err)
		t.Fail()
	}

	fmt.Println("stack", p.stack, "tokens", p.tokens)
}

func TestParseArray(t *testing.T) {

	sample := []byte("[\"elements\", \"elements1\", \"elements2\"]")

	p, _ := NewParser(sample)

	err := p.parseArray()
	if err != nil {
		t.Errorf("error %v", err)
		t.Fail()
	}

	fmt.Println("stack", p.stack, "tokens", p.tokens)
}

func TestParseComplexArray(t *testing.T) {

	sample := []byte("[\"elements\", { \"object\": \"val\", \"key\": { \"another\": \"value\" } }]")

	p, _ := NewParser(sample)

	err := p.parseArray()
	if err != nil {
		t.Errorf("error %v", err)
		t.Fail()
	}

	fmt.Println("stack", p.stack, "tokens", p.tokens)
}

func TestDeeplyNestedObj(t *testing.T) {

	path := "./main/testpass"

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		if !file.IsDir() && file.Name() == "pass2.json" {
			filename := file.Name()
			filePath := filepath.Join(path, filename)
			fmt.Printf("Processing file: %s\n", filePath)

			// Read file content
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", filePath, err)
				continue
			}

			fmt.Println("content", string(content))

			p, err := NewParser(content)

			if err != nil {
				fmt.Printf("Error parsing JSON in file %s: %v\n %v\n", filePath, err, p.tokens)
				continue
			}

			err = p.ParseObj()

			if err != nil {
				fmt.Printf("Error parsing JSON in file %s: %v\n %v", filePath, err, p.tokens)
			} else {
				fmt.Printf("Successfully parsed JSON in file %s\n %v\n, stack: %v", filePath, p.GetTokens(), p.stack)
			}
		}
	}
}
