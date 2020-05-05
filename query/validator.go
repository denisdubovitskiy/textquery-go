package query

type ValidationError struct {
	position  int
	message   string
	character string
}

func (v ValidationError) Error() string {
	return v.message
}

func (v ValidationError) Pos() int {
	return v.position
}

func (v ValidationError) Char() string {
	return v.character
}

func (p *Parser) isBalanced(query string) error {
	stack := parenthesisStack{}
	index := 0

	for index < len(query) {
		character := string(query[index])

		if p.allOpenParens[character] {
			stack.push(parenthesis{character: character, position: index})
		} else if p.allCloseParens[character] {

			if stack.isEmpty() {
				return ValidationError{
					position:  index,
					character: character,
					message:   "open parenthesis is not found",
				}
			}

			openParentheses, _ := stack.pop()
			if p.closeToOpenParens[character] != openParentheses.character {
				return ValidationError{
					position:  index,
					character: character,
					message:   "open parenthesis is not found",
				}
			}
		}

		index++
	}

	if !stack.isEmpty() {
		parens, _ := stack.pop()
		return ValidationError{
			position:  parens.position,
			character: parens.character,
			message:   "close parenthesis is not found",
		}
	}

	return nil
}
