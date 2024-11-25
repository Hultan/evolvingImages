package main

// TODO : pictures should be part of button?
// TODO : Make the zoomed in picture show a loading indicator
// TODO : (Impossible?) Instead of passing x and y for each pixel, pass a slice of all the arguments
// TODO : Make the String functions output valid go code and make a program that will execute it
// TODO : Do a grayscale picture, or an HSV picture, or a black and white image (<0.5)
// TODO :

import (
	"math/rand"
	"os"
	"time"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/hultan/evolvingImage/apt"
	"github.com/hultan/evolvingImage/picture"
)

const (
	mutationRate = 10
)

type stateType int

const (
	stateInit stateType = iota
	stateSelect
	stateZoom
)

var screenWidth, screenHeight int32 = 1600, 900
var rows, cols, numPics int32 = 5, 5, rows * cols
var picWidth, picHeight = int32(float32(screenWidth/cols) * 0.9), int32(float32(screenHeight/rows) * 0.8)
var imageChannel = make(chan ImageResult, numPics)
var buttons = make([]*Button, numPics)
var pictures = make([]*picture.Picture, numPics)
var state GuiState
var evolveButton *Button

type GuiState struct {
	zoom      stateType
	zoomedIn  time.Time
	zoomImage rl.Texture2D
	zoomTree  *picture.Picture
}

type ImageResult struct {
	Image *rl.Image
	index int32
}

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	state = GuiState{zoom: stateInit}

	// Handle parsing of an .apt file
	args := os.Args
	if len(args) > 1 {
		handleArgs(args[1])
	}

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
		// Update
		if rl.IsWindowResized() {
			onGenerateNewImages()
		}

		if state.zoom == stateInit {
			onGenerateNewImages()
			state.zoom = stateSelect
		}

		if evolveButton != nil {
			evolveButton.update()
		}

		if rl.IsKeyPressed(rl.KeyS) && state.zoom == stateZoom {
			state.zoomTree.Save()
		}

		if rl.IsKeyPressed(rl.KeyF5) {
			onGenerateNewImages()
		}

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		if state.zoom == stateZoom {
			if time.Since(state.zoomedIn).Seconds() > 1 && rl.IsMouseButtonPressed(rl.MouseButtonRight) {
				state.zoom = stateSelect
			}
			rl.DrawTexture(state.zoomImage, 0, 0, rl.White)
		} else if state.zoom == stateSelect {
			evolveButton.draw()

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

					rec := rl.Rectangle{
						X:      float32(x),
						Y:      float32(y),
						Width:  float32(picWidth),
						Height: float32(picHeight),
					}
					buttons[img.index] = newButton(img.index, rec, rl.LoadTextureFromImage(img.Image), onFullScreen)
				}
			default:
				// Do nothing
			}

			// Draw textures at the correct position
			for _, b := range buttons {
				if b != nil {
					b.update()
					b.draw()
				}
			}
		}

		x := screenWidth - 600
		rl.DrawText("Left mouse click : select an image.", x, screenHeight-80, 24, rl.LightGray)
		rl.DrawText("Right mouse click : zoom in/out.", x, screenHeight-50, 24, rl.LightGray)

		rl.DrawFPS(25, screenHeight-50)
		rl.EndDrawing()
	}

	// Clean up
	for i := range buttons {
		if buttons[i] != nil {
			rl.UnloadTexture(buttons[i].Texture)
		}
	}

	rl.CloseWindow()
}

func handleArgs(fileName string) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	str := string(bytes)
	pictureNode := apt.BeginLexing(str)
	p := &picture.Picture{
		R: pictureNode.GetChildren()[0],
		G: pictureNode.GetChildren()[1],
		B: pictureNode.GetChildren()[2],
	}
	zoomIn(p)
}

func zoomIn(p *picture.Picture) {
	zoomImage := newImage(p, screenWidth, int32(float32(screenHeight)*0.9))
	state.zoomImage = rl.LoadTextureFromImage(zoomImage)
	state.zoomTree = p
	state.zoom = stateZoom
	state.zoomedIn = time.Now()
}

func onGenerateNewImages() {
	screenWidth = int32(rl.GetScreenWidth())
	screenHeight = int32(rl.GetScreenHeight())
	picWidth = int32(float32(screenWidth/cols) * 0.9)
	picHeight = int32(float32(screenHeight/rows) * 0.8)
	for i := range pictures {
		pictures[i] = picture.NewPicture()
	}

	evolveRect := rl.Rectangle{
		X:      float32(screenWidth)/2 - float32(picWidth)/2,
		Y:      float32(screenHeight) * 0.9,
		Width:  float32(picWidth),
		Height: float32(screenHeight) * 0.08,
	}
	evolveButton = newTextButton(evolveRect, "Evolve!", onEvolveButtonClicked)

	for i := range buttons {
		go func(i int) {
			image := newImage(pictures[i], picWidth, picHeight)
			imageChannel <- ImageResult{
				image,
				int32(i),
			}
		}(i)
	}
}

func onFullScreen(button *Button) {
	if state.zoom == stateSelect {
		zoomIn(pictures[button.Index])
	}
}

func onEvolveButtonClicked() {
	selectedPictures := make([]*picture.Picture, 0)
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
				pixels := newImage(pictures[i], picWidth, picHeight)
				imageChannel <- ImageResult{
					pixels,
					int32(i),
				}
			}(i)
		}
	}
}

func evolve(survivors []*picture.Picture) []*picture.Picture {
	newPics := make([]*picture.Picture, numPics)
	i := 0
	for i < len(survivors) {
		a := survivors[i]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = a.Cross(b)
		i++
	}

	for i < len(newPics) {
		a := survivors[rand.Intn(len(survivors))]
		b := survivors[rand.Intn(len(survivors))]
		newPics[i] = a.Cross(b)
		i++
	}

	for _, pic := range newPics {
		r := rand.Intn(mutationRate)
		for i := 0; i < r; i++ {
			pic.Mutate()
		}
	}

	return newPics
}

func newImage(p *picture.Picture, width, height int32) *rl.Image {
	scale := 128.0
	offset := -1 * scale
	index := 0
	var imageData = make([]byte, width*height*4)

	for y := int32(0); y < height; y++ {
		yy := float64(y)/float64(height)*2 - 1
		for x := int32(0); x < width; x++ {
			xx := float64(x)/float64(width)*2 - 1
			r := p.R.Evaluate(xx, yy)
			g := p.G.Evaluate(xx, yy)
			b := p.B.Evaluate(xx, yy)

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
