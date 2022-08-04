package textquery

import "strings"

type associativity int8

const (
	right associativity = 1
	left  associativity = 2
)

// Node is a binary tree node
type Node struct {
	Key   *Token
	Left  *Node
	Right *Node
}

// Match matches a string against query.
func (n *Node) Match(s string) bool {
	if s == "" || n == nil {
		return false
	}

	l := n.Left.Match(s)
	r := n.Right.Match(s)

	// NOT corner case
	if n.Right != nil &&
		n.Right.Key != nil &&
		n.Right.Key.isOperator() &&
		n.Right.Key.Data == NOT {
		r = !n.Right.Right.Match(s)
	}

	// Operator
	switch n.Key.Data {
	case AND:
		return l && r
	case OR:
		return l || r
	}

	// Operand
	return strings.Contains(s, n.Key.Data)
}

type precedence struct {
	precedence    int
	associativity associativity
}

var precedences = map[string]precedence{
	NOT: {precedence: 3, associativity: right},
	AND: {precedence: 2, associativity: left},
	OR:  {precedence: 1, associativity: left},
}

func hasPrecedence(a, b *Token) bool {
	opA := precedences[a.Data]
	opB := precedences[b.Data]

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
		if token.Data == NOT {
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
	if query == "" {
		return nil
	}
	return constructBinaryTree(parseSearchQuery(tokenize(query)))
}
