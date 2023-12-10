package resp

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func padTokens(tokens []string) string {
	var data []string
	for _, token := range tokens {
		tok := []rune(token)
		tok = append(tok, '\r', '\n')
		data = append(data, string(tok))
	}
	return strings.Join(data, "")
}

func getReader(tokens []string) *bufio.Reader {
	return bufio.NewReader(
		strings.NewReader(
			padTokens(tokens),
		),
	)
}

func TestParseArray(t *testing.T) {

	type TestCase struct {
		name  string
		input []string
		want  [][]string
	}

	testCases := []TestCase{
		{
			name: "Parse array",
			input: []string{
				"*2", "$3", "foo", "$3", "bar",
			},
			want: [][]string{
				{"$3", "foo"},
				{"$3", "bar"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bufReader := getReader(tc.input)
			tok, _ := ReadToken(bufReader)

			response := ParseArray(tok, bufReader)
			if !reflect.DeepEqual(response, tc.want) {
				t.Errorf("ParseArray() = %v, want %v", response, tc.want)
			}
		})
	}
}
