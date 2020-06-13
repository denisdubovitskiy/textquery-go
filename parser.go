package textquery

type associativity int8

const (
	AND              = "AND"
	OR               = "OR"
	NOT              = "NOT"
	fieldDelimiter   = ":"
	leftParen        = "("
	rightParen       = ")"
	replaceDelimiter = "_-_-"

	right associativity = 1
	left  associativity = 2
)

type (
	Node struct {
		Key   Token
		Left  *Node
		Right *Node
	}

	operator struct {
		precedence    int
		associativity associativity
	}
)

var operators = map[string]operator{
	NOT: {precedence: 3, associativity: right},
	AND: {precedence: 2, associativity: left},
	OR:  {precedence: 1, associativity: left},
}

func hasPrecedence(a, b Token) bool {
	opA := operators[a.Key]
	opB := operators[b.Key]

	isLeftAssociative := opB.associativity == left && opA.precedence >= opB.precedence
	isRightAssociative := opB.associativity == right && opA.precedence > opB.precedence

	return isLeftAssociative || isRightAssociative
}

// shunting-yard algorithm implementation
func parseSearchQuery(tokens []Token) []Token {
	var reverseNotation []Token
	operatorsStack := newTokenStack(len(tokens))

	for _, token := range tokens {
		if isOperand(token) {
			reverseNotation = append(reverseNotation, token)
			continue
		}

		if isOpenParen(token) {
			operatorsStack.push(token)
			continue
		}

		if isCloseParen(token) {
			for {
				current := operatorsStack.pop()
				if isOpenParen(current) {
					break
				}

				reverseNotation = append(reverseNotation, current)
			}

			continue
		}

		if isOperator(token) {
			for !operatorsStack.isEmpty() {
				topToken := operatorsStack.peek()

				if !isOperator(topToken) {
					break
				}

				if !hasPrecedence(topToken, token) {
					break
				}

				reverseNotation = append(reverseNotation, operatorsStack.pop())
			}

			operatorsStack.push(token)
			continue
		}

	}

	for !operatorsStack.isEmpty() {
		reverseNotation = append(reverseNotation, operatorsStack.pop())
	}

	return reverseNotation
}

func constructBinaryTree(reverseNotation []Token) Node {
	stack := newNodeStack(len(reverseNotation))

	for _, token := range reverseNotation {
		if !isOperator(token) {
			stack.push(Node{Key: token})
			continue
		}

		right := stack.pop()

		if isNot(token) {
			stack.push(Node{Key: token, Right: &right})
			continue
		}

		left := stack.pop()

		stack.push(Node{
			Key:   token,
			Left:  &left,
			Right: &right,
		})
	}

	return stack.pop()
}

func Parse(query string) Node {
	return constructBinaryTree(parseSearchQuery(tokenize(query)))
}
