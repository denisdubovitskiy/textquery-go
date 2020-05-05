package query

type nodeStack []*Node

func newNodeStack(capacity int) nodeStack {
	return make(nodeStack, 0, capacity)
}

func (s *nodeStack) isEmpty() bool {
	return len(*s) == 0
}

func (s *nodeStack) push(str *Node) {
	*s = append(*s, str)
}

func (s *nodeStack) pop() (*Node, bool) {
	if s.isEmpty() {
		return nil, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}

type parenthesis struct {
	position  int
	character string
}

type parenthesisStack []parenthesis

func (s *parenthesisStack) isEmpty() bool {
	return len(*s) == 0
}

func (s *parenthesisStack) push(str parenthesis) {
	*s = append(*s, str)
}

func (s *parenthesisStack) pop() (parenthesis, bool) {
	if s.isEmpty() {
		return parenthesis{}, false
	} else {
		index := len(*s) - 1
		element := (*s)[index]
		*s = (*s)[:index]
		return element, true
	}
}
