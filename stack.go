package textquery

type tokenStack []Token

func newTokenStack(capacity int) tokenStack {
	return make(tokenStack, 0, capacity)
}

func (s *tokenStack) isEmpty() bool {
	return len(*s) == 0
}

func (s *tokenStack) push(str Token) {
	*s = append(*s, str)
}

func (s *tokenStack) pop() Token {
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}
func (s *tokenStack) peek() Token {
	index := len(*s) - 1
	element := (*s)[index]
	return element
}

func (s *tokenStack) size() int {
	return len(*s)
}

type nodeStack []Node

func newNodeStack(capacity int) nodeStack {
	return make(nodeStack, 0, capacity)
}

func (s *nodeStack) isEmpty() bool {
	return len(*s) == 0
}

func (s *nodeStack) push(str Node) {
	*s = append(*s, str)
}

func (s *nodeStack) pop() Node {
	index := len(*s) - 1
	element := (*s)[index]
	*s = (*s)[:index]
	return element
}
func (s *nodeStack) peek() Node {
	index := len(*s) - 1
	element := (*s)[index]
	return element
}
