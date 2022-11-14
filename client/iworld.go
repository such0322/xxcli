package client

import "github.com/hajimehoshi/ebiten/v2"

type IWorld interface {
	Init()
	Update()
	Draw(image *ebiten.Image)
}
