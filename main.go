package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	. "github.com/hultan/evolvingImage/apt"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var texture rl.Texture2D
var t, index int32
var imageData = make([]byte, screenWidth*screenHeight*4)
var image = rl.NewImage(imageData, screenWidth, screenHeight, 1, rl.UncompressedR8g8b8a8)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	for !rl.WindowShouldClose() {
		// Generate image
		opX := &OperatorX{}
		opY := &OperatorY{}
		opPlus := &OperatorPlus{}
		opSin := &OperatorSin{}
		//opT := &OperatorT{}

		opSin.Child = opX
		opPlus.LeftChild = opSin
		opPlus.RightChild = opY

		generateImage(opPlus, opSin)

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		rl.DrawTexture(texture, 0, 0, rl.White)
		rl.DrawFPS(10, 10)
		rl.EndDrawing()

		t++
	}

	// Clean up
	imageData = nil
	//rl.UnloadImage(image)
	rl.UnloadTexture(texture)

	rl.CloseWindow()
}

func generateImage(node, node2 Node) {
	generateImageData(node, node2)

	image.Data = unsafe.Pointer(unsafe.SliceData(imageData))
	texture = rl.LoadTextureFromImage(image)
}

func generateImageData(node, node2 Node) {

	scale := 128.0
	offset := -1 * scale
	index = 0

	tt := float64(t%10000)/10000*2 - 1
	for y := 0; y < screenHeight; y++ {
		yy := float64(y)/screenHeight*2 - 1
		for x := 0; x < screenWidth; x++ {
			xx := float64(x)/screenWidth*2 - 1
			c := node.Evaluate(xx, yy, tt)
			c2 := node2.Evaluate(xx, yy, tt)

			imageData[index+0] = byte(c*scale - offset)
			imageData[index+1] = byte(c2*scale - offset)
			imageData[index+2] = 0 //byte(c*scale - offset)
			imageData[index+3] = 255
			index += 4
		}
	}
}
