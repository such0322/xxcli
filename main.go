package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	pCLi "github.com/topfreegames/pitaya/v2/client"
	"xiuxianclient/client"
)

var pClient pCLi.PitayaClient

func main() {
	client.InitAgent()
	game := client.NewGame()
	ebiten.SetWindowSize(client.ScreenWidth, client.ScreenHeight)
	ebiten.SetWindowTitle("moz xiuxian client")
	ebiten.RunGame(game)
	client.Disconnect()

	//time.Sleep(2 * time.Second)
}
