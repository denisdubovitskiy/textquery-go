package textquery

import "strings"

var (
	operatorReplacements = []string{
		AND, formatReplacement(AND),
		OR, formatReplacement(OR),
		NOT, formatReplacement(NOT),
		leftParen, formatReplacement(leftParen),
		rightParen, formatReplacement(rightParen),
	}

	operatorReplacer = strings.NewReplacer(operatorReplacements...)

	fieldPartsReplacements = []string{
		"[", formatReplacement("["),
		"]", formatReplacement("]"),
		"{", formatReplacement("{"),
		"}", formatReplacement("}"),
		":", formatReplacement(":"),
	}

	fieldPartsReplacer = strings.NewReplacer(fieldPartsReplacements...)
)

type Token struct {
	Key      string
	Field    string
	Operator string
	Modifier string
}

func isOperator(token Token) bool {
	return token.Key == AND || token.Key == OR || token.Key == NOT
}

func isOperand(token Token) bool {
	return token.Key != AND &&
		token.Key != OR &&
		token.Key != NOT &&
		token.Key != leftParen &&
		token.Key != rightParen
}

func isOpenParen(token Token) bool {
	return token.Key == leftParen
}

func isCloseParen(token Token) bool {
	return token.Key == rightParen
}

func isNot(token Token) bool {
	return token.Key == NOT
}

func hasDelimiter(source string) bool {
	return strings.Contains(source, fieldDelimiter)
}

func replace(query string, replacer *strings.Replacer) []string {
	tokenized := strings.Split(replacer.Replace(query), replaceDelimiter)
	var tokens []string
	for _, token := range tokenized {
		trimmed := strings.TrimSpace(token)
		if trimmed == "" {
			continue
		}
		tokens = append(tokens, trimmed)
	}

	return tokens
}

func tokenizeQueryByOperators(query string) []string {
	return replace(query, operatorReplacer)
}

func tokenizeFieldParts(field string) []string {
	return replace(field, fieldPartsReplacer)
}

func tokenize(query string) []Token {
	var tokens []Token

	for _, token := range tokenizeQueryByOperators(query) {

		if !hasDelimiter(token) {
			tokens = append(tokens, Token{Key: token})
			continue
		}

		// field[operator]:value
		fieldParts := tokenizeFieldParts(token)
		if len(fieldParts) == 6 {
			tokens = append(tokens, Token{
				Key:      fieldParts[5],
				Field:    fieldParts[0],
				Operator: fieldParts[2],
			})
			continue
		}

		// field{modifier}[operator]:value
		if len(fieldParts) == 9 {
			tokens = append(tokens, Token{
				Key:      fieldParts[8],
				Field:    fieldParts[0],
				Operator: fieldParts[5],
				Modifier: fieldParts[2],
			})
			continue
		}
	}

	return tokens
}

func formatReplacement(source string) string {
	return replaceDelimiter + source + replaceDelimiter
}
