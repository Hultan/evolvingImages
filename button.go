package main

import rl "github.com/gen2brain/raylib-go/raylib"

const fontSize = 64

type Button struct {
	Rectangle rl.Rectangle
	Texture   rl.Texture2D
	Text      string
	Selected  bool
	IsClicked func()
}

var font rl.Font

func initFont() {
	font = rl.LoadFontEx("MesloLGLDZNerdFont-Bold.ttf", fontSize, []rune("Evolve!"), 0)
}

func NewButton(rectangle rl.Rectangle, texture rl.Texture2D) *Button {
	if font.BaseSize == 0 {
		initFont()
	}
	return &Button{
		Rectangle: rectangle,
		Texture:   texture,
		Text:      "",
	}
}

func NewTextButton(rectangle rl.Rectangle, text string, isClicked func()) *Button {
	if font.BaseSize == 0 {
		initFont()
	}
	return &Button{
		Rectangle: rectangle,
		Text:      text,
		IsClicked: isClicked,
	}
}

func (b *Button) Update() {
	if rl.CheckCollisionPointRec(rl.GetMousePosition(), b.Rectangle) {
		if rl.IsMouseButtonReleased(rl.MouseButtonLeft) {
			if b.Text == "" {
				b.Selected = !b.Selected
			} else {
				// Button was clicked
				b.IsClicked()
			}
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
