package main

import rl "github.com/gen2brain/raylib-go/raylib"

const fontSize = 24

type Button struct {
	Rectangle rl.Rectangle
	Texture   rl.Texture2D
	Text      string
	Selected  bool
	IsClicked func()
}

func NewButton(rectangle rl.Rectangle, texture rl.Texture2D) *Button {
	return &Button{
		Rectangle: rectangle,
		Texture:   texture,
		Text:      "",
	}
}

func NewTextButton(rectangle rl.Rectangle, text string, isClicked func()) *Button {
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
		tw := rl.MeasureText(b.Text, fontSize)
		x := int32(b.Rectangle.X + b.Rectangle.Width/2 - float32(tw)/2)
		y := int32(b.Rectangle.Y + b.Rectangle.Height/2 - fontSize/2)
		rl.DrawText(b.Text, x, y, fontSize, rl.Black)
	}
	if b.Selected {
		rl.DrawRectangleLinesEx(b.Rectangle, 2, rl.White)
	}
}
