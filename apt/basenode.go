package apt

import "math/rand"

type BaseNode struct {
	Parent   Node
	Children []Node
}

func (b *BaseNode) GetChildren() []Node {
	return b.Children
}

func (b *BaseNode) SetChildren(children []Node) {
	b.Children = children
}

func (b *BaseNode) NodeCount() int {
	count := 1
	for _, child := range b.Children {
		count += child.NodeCount()
	}
	return count
}

func (b *BaseNode) GetParent() Node {
	return b.Parent
}

func (b *BaseNode) SetParent(parent Node) {
	b.Parent = parent
}

func (b *BaseNode) Evaluate(_, _ float64) float64 {
	panic("do not call Evaluate() on a BaseNode")
}

func (b *BaseNode) String() string {
	panic("do not call String() on a BaseNode")
}

func (b *BaseNode) AddRandom(node Node) {
	addIndex := rand.Intn(len(b.Children))
	if b.Children[addIndex] == nil {
		node.SetParent(b)
		b.Children[addIndex] = node
	} else {
		b.Children[addIndex].AddRandom(node)
	}
}

func (b *BaseNode) AddLeaf(leaf Node) bool {
	for i, child := range b.Children {
		if child == nil {
			leaf.SetParent(b)
			b.Children[i] = leaf
			return true
		} else if b.Children[i].AddLeaf(leaf) {
			return true
		}
	}
	return false
}
