// Placed here due to the lack of generics

package textquery

type tokenStack []*Token

func newTokenStack(capacity int) tokenStack {
	return make(tokenStack, 0, capacity)
}

func (stack *tokenStack) isEmpty() bool {
	return len(*stack) == 0
}

func (stack *tokenStack) push(token *Token) {
	*stack = append(*stack, token)
}

func (stack *tokenStack) pop() *Token {
	index := len(*stack) - 1
	token := (*stack)[index]
	*stack = (*stack)[:index]
	return token
}
func (stack *tokenStack) peek() *Token {
	index := len(*stack) - 1
	return (*stack)[index]
}

type nodeStack []*Node

func newNodeStack(capacity int) nodeStack {
	return make(nodeStack, 0, capacity)
}

func (stack *nodeStack) push(node *Node) {
	*stack = append(*stack, node)
}

func (stack *nodeStack) pop() *Node {
	index := len(*stack) - 1
	node := (*stack)[index]
	*stack = (*stack)[:index]
	return node
}
