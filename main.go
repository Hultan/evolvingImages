package main

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	. "github.com/hultan/evolvingImage/apt"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var texture rl.Texture2D
var index int32
var imageData = make([]byte, screenWidth*screenHeight*4)
var image = rl.NewImage(imageData, screenWidth, screenHeight, 1, rl.UncompressedR8g8b8a8)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	// Generate image
	opX := &OperatorX{}
	opY := &OperatorY{}
	opSin := &OperatorSin{}
	opNoise := &OperatorNoise{}
	opAtan2 := &OperatorMult{}
	opPlus := &OperatorPlus{}

	//opT := &OperatorT{}
	opAtan2.LeftChild = opX
	opAtan2.RightChild = opNoise
	opNoise.LeftChild = opX
	opNoise.RightChild = opY
	opSin.Child = opAtan2
	opPlus.LeftChild = opY
	opPlus.RightChild = opSin

	generateImage(opPlus, opPlus, opPlus)

	rl.SetTargetFPS(60)
	for !rl.WindowShouldClose() {
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

func generateImage(red, green, blue Node) {
	generateImageData(red, green, blue)

	image.Data = unsafe.Pointer(unsafe.SliceData(imageData))
	texture = rl.LoadTextureFromImage(image)
}

func generateImageData(red, green, blue Node) {

	scale := 128.0
	offset := -1 * scale
	index = 0

	for y := 0; y < screenHeight; y++ {
		yy := float64(y)/screenHeight*2 - 1
		for x := 0; x < screenWidth; x++ {
			xx := float64(x)/screenWidth*2 - 1
			r := red.Evaluate(xx, yy)
			g := green.Evaluate(xx, yy)
			b := blue.Evaluate(xx, yy)

			imageData[index+0] = byte(r*scale - offset)
			imageData[index+1] = byte(g*scale - offset)
			imageData[index+2] = byte(b*scale - offset)
			imageData[index+3] = 255
			index += 4
		}
	}
}
