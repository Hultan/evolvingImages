package picture

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/hultan/evolvingImage/apt"
)

const (
	imageComplexity    = 25
	imageMinComplexity = 5
	nodes              = imageComplexity
)

type Picture struct {
	R, G, B apt.Node
}

func NewPicture() *Picture {
	p := &Picture{}

	// Generate image
	p.R = p.newNode()
	p.G = p.newNode()
	p.B = p.newNode()

	return p
}

func (p *Picture) newNode() apt.Node {
	// Generate image
	node := apt.GetRandomNode()

	num := rand.Intn(nodes) + imageMinComplexity
	for i := 0; i < num; i++ {
		node.AddRandom(apt.GetRandomNode())
	}

	for node.AddLeaf(apt.GetRandomLeafNode()) {
	}

	return node
}

func (p *Picture) String() string {
	return "( Picture \n" + p.R.String() + " \n" + p.G.String() + " \n" + p.B.String() + " \n)"
}

func (p *Picture) Mutate() {
	r := rand.Intn(3)
	var nodeToMutate apt.Node

	switch r {
	case 0:
		nodeToMutate = p.R
	case 1:
		nodeToMutate = p.G
	case 2:
		nodeToMutate = p.B
	default:
		panic("should not happen")
	}

	count := nodeToMutate.NodeCount()
	r = rand.Intn(count)
	nodeToMutate, count = apt.GetNthNode(nodeToMutate, r, 0)
	// If the node that we mutated is one of the root nodes
	// we need to handle that.
	mutation := apt.Mutate(nodeToMutate)
	if mutation == p.R {
		p.R = mutation
	} else if mutation == p.G {
		p.G = mutation
	} else if mutation == p.B {
		p.B = mutation
	}
}

func (p *Picture) Save() {
	files, err := os.ReadDir("./")
	if err != nil {
		panic(err)
	}

	biggest := 0
	for _, f := range files {
		name := f.Name()
		if strings.HasSuffix(name, ".apt") {
			numberString := strings.TrimSuffix(name, ".apt")
			num, err := strconv.Atoi(numberString)
			if err != nil {
				panic(err)
			}
			if num > biggest {
				biggest = num
			}
		}
	}
	name := fmt.Sprintf("%d.apt", biggest+1)
	file, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, p.String())
	if err != nil {
		panic(err)
	}
}

func (p *Picture) Cross(other *Picture) *Picture {
	aCopy := &Picture{
		apt.CopyTree(p.R, nil),
		apt.CopyTree(p.G, nil),
		apt.CopyTree(p.B, nil),
	}
	aColor := aCopy.pickRandomColor()
	bColor := other.pickRandomColor()

	aIndex := rand.Intn(aColor.NodeCount())
	aNode, _ := apt.GetNthNode(aColor, aIndex, 0)

	bIndex := rand.Intn(bColor.NodeCount())
	bNode, _ := apt.GetNthNode(bColor, bIndex, 0)
	bNodeCopy := apt.CopyTree(bNode, bNode.GetParent())

	apt.ReplaceNode(aNode, bNodeCopy)
	return aCopy
}

func (p *Picture) pickRandomColor() apt.Node {
	r := rand.Intn(3)
	switch r {
	case 0:
		return p.R
	case 1:
		return p.G
	case 2:
		return p.B
	default:
		panic("PickRandomColor : Should not happen!")
	}
}
