package apt

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"

	"github.com/hultan/evolvingImage/noise"
)

// Leaf node (0 children) Example : constants,x,y,t
// Single nodes (1 child) Example : sin/cos
// Double nodes (2 children) Example : +-*/

type Node interface {
	Evaluate(x, y float64) float64
	String() string
	AddRandom(node Node)
	NodeCounts() (nodeCount, nilCount int)
}

type LeafNode struct{}

func (l *LeafNode) AddRandom(n Node) {
	//panic("Error : You tried to add a node to a leaf node!")
	fmt.Println("Bug:Added a node to a leaf node")
}

func (l *LeafNode) NodeCounts() (int, int) {
	return 1, 0
}

type SingleNode struct {
	Child Node
}

func (s *SingleNode) AddRandom(n Node) {
	if s.Child == nil {
		s.Child = n
	} else {
		s.Child.AddRandom(n)
	}
}

func (s *SingleNode) NodeCounts() (int, int) {
	if s.Child == nil {
		return 1, 1
	}
	nodeCount, nilCount := s.Child.NodeCounts()

	return nodeCount + 1, nilCount
}

type DoubleNode struct {
	LeftChild  Node
	RightChild Node
}

func (d *DoubleNode) AddRandom(n Node) {
	r := rand.Intn(2)
	if r == 0 {
		if d.LeftChild == nil {
			d.LeftChild = n
		} else {
			d.LeftChild.AddRandom(n)
		}
	} else {
		if d.RightChild == nil {
			d.RightChild = n
		} else {
			d.RightChild.AddRandom(n)
		}
	}
}

func (d *DoubleNode) NodeCounts() (int, int) {
	var leftNodeCount, leftNilCount, rightNodeCount, rightNilCount int

	if d.LeftChild == nil {
		leftNodeCount = 0
		leftNilCount = 1
	} else {
		leftNodeCount, leftNilCount = d.LeftChild.NodeCounts()
	}

	if d.RightChild == nil {
		rightNodeCount = 0
		rightNilCount = 1
	} else {
		rightNodeCount, rightNilCount = d.RightChild.NodeCounts()
	}

	return leftNodeCount + rightNodeCount + 1, leftNilCount + rightNilCount
}

type OperatorPlus struct {
	DoubleNode
}

func (op *OperatorPlus) Evaluate(x, y float64) float64 {
	return op.LeftChild.Evaluate(x, y) + op.RightChild.Evaluate(x, y)
}

func (op *OperatorPlus) String() string {
	return "( + " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorMinus struct {
	DoubleNode
}

func (op *OperatorMinus) Evaluate(x, y float64) float64 {
	return op.LeftChild.Evaluate(x, y) - op.RightChild.Evaluate(x, y)
}

func (op *OperatorMinus) String() string {
	return "( - " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorMult struct {
	DoubleNode
}

func (op *OperatorMult) Evaluate(x, y float64) float64 {
	return op.LeftChild.Evaluate(x, y) * op.RightChild.Evaluate(x, y)
}

func (op *OperatorMult) String() string {
	return "( * " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorDiv struct {
	DoubleNode
}

func (op *OperatorDiv) Evaluate(x, y float64) float64 {
	return op.LeftChild.Evaluate(x, y) / op.RightChild.Evaluate(x, y)
}

func (op *OperatorDiv) String() string {
	return "( / " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorAtan2 struct {
	DoubleNode
}

func (op *OperatorAtan2) Evaluate(x, y float64) float64 {
	return math.Atan2(y, x)
}

func (op *OperatorAtan2) String() string {
	return "( Atan2 " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorNoise struct {
	DoubleNode
}

func (op *OperatorNoise) Evaluate(x, y float64) float64 {
	return 80*noise.Snoise2(op.LeftChild.Evaluate(x, y), op.RightChild.Evaluate(x, y)) - 2.0
}

func (op *OperatorNoise) String() string {
	return "( SimplexNoise " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorSin struct {
	SingleNode
}

func (op *OperatorSin) Evaluate(x, y float64) float64 {
	return math.Sin(op.Child.Evaluate(x, y))
}

func (op *OperatorSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OperatorCos struct {
	SingleNode
}

func (op *OperatorCos) Evaluate(x, y float64) float64 {
	return math.Cos(op.Child.Evaluate(x, y))
}

func (op *OperatorCos) String() string {
	return "( Cos " + op.Child.String() + " )"
}

type OperatorAtan struct {
	SingleNode
}

func (op *OperatorAtan) Evaluate(x, y float64) float64 {
	return math.Atan(op.Child.Evaluate(x, y))
}

func (op *OperatorAtan) String() string {
	return "( Atan " + op.Child.String() + " )"
}

type OperatorX struct {
	LeafNode
}

func (op *OperatorX) Evaluate(x, _ float64) float64 {
	return x
}

func (op *OperatorX) String() string {
	return "x"
}

type OperatorY struct {
	LeafNode
}

func (op *OperatorY) Evaluate(_, y float64) float64 {
	return y
}

func (op *OperatorY) String() string {
	return "y"
}

type OperatorConstant struct {
	LeafNode
	Value float64
}

func (op *OperatorConstant) Evaluate(_, _ float64) float64 {
	return op.Value
}

func (op *OperatorConstant) String() string {
	return strconv.FormatFloat(op.Value, 'f', 9, 64)
}

func GetRandomNode() Node {
	r := rand.Intn(9)
	switch r {
	case 0:
		return &OperatorPlus{}
	case 1:
		return &OperatorMinus{}
	case 2:
		return &OperatorMult{}
	case 3:
		return &OperatorDiv{}
	case 4:
		return &OperatorAtan2{}
	case 5:
		return &OperatorAtan{}
	case 6:
		return &OperatorCos{}
	case 7:
		return &OperatorSin{}
	case 8:
		return &OperatorNoise{}
	default:
		panic("GetRandomNode failed")
	}
}

func GetRandomLeafNode() Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return &OperatorX{}
	case 1:
		return &OperatorY{}
	case 2:
		return &OperatorConstant{
			LeafNode: LeafNode{},
			Value:    rand.Float64()*2 - 1,
		}
	default:
		panic("GetRandomLeafNode failed")
	}
}
