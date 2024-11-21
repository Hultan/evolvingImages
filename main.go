package main

import (
	"math/rand"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	. "github.com/hultan/evolvingImage/apt"
)

const (
	imageComplexity    = 25
	imageMinComplexity = 5
	mutationRate       = 10
)

var screenWidth, screenHeight int32 = 1600, 900
var rows, cols, numPics int32 = 5, 5, rows * cols
var picWidth, picHeight int32
var imageChannel chan ImageResult
var buttons []*Button
var pictures []*Picture

type ImageResult struct {
	Image *rl.Image
	index int32
}

type Picture struct {
	r, g, b Node
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	picWidth = int32(float32(screenWidth/cols) * 0.9)
	picHeight = int32(float32(screenHeight/rows) * 0.8)

	pictures = make([]*Picture, numPics)
	for i := range pictures {
		pictures[i] = NewPicture()
	}

	buttons = make([]*Button, numPics)

	evolveRect := rl.Rectangle{
		X:      float32(screenWidth)/2 - float32(picWidth)/2,
		Y:      float32(screenHeight) * 0.9,
		Width:  float32(picWidth),
		Height: float32(screenHeight) * 0.08,
	}
	evolveButton := NewTextButton(evolveRect, "Evolve!", evolveButtonClicked)

	imageChannel = make(chan ImageResult, numPics)
	for i := range buttons {
		go func(i int) {
			image := generateImage(pictures[i], picWidth, picHeight)
			imageChannel <- ImageResult{
				image,
				int32(i),
			}
		}(i)
	}

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// Update
		//screenWidth = int32(rl.GetScreenWidth())
		//screenHeight = int32(rl.GetScreenHeight())
		//picWidth = int32(float32(screenWidth/cols) * 0.9)
		//picHeight = int32(float32(screenHeight/rows) * 0.8)

		evolveButton.Update()

		if rl.IsKeyPressed(rl.KeyF5) {

		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		evolveButton.Draw()

		select {
		case img, ok := <-imageChannel:
			if ok {
				// Calculate image x,y position (1-3,1-3)
				xi := img.index % cols
				yi := (img.index - xi) / cols
				// Calculate image screen x,y position (in pixels)
				x := xi * picWidth
				y := yi * picHeight
				// Calculate padding around images
				xPadding := int32(float32(screenWidth) * 0.1 / float32(cols+1))
				yPadding := int32(float32(screenHeight) * 0.1 / float32(rows+1))
				// Add padding to the screen position
				x += xPadding * (int32(xi) + 1)
				y += yPadding * (int32(yi) + 1)

				if buttons[img.index] == nil {
					rec := rl.Rectangle{
						X:      float32(x),
						Y:      float32(y),
						Width:  float32(picWidth),
						Height: float32(picHeight),
					}
					buttons[img.index] = NewButton(rec, rl.LoadTextureFromImage(img.Image))
					//buttons[img.index] = NewTextButton(rec, "Per", func() {
					//	fmt.Println("Button was clicked!")
					//})
				}

				buttons[img.index].Draw()
			}
		default:
			// Do nothing
		}

		// Draw textures at the correct position
		for _, button := range buttons {
			if button != nil {
				button.Update()
				button.Draw()
			}
		}

		rl.DrawFPS(25, screenHeight-40)
		rl.EndDrawing()
	}

	// Clean up
	for i := range buttons {
		rl.UnloadTexture(buttons[i].Texture)
	}

	rl.CloseWindow()
}

func evolveButtonClicked() {
	selectedPictures := make([]*Picture, 0)
	for i, button := range buttons {
		if button.Selected {
			selectedPictures = append(selectedPictures, pictures[i])
		}
	}

	if len(selectedPictures) != 0 {
		for i := range buttons {
			buttons[i] = nil
		}

		pictures = evolve(selectedPictures)
		for i := range pictures {
			go func(i int) {
				pixels := generateImage(pictures[i], picWidth, picHeight)
				imageChannel <- ImageResult{
					pixels,
					int32(i),
				}
			}(i)
		}
	}
}

func evolve(survivors []*Picture) []*Picture {
	newPics := make([]*Picture, numPics)
	i := 0
	for i < len(survivors) {
		a := survivors[i]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = cross(a, b)
		i++
	}

	for i < len(newPics) {
		a := survivors[rand.Intn(len(survivors))]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = cross(a, b)
		i++
	}

	for _, pic := range newPics {
		r := rand.Intn(mutationRate)
		for i := 0; i < r; i++ {
			pic.mutate()
		}
	}

	return newPics
}

func (p *Picture) pickRandomColor() Node {
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

func cross(a, b *Picture) *Picture {
	aCopy := &Picture{
		CopyTree(a.r, nil),
		CopyTree(a.g, nil),
		CopyTree(a.b, nil),
	}
	aColor := aCopy.pickRandomColor()
	bColor := b.pickRandomColor()

	aIndex := rand.Intn(aColor.NodeCount())
	aNode, _ := GetNthNode(aColor, aIndex, 0)

	bIndex := rand.Intn(bColor.NodeCount())
	bNode, _ := GetNthNode(bColor, bIndex, 0)
	bNodeCopy := CopyTree(bNode, bNode.GetParent())

	ReplaceNode(aNode, bNodeCopy)
	return aCopy
}

func generateImage(p *Picture, width, height int32) *rl.Image {
	scale := 128.0
	offset := -1 * scale
	index := 0
	var imageData = make([]byte, width*height*4)

	for y := int32(0); y < height; y++ {
		yy := float64(y)/float64(height)*2 - 1
		for x := int32(0); x < width; x++ {
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

	var image = rl.NewImage(imageData, width, height, 1, rl.UncompressedR8g8b8a8)
	image.Data = unsafe.Pointer(unsafe.SliceData(imageData))
	return image
}

func NewPicture() *Picture {
	p := &Picture{}

	// Generate image
	p.r = GetRandomNode()
	p.g = GetRandomNode()
	p.b = GetRandomNode()

	const nodes = imageComplexity

	num := rand.Intn(nodes) + imageMinComplexity
	for i := 0; i < num; i++ {
		p.r.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes) + imageMinComplexity
	for i := 0; i < num; i++ {
		p.g.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes) + imageMinComplexity
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

func (p *Picture) String() string {
	return "R:" + p.r.String() + "\nG:" + p.g.String() + "\nB:" + p.b.String()
}

func (p *Picture) mutate() {
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
