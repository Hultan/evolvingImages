package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/hultan/evolvingImage/apt"
)

const nodes = imageComplexity

type Picture struct {
	r, g, b apt.Node
}

func newPicture() *Picture {
	p := &Picture{}

	// Generate image
	p.r = p.newNode()
	p.g = p.newNode()
	p.b = p.newNode()

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
	return "( Picture \n" + p.r.String() + " \n" + p.g.String() + " \n" + p.b.String() + " \n)"
}

func (p *Picture) mutate() {
	r := rand.Intn(3)
	var nodeToMutate apt.Node

	switch r {
	case 0:
		nodeToMutate = p.r
	case 1:
		nodeToMutate = p.g
	case 2:
		nodeToMutate = p.b
	default:
		panic("should not happen")
	}

	count := nodeToMutate.NodeCount()
	r = rand.Intn(count)
	nodeToMutate, count = apt.GetNthNode(nodeToMutate, r, 0)
	// If the node that we mutated is one of the root nodes
	// we need to handle that.
	mutation := apt.Mutate(nodeToMutate)
	if mutation == p.r {
		p.r = mutation
	} else if mutation == p.g {
		p.g = mutation
	} else if mutation == p.b {
		p.b = mutation
	}
}

func (p *Picture) save() {
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

func (p *Picture) cross(other *Picture) *Picture {
	aCopy := &Picture{
		apt.CopyTree(p.r, nil),
		apt.CopyTree(p.g, nil),
		apt.CopyTree(p.b, nil),
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
		return p.r
	case 1:
		return p.g
	case 2:
		return p.b
	default:
		panic("PickRandomColor : Should not happen!")
	}
}
