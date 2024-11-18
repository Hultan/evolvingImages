package main

import (
	"math/rand"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	. "github.com/hultan/evolvingImage/apt"
)

var screenWidth, screenHeight int32 = 800, 600
var rows, cols, numPics int32 = 5, 5, rows * cols

type imageResult struct {
	image *rl.Image
	index int32
}

type picture struct {
	r, g, b Node
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	pictures := make([]*picture, numPics)
	for i := range pictures {
		pictures[i] = NewPicture()
	}

	picWidth := int32(float32(screenWidth/cols) * 0.9)
	picHeight := int32(float32(screenHeight/rows) * 0.9)

	buttons := make([]*Button, numPics)
	imageChannel := make(chan imageResult, numPics)
	for i := range buttons {
		go func(i int) {
			image := generateImage(pictures[i], picWidth, picHeight)
			imageChannel <- imageResult{
				image,
				int32(i),
			}
		}(i)
	}

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

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
					buttons[img.index] = NewButton(rec, rl.LoadTextureFromImage(img.image))
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

		rl.DrawRectangle(5, 5, 120, 30, rl.Black)
		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}

	// Clean up
	for i := range buttons {
		rl.UnloadTexture(buttons[i].Texture)
	}

	rl.CloseWindow()
}

func generateImage(p *picture, width, height int32) *rl.Image {
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

func NewPicture() *picture {
	p := &picture{}

	// Generate image
	p.r = GetRandomNode()
	p.g = GetRandomNode()
	p.b = GetRandomNode()

	const nodes = 20

	num := rand.Intn(nodes) + 5
	for i := 0; i < num; i++ {
		p.r.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes) + 5
	for i := 0; i < num; i++ {
		p.g.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes) + 5
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
