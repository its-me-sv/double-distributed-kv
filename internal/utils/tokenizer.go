package utils

import (
	"errors"
	"strings"
	"unicode"
)

type QuoteState uint8

const (
	QuoteStateOutside QuoteState = iota
	QuoteStateInsideDouble
	QuoteStateInsideSingle
)

func Tokenize(input string) ([]string, error) {
	tokens := []string{}
	var currToken strings.Builder
	var quoteState QuoteState

	tokenStarted, escaped := false, false

	for _, char := range input {
		if escaped {
			currToken.WriteRune(char)
			escaped, tokenStarted = false, false
			continue
		}

		if char == '\\' && quoteState != QuoteStateInsideSingle {
			escaped, tokenStarted = true, true
		}

		if quoteState != QuoteStateOutside {
			switch quoteState {
			case QuoteStateInsideDouble:
				if char == '"' {
					quoteState = QuoteStateOutside
				} else {
					currToken.WriteRune(char)
				}

			case QuoteStateInsideSingle:
				if char == '\'' {
					quoteState = QuoteStateOutside
				} else {
					currToken.WriteRune(char)
				}
			}

			continue
		}

		switch char {
		case '"':
			quoteState = QuoteStateInsideDouble
			tokenStarted = true
			continue

		case '\'':
			quoteState = QuoteStateInsideSingle
			tokenStarted = true
			continue
		}

		if unicode.IsSpace(char) {
			if tokenStarted {
				tokens = append(tokens, currToken.String())
				currToken.Reset()
				tokenStarted = false
			}
			continue
		}

		currToken.WriteRune(char)
		tokenStarted = true
	}

	if escaped || quoteState != QuoteStateOutside {
		return nil, errors.New("invalid")
	}

	if tokenStarted {
		tokens = append(tokens, currToken.String())
	}

	return tokens, nil
}
