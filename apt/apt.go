package apt

import (
	"math"
	"strconv"
	"time"

	simplex "github.com/ojrac/opensimplex-go"
)

// Leaf node (0 children) Example : constants,x,y,t
// Single nodes (1 child) Example : sin/cos
// Double nodes (2 children) Example : +-*/

var noise = simplex.NewNormalized(time.Now().UTC().UnixNano())

type Node interface {
	Evaluate(x, y float64) float64
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
	return noise.Eval2(op.LeftChild.Evaluate(x, y), op.RightChild.Evaluate(x, y))*2 - 1
}

func (op *OperatorNoise) String() string {
	return "( SimplexNoise " + op.LeftChild.String() + " " + op.RightChild.String() + " )"
}

type OperatorSin SingleNode

func (op *OperatorSin) Evaluate(x, y float64) float64 {
	return math.Sin(op.Child.Evaluate(x, y))
}

func (op *OperatorSin) String() string {
	return "( Sin " + op.Child.String() + " )"
}

type OperatorCos SingleNode

func (op *OperatorCos) Evaluate(x, y float64) float64 {
	return math.Cos(op.Child.Evaluate(x, y))
}

func (op *OperatorCos) String() string {
	return "( Cos " + op.Child.String() + " )"
}

type OperatorAtan SingleNode

func (op *OperatorAtan) Evaluate(x, y float64) float64 {
	return math.Atan(op.Child.Evaluate(x, y))
}

func (op *OperatorAtan) String() string {
	return "( Atan " + op.Child.String() + " )"
}

type OperatorX LeafNode

func (op *OperatorX) Evaluate(x, _ float64) float64 {
	return x
}

func (op *OperatorX) String() string {
	return "x"
}

type OperatorY LeafNode

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
