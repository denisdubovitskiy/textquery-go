package query

type NodeType int

type Node struct {
	Data     Part
	Children []*Node
}

func (n *Node) insertLeft() {
	n.Children = append([]*Node{{}}, n.Children...)
}

func (n *Node) accessLeft() *Node {
	return n.Children[0]
}

func (n *Node) insertRight() {
	n.Children = append(n.Children, &Node{})
}

func (n *Node) accessRight() *Node {
	return n.Children[len(n.Children)-1]
}
