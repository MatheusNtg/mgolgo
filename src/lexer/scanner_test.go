package lexer

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}

func TestScanNumToken(t *testing.T) {
	testCases := []struct {
		name           string
		preparedText   string
		expectedTokens []Token
	}{
		{
			name:           "Integer number",
			preparedText:   "1",
			expectedTokens: []Token{NewToken(NUM, "1", INTEGER)},
		},
		{
			name:           "Integer number with N digits",
			preparedText:   "123",
			expectedTokens: []Token{NewToken(NUM, "123", INTEGER)},
		},
		{
			name:           "Real number",
			preparedText:   "1.0",
			expectedTokens: []Token{NewToken(NUM, "1.0", REAL)},
		},
		{
			name:           "Real number with N digits after point",
			preparedText:   "1.000",
			expectedTokens: []Token{NewToken(NUM, "1.000", REAL)},
		},
		{
			name:           "Real number with N digits before point",
			preparedText:   "123.0",
			expectedTokens: []Token{NewToken(NUM, "123.0", REAL)},
		},
		{
			name:           "Real number with N digits before and after point",
			preparedText:   "123.000",
			expectedTokens: []Token{NewToken(NUM, "123.000", REAL)},
		},
		{
			name:           "Integer with capital exponential",
			preparedText:   "1E0",
			expectedTokens: []Token{NewToken(NUM, "1E0", INTEGER)},
		},
		{
			name:           "Integer with lower exponential",
			preparedText:   "1e0",
			expectedTokens: []Token{NewToken(NUM, "1e0", INTEGER)},
		},
		{
			name:           "Incomplete real number with capital exponential",
			preparedText:   "1.0E0",
			expectedTokens: []Token{NewToken(NUM, "1.0E0", REAL)},
		},
		{
			name:           "Incomplete real number with capital exponential",
			preparedText:   "1.0e0",
			expectedTokens: []Token{NewToken(NUM, "1.0e0", REAL)},
		},
		{
			name:           "Integer with capital exponential positive",
			preparedText:   "1E+0",
			expectedTokens: []Token{NewToken(NUM, "1E+0", INTEGER)},
		},
		{
			name:           "Integer with lower exponential positive",
			preparedText:   "1e+0",
			expectedTokens: []Token{NewToken(NUM, "1e+0", INTEGER)},
		},
		{
			name:           "Integer with capital exponential negative",
			preparedText:   "1E-0",
			expectedTokens: []Token{NewToken(NUM, "1E-0", INTEGER)},
		},
		{
			name:           "Integer with lower exponential negative",
			preparedText:   "1e-0",
			expectedTokens: []Token{NewToken(NUM, "1e-0", INTEGER)},
		},
		{
			name:           "Incomplete real number with capital exponential positive",
			preparedText:   "1.0E+0",
			expectedTokens: []Token{NewToken(NUM, "1.0E+0", REAL)},
		},
		{
			name:           "Incomplete real number with capital exponential positive",
			preparedText:   "1.0e+0",
			expectedTokens: []Token{NewToken(NUM, "1.0e+0", REAL)},
		},
		{
			name:           "Incomplete real number with capital exponential negative",
			preparedText:   "1.0E-0",
			expectedTokens: []Token{NewToken(NUM, "1.0E-0", REAL)},
		},
		{
			name:           "Incomplete real number with lower exponential negative",
			preparedText:   "1.0e-0",
			expectedTokens: []Token{NewToken(NUM, "1.0e-0", REAL)},
		},
		{
			name:         "Error incomplete real number with capital exponential positive",
			preparedText: "1.E+0",
			expectedTokens: []Token{
				ERROR_TOKEN,
				NewToken(IDENTIFIER, "E", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(NUM, "0", INTEGER),
			},
		},
		{
			name:         "Error incomplete real number with capital exponential positive",
			preparedText: "1.e+0",
			expectedTokens: []Token{
				ERROR_TOKEN,
				NewToken(IDENTIFIER, "e", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(NUM, "0", INTEGER),
			},
		},
		{
			name:         "Error incomplete real number with capital exponential negative",
			preparedText: "1.E-0",
			expectedTokens: []Token{
				ERROR_TOKEN,
				NewToken(IDENTIFIER, "E", NULL),
				NewToken(ARIT_OP, "-", NULL),
				NewToken(NUM, "0", INTEGER),
			},
		},
		{
			name:         "Error incomplete real number with lower exponential negative",
			preparedText: "1.e-0",
			expectedTokens: []Token{
				ERROR_TOKEN,
				NewToken(IDENTIFIER, "e", NULL),
				NewToken(ARIT_OP, "-", NULL),
				NewToken(NUM, "0", INTEGER),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			tokens := []Token{}
			for {
				token, _, _ := scanner.Scan()
				if token == EOF_TOKEN {
					break
				}
				tokens = append(tokens, token)
			}

			require.Equal(t, tc.expectedTokens, tokens)
		})
	}
}

func TestScanIdToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken Token
	}{
		{
			name:          "Identifier with number",
			preparedText:  "id1",
			expectedToken: NewToken(IDENTIFIER, "id1", NULL),
		},
		{
			name:          "Identifier with underline and number",
			preparedText:  "id_1",
			expectedToken: NewToken(IDENTIFIER, "id_1", NULL),
		},
		{
			name:          "Identifier with underline and more than one number",
			preparedText:  "id_123",
			expectedToken: NewToken(IDENTIFIER, "id_123", NULL),
		},
		{
			name:          "Identifier with underline and more than one number and more than one character",
			preparedText:  "id_123_id",
			expectedToken: NewToken(IDENTIFIER, "id_123_id", NULL),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			token, _, _ := scanner.Scan()

			require.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestScanCommentToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken []Token
	}{
		{
			name:         "Valid comment with N open brackets",
			preparedText: "{{{ab}",
			expectedToken: []Token{
				COMMENT_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Close comment twice with characters in between",
			preparedText: "{ab}ab}",
			expectedToken: []Token{
				COMMENT_TOKEN,
				ERROR_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Close comment twice",
			preparedText: "{{abab}}",
			expectedToken: []Token{
				COMMENT_TOKEN,
				ERROR_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Comment not closed",
			preparedText: "{{abab",
			expectedToken: []Token{
				ERROR_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Comment not open",
			preparedText: "abab}}",
			expectedToken: []Token{
				ERROR_TOKEN,
				ERROR_TOKEN,
				EOF_TOKEN,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())

			for _, expectedToken := range tc.expectedToken {
				token, _, _ := scanner.Scan()
				require.Equal(t, expectedToken, token)
			}
		})
	}
}

func TestScanLiteralConstantToken(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken Token
	}{
		{
			name:          "Simple Constant Literal",
			preparedText:  `"This is a constant literal"`,
			expectedToken: NewToken(LITERAL_CONST, `"This is a constant literal"`, LITERAL),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())
			token, _, _ := scanner.Scan()

			require.Equal(t, tc.expectedToken, token)
		})
	}
}

func TestScanGeneralCases(t *testing.T) {
	testCases := []struct {
		name          string
		preparedText  string
		expectedToken []Token
	}{
		{
			name:         "Assignment",
			preparedText: "A<-B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				ATTR_TOKEN,
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Assignment with sum",
			preparedText: "A<-B+C",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				ATTR_TOKEN,
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Assignment with subtraction",
			preparedText: "A<-B-C",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				ATTR_TOKEN,
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(ARIT_OP, "-", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Less than or greater than",
			preparedText: "A<>B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Less than or equal",
			preparedText: "A<=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Greater than or equal",
			preparedText: "A>=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, ">=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Equal",
			preparedText: "A=B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "=", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Less than",
			preparedText: "A<B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, "<", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Greater than",
			preparedText: "A>B",
			expectedToken: []Token{
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(REL_OP, ">", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				EOF_TOKEN,
			},
		},
		{
			name:         "Operation with comparison between parentheses",
			preparedText: "(A+B<>C)",
			expectedToken: []Token{
				OPEN_PAR_TOKEN,
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				CLOSE_PAR_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Two Operations with comparisons between parentheses and semicolon",
			preparedText: "(A+B<>C/D);",
			expectedToken: []Token{
				OPEN_PAR_TOKEN,
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(ARIT_OP, "/", NULL),
				NewToken(IDENTIFIER, "D", NULL),
				CLOSE_PAR_TOKEN,
				SEMICOLON_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Two Operations with comparisons between parentheses and semicolon",
			preparedText: "se(A+B<>C/D);",
			expectedToken: []Token{
				NewToken("se", "se", "se"),
				OPEN_PAR_TOKEN,
				NewToken(IDENTIFIER, "A", NULL),
				NewToken(ARIT_OP, "+", NULL),
				NewToken(IDENTIFIER, "B", NULL),
				NewToken(REL_OP, "<>", NULL),
				NewToken(IDENTIFIER, "C", NULL),
				NewToken(ARIT_OP, "/", NULL),
				NewToken(IDENTIFIER, "D", NULL),
				CLOSE_PAR_TOKEN,
				SEMICOLON_TOKEN,
				EOF_TOKEN,
			},
		},
		{
			name:         "Escreva with jump line",
			preparedText: `escreva "\nA=\n";`,
			expectedToken: []Token{
				NewToken("escreva", "escreva", "escreva"),
				NewToken(LITERAL_CONST, `"\nA=\n"`, LITERAL),
				SEMICOLON_TOKEN,
				EOF_TOKEN,
			},
		},
	}

	symbolTable := GetSymbolTableInstance()

	FillSymbolTable(symbolTable)
	defer symbolTable.Cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())

			for _, expectedToken := range tc.expectedToken {
				token, _, _ := scanner.Scan()
				require.Equal(t, expectedToken, token)
			}
		})
	}
}

func TestStdoutErrorLog(t *testing.T) {
	testCases := []struct {
		name           string
		preparedText   string
		expectedOutput []string
	}{
		{
			name:         "Character does not exits in the alphabet in one line",
			preparedText: "abc %",
			expectedOutput: []string{
				"",
				"erro na linha 1 coluna 5, palavra % inexistente na linguagem",
			},
		},
		{
			name:         "Character does not exits in the alphabet in one line inside a word",
			preparedText: "abc%",
			expectedOutput: []string{
				"erro na linha 1 coluna 4, palavra abc% inexistente na linguagem",
			},
		},
		{
			name:         "Character does not exits in the alphabet in the second line",
			preparedText: "A<-3;\nB %",
			expectedOutput: []string{
				"",
				"",
				"",
				"",
				"",
				"erro na linha 2 coluna 3, palavra % inexistente na linguagem",
			},
		},
		{
			name:         "Character does not exits in the alphabet with N line breaks, first column",
			preparedText: "\n\n\n$",
			expectedOutput: []string{
				"erro na linha 4 coluna 1, palavra $ inexistente na linguagem",
			},
		},
		{
			name:         "Character does not exits in the alphabet with N line breaks, Mth column",
			preparedText: "A\nB\nC\nD<-E ; $",
			expectedOutput: []string{
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"erro na linha 4 coluna 8, palavra $ inexistente na linguagem",
			},
		},
		{
			name:         "Malformated number",
			preparedText: "1.e3",
			expectedOutput: []string{
				"erro na linha 1 coluna 3, número 1. inválido",
				"",
			},
		},
		{
			name:         "Malformated number with double points",
			preparedText: "1..0",
			expectedOutput: []string{
				"erro na linha 1 coluna 3, número 1. inválido",
				"erro na linha 1 coluna 4, palavra . inexistente na linguagem",
			},
		},
		{
			name:         "Malformated literal",
			preparedText: `"this is a malformated literal`,
			expectedOutput: []string{
				"erro na linha 1 coluna 30, literal \"this is a malformated literal inválido",
				"",
			},
		},
		{
			name:         "Malformated comment",
			preparedText: "{this is malformated commment",
			expectedOutput: []string{
				"erro na linha 1 coluna 29, comentário {this is malformated commment inválido",
				"",
			},
		},
		{
			name:         "State 0 with no transition and lexembuffer empty",
			preparedText: "!!",
			expectedOutput: []string{
				"erro na linha 1 coluna 1, palavra ! inexistente na linguagem",
				"erro na linha 1 coluna 2, palavra ! inexistente na linguagem",
			},
		},
		{
			name:         "State 0 with no transition and lexembuffer empty after reading something",
			preparedText: "123!!",
			expectedOutput: []string{
				"",
				"erro na linha 1 coluna 4, palavra ! inexistente na linguagem",
				"erro na linha 1 coluna 5, palavra ! inexistente na linguagem",
			},
		},
	}

	symbolTable := GetSymbolTableInstance()

	FillSymbolTable(symbolTable)
	defer symbolTable.Cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			file, err := ioutil.TempFile("", "scan-test")
			require.NoError(t, err)
			defer file.Close()

			_, err = file.WriteString(tc.preparedText)
			require.NoError(t, err)

			file.Seek(0, io.SeekStart)

			scanner := NewScanner(file, GetSymbolTableInstance())

			for _, expectedOutput := range tc.expectedOutput {
				output := captureOutput(func() { scanner.Scan() })
				// Remove date, hour and line break
				if output != "" {
					output = output[20 : len(output)-1]
				}
				require.Equal(t, expectedOutput, output)
			}
		})
	}
}
