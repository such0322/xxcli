package client

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"image/color"
	"time"
	"xiuxian/protocol"
)

//go:embed excel.ttf
var excelFont []byte

const (
	//ScreenWidth  = 640
	//ScreenHeight = 480

	ScreenWidth  = 1200
	ScreenHeight = 1200
)

type Game struct {
	World         IWorld
	Width, Height int
	FontFace      font.Face
}

func NewGame() *Game {
	g := &Game{
		Width:  ScreenWidth,
		Height: ScreenHeight,
	}
	fontData, _ := truetype.Parse(excelFont)
	g.FontFace = truetype.NewFace(fontData, &truetype.Options{Size: 10})
	g.World = NewNavWorld(g)
	g.World.Init()
	return g

	//g.World = NewWorld(g)
	//

	//
	//if err := g.tryConnect("127.0.0.1:3250"); err != nil {
	//	panic("servers not connect")
	//}
	//go ReadServerMessages(g)
	//
	//g.afterConnect()
	//g.World.Init()
	//Request(WorldGetState, &protocol.MsgWorldGetStateReq{})
	////printFPS()
	//return g
}

func printFPS() {
	go func() {
		for {
			fmt.Println("FPS: ", ebiten.ActualFPS())
			fmt.Println("Ticks: ", ebiten.ActualTPS())
			time.Sleep(time.Second)
		}
	}()
}

func (g *Game) Update() error {
	var quit error
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		quit = errors.New("quit")
	}
	g.World.Update()

	return quit
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen)
}

func (g *Game) DrawText(screen *ebiten.Image, x, y int, textLines ...string) {
	rectHeight := 10
	for _, txt := range textLines {
		w := float64(font.MeasureString(g.FontFace, txt).Round())
		ebitenutil.DrawRect(screen, float64(x), float64(y-8), w, float64(rectHeight), color.RGBA{0, 0, 0, 192})
		text.Draw(screen, txt, g.FontFace, x+1, y+1, color.RGBA{0, 0, 150, 255})
		text.Draw(screen, txt, g.FontFace, x, y, color.RGBA{100, 150, 255, 255})
		y += rectHeight

	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.Width, g.Height
}

func (g *Game) tryConnect(addr string) error {
	if err := agent.Cli.ConnectTo(addr); err != nil {
		return err
	}
	return nil
}

func (g *Game) afterConnect() {
	time.Sleep(1 * time.Second)
	agent.request(GateLogin, &protocol.MsgGateLoginReq{
		Account: "moz23",
	})
	time.Sleep(time.Millisecond * 100)
	Request(PublicSelectRole, &protocol.MsgPublicSelectRoleReq{})
	time.Sleep(time.Millisecond * 100)
	fmt.Println("Request PublicSelectRole done")
	Request(GameStartGame, &protocol.MsgGameStartGameReq{})
	time.Sleep(time.Millisecond * 400)
	fmt.Println("Request GameStartGame done")
}
