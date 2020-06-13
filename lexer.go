package textquery

import "strings"

const (
	// AND ...
	AND = "AND"
	// OR ...
	OR = "OR"
	// NOT ...
	NOT = "NOT"

	// fieldDelimiter is used inside <field_name>:<field_value> constructs
	fieldDelimiter = ":"

	leftParen  = "("
	rightParen = ")"

	// replaceDelimiter is used in splitting query into parts
	replaceDelimiter = "_-_-"
)

var (
	// operatorReplacements has set of operators
	// wrapped by synthetic delimiter which is then
	// going to be used in query splitting
	operatorReplacements = []string{
		AND, formatReplacement(AND),
		OR, formatReplacement(OR),
		NOT, formatReplacement(NOT),

		leftParen, formatReplacement(leftParen),
		rightParen, formatReplacement(rightParen),
	}

	operatorReplacer = strings.NewReplacer(operatorReplacements...)

	// fieldPartsReplacements is just like operatorReplacements but
	// only used when the field delimiter has been found
	fieldPartsReplacements = []string{
		"[", formatReplacement("["),
		"]", formatReplacement("]"),
		"{", formatReplacement("{"),
		"}", formatReplacement("}"),
		":", formatReplacement(":"),
	}

	fieldPartsReplacer = strings.NewReplacer(fieldPartsReplacements...)
)

// Token ...
type Token struct {
	Key      string
	Field    string
	Operator string
	Modifier string
}

func (t *Token) isOperator() bool {
	return t.Key == AND || t.Key == OR || t.Key == NOT
}

func (t *Token) isOperand() bool {
	return !(t.isOperator() || t.isParen())
}

func (t *Token) isParen() bool {
	return t.isOpenParen() || t.isCloseParen()
}

func (t *Token) isOpenParen() bool {
	return t.Key == leftParen
}

func (t *Token) isCloseParen() bool {
	return t.Key == rightParen
}

// hasFieldDelimiter checks if the source string has a delimiter
// like field_name:value
func hasFieldDelimiter(source string) bool {
	return strings.Contains(source, fieldDelimiter)
}

// replace is used to replace operators to the ones
// with synthetic marks around them for easy splitting
func replace(query string, replacer *strings.Replacer) []string {
	tokenized := strings.Split(replacer.Replace(query), replaceDelimiter)
	tokens := make([]string, 0, len(tokenized))

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

// tokenizeFieldParts
func tokenizeFieldParts(field string) []string {
	return replace(field, fieldPartsReplacer)
}

// tokenize splits the whole query to the individual tokens
// which include operators, parens, search phrase parts
func tokenize(query string) []*Token {
	var tokens []*Token

	for _, token := range tokenizeQueryByOperators(query) {

		if !hasFieldDelimiter(token) {
			tokens = append(tokens, &Token{Key: token})
			continue
		}

		// field[precedence]:value
		fieldParts := tokenizeFieldParts(token)
		if len(fieldParts) == 6 {
			tokens = append(tokens, &Token{
				Key:      fieldParts[5],
				Field:    fieldParts[0],
				Operator: fieldParts[2],
			})
			continue
		}

		// field{modifier}[precedence]:value
		if len(fieldParts) == 9 {
			tokens = append(tokens, &Token{
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

// formatReplacement wraps the source by synthetic wrapper
// which is going to be used to split query
func formatReplacement(source string) string {
	return replaceDelimiter + source + replaceDelimiter
}
