package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"image/color"
	"xiuxian/protocol"
)

type World struct {
	Game      *Game
	Run       bool
	Contact   *resolv.ContactSet
	Solid     *resolv.ConvexPolygon
	CircleOne *resolv.Circle
	CircleTwo *resolv.Circle
}

func NewWorld(game *Game) *World {
	return &World{
		Game: game,
	}
}

func (w *World) Init() {
	w.Solid = resolv.NewConvexPolygon(
		100, 100,
		250, 80,
		300, 150,
		250, 250,
		150, 300,
		80, 150,
	)

	w.CircleOne = resolv.NewCircle(500, 200, 16)
	w.CircleTwo = resolv.NewCircle(400, 250, 32)
	w.Run = true

}

func (w *World) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Println("pos =", x, y)
		Notify(WorldMoveToNtf, &protocol.MsgWorldMoveToNtf{
			DstPos: &protocol.Pos{
				X: float64(x),
				Y: 0,
				Z: float64(y),
			},
		})
	}

}

func (w *World) Draw(screen *ebiten.Image) {
	controllingColor := color.RGBA{0, 255, 80, 255}
	if w.Contact != nil {
		controllingColor = color.RGBA{160, 0, 0, 255}
	}

	DrawPolygon(screen, w.Solid, color.White)
	DrawCircle(screen, w.CircleOne, controllingColor)
	DrawCircle(screen, w.CircleTwo, color.White)

	w.Game.DrawText(screen, 16, 16,
		fmt.Sprintf("%d FPS (frames per second)", int(ebiten.ActualFPS())),
		fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.ActualTPS())),
	)
}
