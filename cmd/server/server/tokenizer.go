package server

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"strconv"
)

type RedisTokenizer interface {
	GetTerminal() string
	NextToken(reader *bufio.Reader) (data []byte, err error)
	GetTokens(reader *bufio.Reader) (tokens []string, err error)
}

// Tokenizer is used to create tokens from the reader
type Tokenizer struct {
	KTerminal string // terminate the sequence
	*zap.Logger
}

func DefaultTokenizer() *Tokenizer {
	tok := NewTokenizer("\r\n")
	return tok
}

func NewTokenizer(terminal string) *Tokenizer {
	logger, _ := zap.NewProduction()
	return &Tokenizer{
		KTerminal: terminal,
		Logger:    logger,
	}
}

func (tok *Tokenizer) GetTerminal() string {
	return tok.KTerminal
}

// NextToken reads the next section of data from the reader.
// A section is a sequence of bytes that ends with the KTerminal value
// relative to the current offset.
func (tok *Tokenizer) NextToken(reader *bufio.Reader) (data []byte, err error) {
	// e.g +\r\nValue\r\n
	// fist call to NextToken() returns []byte("+")
	// second call to NextToken() returns []byte("Value")
	// third call will return nothing

	// store prefix of @kTerminal, if @kTerminal is absent in the reader
	// @ignored is appended to @data
	var ignored []byte
	var ch byte

	kTerminal := tok.GetTerminal()
	start, end := 0, len(kTerminal)
	current := start

	//e.g kTerminal = "<>", reader = "some_random<_value<>"
	// @ignored stores the prefix of @kTerminal
	//		data = "some_random", ignored = "<"
	// @ignored is appended to @data
	// 		data = "some_random<", ignored = ""
	// @ignored stores @kTerminal
	// 		data = "some_random<value", ignored = "<>"

	fmt.Println("-tok-", start, end)
	// read the tokens until @kTerminal is found.
	// It doesn't skip the @kTerminal. In case of strings with
	// @kTerminal as substring, NextToken() will return the substring,
	// multiple calls maybe required.
	for current < end {
		ch, err = reader.ReadByte()
		// end of reader
		if err != nil {
			if err == io.EOF {
				tok.Info("EOF")
				return data, err
			}
			tok.Error("error while reading byte", zap.Error(err))
			return nil, err
		}
		// current byte leads to @kTerminal (i.e. parsing a prefix of @kTerminal)
		// add current byte to ignored bytes
		if ch == kTerminal[current] {
			ignored = append(ignored, ch)
			current++
		} else {
			// current byte doesn't lead to @kTerminal (not a prefix of @kTerminal)
			// append the @ignored bytes to @data
			data = append(data, ignored...)
			ignored = []byte{}
			current = start
			if ch == kTerminal[current] {
				// if current byte leads to @kTerminal (prefix of @kTerminal), unread the byte.
				// as this should be added to @ignore slice and not @data slice.
				if err = reader.UnreadByte(); err != nil {
					tok.Error("error while unread byte", zap.Error(err))
					return nil, err
				}
			} else {
				// add the current byte to the data
				data = append(data, ch)
			}
		}
	}
	fmt.Println("--<>--", string(data))
	return data, err
}

func (tok *Tokenizer) GetTokens(reader *bufio.Reader) (tokens []string, err error) {
	// Read the <metadata, data>
	var data, metadata []byte

	// read the <metadata>
	if metadata, err = tok.NextToken(reader); err != nil {
		tok.Error("error while tokenizing metadata", zap.Error(err))
		return tokens, err
	}
	tokens = append(tokens, string(metadata))

	defer func() {
		tok.Info(
			"token",
			zap.String("metadata", string(metadata)),
			zap.String("data", string(data)),
			zap.Strings("tokens", tokens),
		)
	}()
	// read the data i.e the content
	if metadata[0] == '*' {
		// An array is a collection of tokens (<metadata>\r\n<data>\r\n)*
		var temp []string
		var arraySize int

		if arraySize, err = ToInt(metadata[1:]); err != nil {
			tok.Error("error while converting metadata to int", zap.Error(err))
			return nil, err
		}
		// read the next @arraySize tokens (<metadata, data>) pairs
		for arraySize > 0 {
			if temp, err = tok.GetTokens(reader); err != nil {
				tok.Error("error while tokenizing data", zap.Error(err))
				return nil, err
			}
			tokens = append(tokens, temp...)
			arraySize--
		}
		return tokens, err
	}

	// read the <data> (content)
	data, err = tok.NextToken(reader)
	if err != nil {
		if err == io.EOF {
			tok.Info("EOF received")
			return tokens, nil
		}
		tok.Error("error while tokenizing data", zap.Error(err))
		return tokens, err
	}
	tokens = append(tokens, string(data))
	return tokens, err
}

var ToInt = func(s []byte) (int, error) {
	a, err := strconv.ParseInt(string(s), 10, 0)
	return int(a), err
}
