package apt

import (
	"math/rand"
	"reflect"
)

// Leaf node (0 children) Example : constants,x,y,t
// Single nodes (1 child) Example : sin/cos
// Double nodes (2 children) Example : +-*/
// Triple nodes (3 children) Example : fmb, turbulence

type Node interface {
	Evaluate(x, y float64) float64
	String() string
	SetParent(parent Node)
	GetParent() Node
	SetChildren([]Node)
	GetChildren() []Node
	AddRandom(node Node)
	AddLeaf(leaf Node) bool
	NodeCount() int
}

// CopyTree copies a tree (or subtree) and returns the copy.
// Since Node is an interface, we need to use reflection in CopyTree
func CopyTree(node, parent Node) Node {
	nodeCopy := reflect.New(reflect.ValueOf(node).Elem().Type()).Interface().(Node)

	// Make sure that constants have their value copied
	switch n := node.(type) {
	case *OperatorConstant:
		nodeCopy.(*OperatorConstant).Value = n.Value
	}

	nodeCopy.SetParent(parent)
	copyChildren := make([]Node, len(node.GetChildren()))
	nodeCopy.SetChildren(copyChildren)

	for i := range copyChildren {
		copyChildren[i] = CopyTree(node.GetChildren()[i], nodeCopy)
	}

	return nodeCopy
}

func ReplaceNode(old, new Node) {
	oldParent := old.GetParent()

	if oldParent != nil {
		for i, node := range oldParent.GetChildren() {
			if node == old {
				oldParent.GetChildren()[i] = new
			}
		}
	}

	new.SetParent(oldParent)
}

func GetNthNode(node Node, n, count int) (Node, int) {
	if n == count {
		return node, count
	}
	var result Node
	for _, child := range node.GetChildren() {
		count++
		result, count = GetNthNode(child, n, count)
		if result != nil {
			return result, count
		}
	}

	return nil, count
}

func Mutate(node Node) Node {
	r := rand.Intn(23)
	var mutatedNode Node

	if r <= 19 {
		mutatedNode = GetRandomNode()
	} else {
		mutatedNode = GetRandomLeafNode()
	}

	// Fix parents child pointer
	if node.GetParent() != nil {
		for i, parentChild := range node.GetParent().GetChildren() {
			if parentChild == node {
				node.GetParent().GetChildren()[i] = mutatedNode
			}
		}
	}

	// Add children from the old node to the mutated node
	for i, child := range node.GetChildren() {
		if i >= len(mutatedNode.GetChildren()) {
			break
		}
		mutatedNode.GetChildren()[i] = child
		child.SetParent(mutatedNode)
	}

	// Any nil children are filled with random leafs
	for i, child := range mutatedNode.GetChildren() {
		if child == nil {
			leaf := GetRandomLeafNode()
			leaf.SetParent(mutatedNode)
			mutatedNode.GetChildren()[i] = leaf
		}
	}

	mutatedNode.SetParent(node.GetParent())

	return mutatedNode
}

func GetRandomNode() Node {
	r := rand.Intn(21)
	switch r {
	case 0:
		return NewPlus()
	case 1:
		return NewMinus()
	case 2:
		return NewMult()
	case 3:
		return NewDiv()
	case 4:
		return NewAtan2()
	case 5:
		return NewAtan()
	case 6:
		return NewCos()
	case 7:
		return NewSin()
	case 8:
		return NewNoise()
	case 9:
		return NewSquare()
	case 10:
		return NewLog2()
	case 11:
		return NewNegate()
	case 12:
		return NewCeil()
	case 13:
		return NewFloor()
	case 14:
		return NewAbs()
	case 15:
		return NewClip()
	case 16:
		return NewWrap()
	case 17:
		return NewLerp()
	case 18:
		return NewFBM()
	case 19:
		return NewTurbulence()
	case 20:
		return NewSwirl()

	default:
		panic("GetRandomNode failed")
	}
}

func GetRandomLeafNode() Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return NewX()
	case 1:
		return NewY()
	case 2:
		return NewConstant()
	default:
		panic("GetRandomBaseNode failed")
	}
}
