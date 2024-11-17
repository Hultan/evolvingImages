package main

import (
	"fmt"
	"math/rand"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	. "github.com/hultan/evolvingImage/apt"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var texture rl.Texture2D
var index int32
var imageData = make([]byte, screenWidth*screenHeight*4)
var image = rl.NewImage(imageData, screenWidth, screenHeight, 1, rl.UncompressedR8g8b8a8)

type picture struct {
	r, g, b Node
}

func NewPicture() *picture {
	p := &picture{}

	// Generate image
	p.r = GetRandomNode()
	p.g = GetRandomNode()
	p.b = GetRandomNode()

	const nodes = 4

	num := rand.Intn(nodes)
	for i := 0; i < num; i++ {
		p.r.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes)
	for i := 0; i < num; i++ {
		p.g.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes)
	for i := 0; i < num; i++ {
		p.b.AddRandom(GetRandomNode())
	}

	for p.r.AddLeaf(GetRandomLeafNode()) {
	}
	for p.b.AddLeaf(GetRandomLeafNode()) {
	}
	for p.g.AddLeaf(GetRandomLeafNode()) {
	}

	return p
}

func (p *picture) String() string {
	return "R:" + p.r.String() + "\nG:" + p.g.String() + "\nB:" + p.b.String()
}

func (p *picture) Mutate() {
	r := rand.Intn(3)
	var nodeToMutate Node

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
	nodeToMutate, count = GetNthNode(nodeToMutate, r, 0)
	// If the node that we mutated is one of the root nodes
	// we need to handle that.
	mutation := Mutate(nodeToMutate)
	if mutation == p.r {
		p.r = mutation
	} else if mutation == p.g {
		p.g = mutation
	} else if mutation == p.b {
		p.b = mutation
	}
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	p := NewPicture()
	fmt.Println(p.String())
	fmt.Println()

	generateImage(p, screenWidth, screenHeight)

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// Update
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			for i := 0; i < 5; i++ {
				p.Mutate()
			}
			fmt.Println(p.String())
			fmt.Println()
			generateImage(p, screenWidth, screenHeight)
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.DrawRectangle(5, 5, 120, 30, rl.Black)
		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

	// Clean up
	imageData = nil
	//rl.UnloadImage(image)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

func generateImage(p *picture, width, height int) {
	scale := 128.0
	offset := -1 * scale
	index = 0

	for y := 0; y < height; y++ {
		yy := float64(y)/float64(height)*2 - 1
		for x := 0; x < width; x++ {
			xx := float64(x)/float64(width)*2 - 1
			r := p.r.Evaluate(xx, yy)
			g := p.g.Evaluate(xx, yy)
			b := p.b.Evaluate(xx, yy)

			imageData[index+0] = byte(r*scale - offset)
			imageData[index+1] = byte(g*scale - offset)
			imageData[index+2] = byte(b*scale - offset)
			imageData[index+3] = 255
			index += 4
		}
	}

	image.Data = unsafe.Pointer(unsafe.SliceData(imageData))
	texture = rl.LoadTextureFromImage(image)
}
