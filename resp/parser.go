package resp

import (
	"bufio"
	. "github.com/vamsaty/cc-redis-server/utils"
	"strconv"
)

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

func ParseArray(token []rune, reader *bufio.Reader) [][]string {
	arraySize, err := strconv.Atoi(string(token[1:]))
	PanicIf(err)

	var tokens [][]string
	for arraySize > 0 {
		metadata, err := ReadToken(reader)
		PanicIf(err)

		data, err := ReadToken(reader)
		PanicIf(err)

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
