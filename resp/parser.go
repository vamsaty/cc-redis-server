package resp

import (
	"bufio"
	ccUtils "github.com/vamsaty/cc-utils"
	"strconv"
)

// ReadToken reads a token from the reader. A token is a sequence of bytes
// terminated by \r\n. The token is returned as a slice of runes.
func ReadToken(reader *bufio.Reader) ([]rune, error) {
	var token []rune
	prevChar := rune(0)
	for {
		if char, _, err := reader.ReadRune(); err != nil {
			return nil, err
		} else if char == '\n' && prevChar == '\r' {
			break
		} else {
			token = append(token, char)
			prevChar = char
		}
	}
	return token[:len(token)-1], nil
}

// ParseArray parses a slice of runes into a 2D array of strings.
// The @token is a slice of runes returned by ReadToken function.
// The @reader parameter is the reader from which the tokens are read.
// It is the same @reader in ReadToken function
func ParseArray(token []rune, reader *bufio.Reader) [][]string {
	arraySize, err := strconv.Atoi(string(token[1:]))
	ccUtils.PanicIf(err)

	var tokens [][]string
	for arraySize > 0 {
		metadata, err := ReadToken(reader)
		ccUtils.PanicIf(err)

		data, err := ReadToken(reader)
		ccUtils.PanicIf(err)

		tokens = append(
			tokens,
			[]string{
				string(metadata),
				string(data),
			})
		arraySize--
	}
	return tokens
}
