package main

import rl "github.com/gen2brain/raylib-go/raylib"

const fontSize = 64

type Button struct {
	Index          int32
	Rectangle      rl.Rectangle
	Texture        rl.Texture2D
	Text           string
	Selected       bool
	IsLeftClicked  func()
	IsRightClicked func(*Button)
}

var font rl.Font

func initFont() {
	font = rl.LoadFontEx("MesloLGLDZNerdFont-Bold.ttf", fontSize, []rune("Evolve!"), 0)
}

func NewButton(index int32, rectangle rl.Rectangle, texture rl.Texture2D, isRightClicked func(*Button)) *Button {
	if font.BaseSize == 0 {
		initFont()
	}
	return &Button{
		Index:          index,
		Rectangle:      rectangle,
		Texture:        texture,
		Text:           "",
		IsRightClicked: isRightClicked,
		IsLeftClicked:  nil,
	}
}

func NewTextButton(rectangle rl.Rectangle, text string, isLeftClicked func()) *Button {
	if font.BaseSize == 0 {
		initFont()
	}
	return &Button{
		Index:          -1,
		Rectangle:      rectangle,
		Text:           text,
		IsLeftClicked:  isLeftClicked,
		IsRightClicked: nil,
	}
}

func (b *Button) Update() {
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), b.Rectangle) {
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			if b.Text == "" {
				b.Selected = !b.Selected
			} else {
				// Button was clicked
				b.IsLeftClicked()
			}
		} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
			b.IsRightClicked(b)
		}
	}
}

func (b *Button) Draw() {
	if b.Text == "" {
		rl.DrawTexture(b.Texture, int32(b.Rectangle.X), int32(b.Rectangle.Y), rl.White)
	} else {
		rl.DrawRectangleRec(b.Rectangle, rl.White)
		tw := rl.MeasureTextEx(font, b.Text, fontSize, 0)
		x := b.Rectangle.X + b.Rectangle.Width/2 - tw.X/2
		y := b.Rectangle.Y + b.Rectangle.Height/2 - tw.Y/2
		r := rl.Vector2{
			X: x,
			Y: y,
		}
		//rl.DrawText(b.Text, x, y, fontSize, rl.Black)
		rl.DrawTextEx(font, b.Text, r, fontSize, 0, rl.Black)
	}
	if b.Selected {
		rl.DrawRectangleLinesEx(b.Rectangle, 2, rl.White)
	}
}
