package apt

import (
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
	SetParent(parent Node)
	GetParent() Node
	GetChildren() []Node
	AddRandom(node Node)
	AddLeaf(leaf Node) bool
	NodeCount() int
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

type BaseNode struct {
	Parent   Node
	Children []Node
}

func (b *BaseNode) GetChildren() []Node {
	return b.Children
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

func (b *BaseNode) Evaluate(x, y float64) float64 {
	panic("should call Evaluate() on a BaseNode")
}

func (b *BaseNode) String() string {
	panic("should call String() on a BaseNode")
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

type OperatorPlus struct {
	BaseNode
}

func NewPlus() *OperatorPlus {
	return &OperatorPlus{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorPlus) Evaluate(x, y float64) float64 {
	return op.Children[0].Evaluate(x, y) + op.Children[1].Evaluate(x, y)
}

func (op *OperatorPlus) String() string {
	return "( + " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorMinus struct {
	BaseNode
}

func NewMinus() *OperatorMinus {
	return &OperatorMinus{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorMinus) Evaluate(x, y float64) float64 {
	return op.Children[0].Evaluate(x, y) - op.Children[1].Evaluate(x, y)
}

func (op *OperatorMinus) String() string {
	return "( - " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorMult struct {
	BaseNode
}

func NewMult() *OperatorMult {
	return &OperatorMult{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorMult) Evaluate(x, y float64) float64 {
	return op.Children[0].Evaluate(x, y) * op.Children[1].Evaluate(x, y)
}

func (op *OperatorMult) String() string {
	return "( * " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorDiv struct {
	BaseNode
}

func NewDiv() *OperatorDiv {
	return &OperatorDiv{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorDiv) Evaluate(x, y float64) float64 {
	return op.Children[0].Evaluate(x, y) / op.Children[1].Evaluate(x, y)
}

func (op *OperatorDiv) String() string {
	return "( / " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorAtan2 struct {
	BaseNode
}

func NewAtan2() *OperatorAtan2 {
	return &OperatorAtan2{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorAtan2) Evaluate(x, y float64) float64 {
	return math.Atan2(y, x)
}

func (op *OperatorAtan2) String() string {
	return "( Atan2 " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorNoise struct {
	BaseNode
}

func NewNoise() *OperatorNoise {
	return &OperatorNoise{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorNoise) Evaluate(x, y float64) float64 {
	return 80*noise.Snoise2(op.Children[0].Evaluate(x, y), op.Children[1].Evaluate(x, y)) - 2.0
}

func (op *OperatorNoise) String() string {
	return "( SimplexNoise " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorSquare struct {
	BaseNode
}

func NewSquare() *OperatorSquare {
	return &OperatorSquare{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorSquare) Evaluate(x, y float64) float64 {
	value := op.Children[0].Evaluate(x, y)
	return value * value
}

func (op *OperatorSquare) String() string {
	return "( Square " + op.Children[0].String() + " )"
}

type OperatorLog2 struct {
	BaseNode
}

func NewLog2() *OperatorLog2 {
	return &OperatorLog2{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorLog2) Evaluate(x, y float64) float64 {
	return math.Log2(op.Children[0].Evaluate(x, y))
}

func (op *OperatorLog2) String() string {
	return "( Log2 " + op.Children[0].String() + " )"
}

type OperatorNegate struct {
	BaseNode
}

func NewNegate() *OperatorNegate {
	return &OperatorNegate{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorNegate) Evaluate(x, y float64) float64 {
	return -op.Children[0].Evaluate(x, y)
}

func (op *OperatorNegate) String() string {
	return "( Negate " + op.Children[0].String() + " )"
}

type OperatorCeil struct {
	BaseNode
}

func NewCeil() *OperatorCeil {
	return &OperatorCeil{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorCeil) Evaluate(x, y float64) float64 {
	return math.Ceil(op.Children[0].Evaluate(x, y))
}

func (op *OperatorCeil) String() string {
	return "( Ceil " + op.Children[0].String() + " )"
}

type OperatorFloor struct {
	BaseNode
}

func NewFloor() *OperatorFloor {
	return &OperatorFloor{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorFloor) Evaluate(x, y float64) float64 {
	return math.Floor(op.Children[0].Evaluate(x, y))
}

func (op *OperatorFloor) String() string {
	return "( Floor " + op.Children[0].String() + " )"
}

type OperatorAbs struct {
	BaseNode
}

func NewAbs() *OperatorAbs {
	return &OperatorAbs{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorAbs) Evaluate(x, y float64) float64 {
	return math.Abs(op.Children[0].Evaluate(x, y))
}

func (op *OperatorAbs) String() string {
	return "( Abs " + op.Children[0].String() + " )"
}

type OperatorClip struct {
	BaseNode
}

func NewClip() *OperatorClip {
	return &OperatorClip{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 2),
		},
	}
}

func (op *OperatorClip) Evaluate(x, y float64) float64 {
	value := op.Children[0].Evaluate(x, y)

	maxVal := math.Abs(op.Children[1].Evaluate(x, y))
	if value > maxVal {
		return maxVal
	} else if value < -maxVal {
		return -maxVal
	}
	return value
}

func (op *OperatorClip) String() string {
	return "( Clip " + op.Children[0].String() + " " + op.Children[1].String() + " )"
}

type OperatorWrap struct {
	BaseNode
}

func NewWrap() *OperatorWrap {
	return &OperatorWrap{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorWrap) Evaluate(x, y float64) float64 {
	f := op.Children[0].Evaluate(x, y)
	temp := (f - -1.0) / (2.0)
	return -1.0 + 2.0*(temp-math.Floor(temp))
}

func (op *OperatorWrap) String() string {
	return "(Wrap " + op.Children[0].String() + ")"
}

type OperatorFBM struct {
	BaseNode
}

func NewFBM() *OperatorFBM {
	return &OperatorFBM{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 3),
		},
	}
}

func (op *OperatorFBM) Evaluate(x, y float64) float64 {
	return 2*3.627*noise.Fbm2(op.Children[0].Evaluate(x, y), op.Children[1].Evaluate(x, y), 5*op.Children[2].Evaluate(x,
		y), 0.5, 2, 3) + .492 - 1
}

func (op *OperatorFBM) String() string {
	return "( FBM " + op.Children[0].String() + " " + op.Children[1].String() + " " + op.Children[2].String() + " )"
}

type OperatorTurbulence struct {
	BaseNode
}

func NewTurbulence() *OperatorTurbulence {
	return &OperatorTurbulence{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 3),
		},
	}
}

func (op *OperatorTurbulence) Evaluate(x, y float64) float64 {
	return 2*6.96*noise.Turbulence(op.Children[0].Evaluate(x, y), op.Children[1].Evaluate(x, y),
		5*op.Children[2].Evaluate(x, y), 0.5, 2, 3) - 1
}

func (op *OperatorTurbulence) String() string {
	return "( Turbulence " + op.Children[0].String() + " " + op.Children[1].String() + " " + op.Children[2].String() + " )"
}

type OperatorLerp struct {
	BaseNode
}

func NewLerp() *OperatorLerp {
	return &OperatorLerp{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 3),
		},
	}
}

func (op *OperatorLerp) Evaluate(x, y float64) float64 {
	a := op.Children[0].Evaluate(x, y)
	b := op.Children[1].Evaluate(x, y)
	pct := op.Children[2].Evaluate(x, y)
	return a + pct*(b-a)
}

func (op *OperatorLerp) String() string {
	return "( Lerp " + op.Children[0].String() + " " + op.Children[1].String() + " " + op.Children[2].String() + " )"
}

type OperatorSin struct {
	BaseNode
}

func NewSin() *OperatorSin {
	return &OperatorSin{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorSin) Evaluate(x, y float64) float64 {
	return math.Sin(op.Children[0].Evaluate(x, y))
}

func (op *OperatorSin) String() string {
	return "( Sin " + op.Children[0].String() + " )"
}

type OperatorCos struct {
	BaseNode
}

func NewCos() *OperatorCos {
	return &OperatorCos{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorCos) Evaluate(x, y float64) float64 {
	return math.Cos(op.Children[0].Evaluate(x, y))
}

func (op *OperatorCos) String() string {
	return "( Cos " + op.Children[0].String() + " )"
}

type OperatorAtan struct {
	BaseNode
}

func NewAtan() *OperatorAtan {
	return &OperatorAtan{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 1),
		},
	}
}

func (op *OperatorAtan) Evaluate(x, y float64) float64 {
	return math.Atan(op.Children[0].Evaluate(x, y))
}

func (op *OperatorAtan) String() string {
	return "( Atan " + op.Children[0].String() + " )"
}

type OperatorX struct {
	BaseNode
}

func NewX() *OperatorX {
	return &OperatorX{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 0),
		},
	}
}

func (op *OperatorX) Evaluate(x, _ float64) float64 {
	return x
}

func (op *OperatorX) String() string {
	return "x"
}

type OperatorY struct {
	BaseNode
}

func NewY() *OperatorY {
	return &OperatorY{
		BaseNode{
			Parent:   nil,
			Children: make([]Node, 0),
		},
	}
}

func (op *OperatorY) Evaluate(_, y float64) float64 {
	return y
}

func (op *OperatorY) String() string {
	return "y"
}

type OperatorConstant struct {
	BaseNode
	Value float64
}

func NewConstant() *OperatorConstant {
	return &OperatorConstant{
		BaseNode: BaseNode{
			Parent:   nil,
			Children: make([]Node, 0),
		},
		Value: rand.Float64()*2 - 1,
	}
}

func (op *OperatorConstant) Evaluate(_, _ float64) float64 {
	return op.Value
}

func (op *OperatorConstant) String() string {
	return strconv.FormatFloat(op.Value, 'f', 9, 64)
}

func GetRandomNode() Node {
	r := rand.Intn(19)
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
