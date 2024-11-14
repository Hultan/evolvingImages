package apt

import "math"

// Leaf node (0 children) Example : constants,x,y,t
// Single nodes (1 child) Example : sin/cos
// Double nodes (2 children) Example : +-*/

type Node interface {
	Evaluate(x, y, t float64) float64
	String() string
}

type LeafNode struct{}

type SingleNode struct {
	Child Node
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

type OperatorPlus struct {
	DoubleNode
}

func (op *OperatorPlus) Evaluate(x, y, t float64) float64 {
	return op.LeftChild.Evaluate(x, y, t) + op.RightChild.Evaluate(x, y, t)
}

func (op *OperatorPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorSin SingleNode

func (op *OperatorSin) Evaluate(x, y, t float64) float64 {
	return math.Sin(op.Child.Evaluate(x, y, t))
}

func (op *OperatorSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OperatorX LeafNode

func (op *OperatorX) Evaluate(x, _, _ float64) float64 {
	return x
}

func (op *OperatorX) String() string {
	return "x"
}

type OperatorY LeafNode

func (op *OperatorY) Evaluate(_, y, _ float64) float64 {
	return y
}

func (op *OperatorY) String() string {
	return "y"
}

type OperatorT LeafNode

func (op *OperatorT) Evaluate(_, _, t float64) float64 {
	return t
}

func (op *OperatorT) String() string {
	return "t"
}
