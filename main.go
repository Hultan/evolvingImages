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

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Evolving Images")
	rl.SetTraceLogLevel(rl.LogNone)

	// Generate image
	aptR := GetRandomNode()
	aptG := GetRandomNode()
	aptB := GetRandomNode()

	const nodes = 20

	num := rand.Intn(nodes)
	for i := 0; i < num; i++ {
		aptR.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes)
	for i := 0; i < num; i++ {
		aptG.AddRandom(GetRandomNode())
	}

	num = rand.Intn(nodes)
	for i := 0; i < num; i++ {
		aptB.AddRandom(GetRandomNode())
	}

	for {
		_, nilCount := aptR.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptR.AddRandom(GetRandomLeafNode())
	}

	for {
		_, nilCount := aptG.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptG.AddRandom(GetRandomLeafNode())
	}
	for {
		_, nilCount := aptB.NodeCounts()
		if nilCount == 0 {
			break
		}
		aptB.AddRandom(GetRandomLeafNode())
	}

	fmt.Println(aptR.String())
	fmt.Println()
	fmt.Println(aptG.String())
	fmt.Println()
	fmt.Println(aptB.String())
	fmt.Println()

	generateImage(aptR, aptG, aptB, screenWidth, screenHeight)

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

func generateImage(red, green, blue Node, width, height int) {
	scale := 128.0
	offset := -1 * scale
	index = 0

	for y := 0; y < height; y++ {
		yy := float64(y)/float64(height)*2 - 1
		for x := 0; x < width; x++ {
			xx := float64(x)/float64(width)*2 - 1
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

	image.Data = unsafe.Pointer(unsafe.SliceData(imageData))
	texture = rl.LoadTextureFromImage(image)
}
