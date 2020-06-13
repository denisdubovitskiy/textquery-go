package textquery

type associativity int8

const (
	right associativity = 1
	left  associativity = 2
)

type (
	// Node is a binary tree node
	Node struct {
		Key   *Token
		Left  *Node
		Right *Node
	}

	precedence struct {
		precedence    int
		associativity associativity
	}
)

var precedences = map[string]precedence{
	NOT: {precedence: 3, associativity: right},
	AND: {precedence: 2, associativity: left},
	OR:  {precedence: 1, associativity: left},
}

func hasPrecedence(a, b *Token) bool {
	opA := precedences[a.Key]
	opB := precedences[b.Key]

	isLeftAssociative := opB.associativity == left && opA.precedence >= opB.precedence
	isRightAssociative := opB.associativity == right && opA.precedence > opB.precedence

	return isLeftAssociative || isRightAssociative
}

// shunting-yard algorithm implementation
func parseSearchQuery(tokens []*Token) []*Token {
	// reverse polish notation
	var reverseNotation []*Token

	operatorsStack := newTokenStack(len(tokens))

	for _, token := range tokens {
		if token.isOperand() {
			reverseNotation = append(reverseNotation, token)
			continue
		}

		if token.isOpenParen() {
			operatorsStack.push(token)
			continue
		}

		if token.isCloseParen() {
			for {
				currentToken := operatorsStack.pop()

				if currentToken.isOpenParen() {
					break
				}

				reverseNotation = append(reverseNotation, currentToken)
			}

			continue
		}

		if token.isOperator() {
			for !operatorsStack.isEmpty() {
				topToken := operatorsStack.peek()

				if !topToken.isOperator() {
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

func constructBinaryTree(reverseNotation []*Token) *Node {
	stack := newNodeStack(len(reverseNotation))

	for _, token := range reverseNotation {
		if !token.isOperator() {
			stack.push(&Node{Key: token})
			continue
		}

		right := stack.pop()

		// NOT nodes have only right leaf
		if token.Key == NOT {
			stack.push(&Node{Key: token, Right: right})
			continue
		}

		left := stack.pop()

		stack.push(&Node{
			Key:   token,
			Left:  left,
			Right: right,
		})
	}

	return stack.pop()
}

// Parse parses query into the binary tree
func Parse(query string) *Node {
	return constructBinaryTree(parseSearchQuery(tokenize(query)))
}
