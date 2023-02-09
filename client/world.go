package client

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kvartborg/vector"
	"github.com/solarlune/resolv"
	"image/color"
	"math"
	"xiuxian/common/consts"
	"xiuxian/protocol"
)

type World struct {
	Game        *Game
	Run         bool
	Contact     *resolv.ContactSet
	Solid       *resolv.ConvexPolygon
	SolidOther  *resolv.ConvexPolygon
	CircleOne   *resolv.Circle
	CircleTwo   *resolv.Circle
	CircleOther []*resolv.Circle
}

func NewWorld(game *Game) *World {
	return &World{
		Game: game,
	}
}

func (w *World) Init() {
	//多边形
	//solid := resolv.NewConvexPolygon(0, 0,
	//	100, 100,
	//	250, 80,
	//	300, 150,
	//	250, 250,
	//	150, 300,
	//	80, 150,
	//)
	size := 1.0
	v0 := vector.Vector{600, 600}
	p1 := vector.Vector{-140.61533684577087, -297.97480971567506}
	p1 = p1.Add(v0).Scale(size)
	fmt.Println("p1 = ", p1)
	p2 := vector.Vector{-17.041775554737015, -382.1022845505485}
	p2 = p2.Add(v0).Scale(size)
	fmt.Println("p2 = ", p2)

	//var pp = []vector.Vector{
	//	{158.43974761489366, -364.7796441602451},
	//	{-4.180154619212168, -300.61832800238},
	//	{-17.041775554737015, -382.1022845505485},
	//	//{47.593449964682584, 572.3500316452687},
	//	//{47.38993480677577, 572.7553909639622},
	//	//{28.512534420682645, 89.26813107739937},
	//}
	var other []*resolv.Circle
	//
	//for _, p := range pp {
	//	p = p.Add(v0).Scale(size)
	//	c := resolv.NewCircle(p.X(), p.Y(), 12)
	//	fmt.Println(c.X, c.Y)
	//	other = append(other, c)
	//
	//}

	//ppp := []float64{
	//	-40.61533684577087, -297.97480971567506,
	//	-250.31399137940355, -512.5126252326233,
	//	-209.87409692155882, -545.667102970041,
	//	-164.29136068787642, -571.2955757152729,
	//	-114.9507910528395, -588.619335292588,
	//	-63.35157625727557, -597.1120080334749,
	//	-11.061532332285296, -596.5155483742017,
	//	40.33053419750935, -586.848079439692,
	//	89.26310139624105, -568.4033423815457,
	//	134.24937797065013, -541.7417712017874,
	//	173.9224786711775, -507.67346424930776,
	//	207.07695640859515, -467.233569791463,
	//	232.70542915382697, -421.65083355778063,
	//	250.02918873114209, -372.31026392274373,
	//}
	//w.SolidOther = resolv.NewConvexPolygon(600, 600,
	//	ppp...,
	//)
	//fmt.Println("++++++1", w.SolidOther.Transformed())
	//fmt.Println("------1", w.SolidOther.PointInside(vector.Vector{-17.041775554737015, -382.1022845505485}))
	//
	//fmt.Println("========================================================================")
	////w.SolidOther.SetScale(2, 2)
	//fmt.Println("++++++2", w.SolidOther.Transformed())
	////fmt.Println("------2", w.SolidOther.PointInside(vector.Vector{36.296569013813373, -16.831325233651895}))
	////w.SolidOther.MoveVec(p1)
	////fmt.Println("++++++2", w.SolidOther.Transformed())

	//长方形带旋转
	//solid := resolv.NewConvexPolygon(0, 0,
	//	0, 0,
	//	0, 1,
	//	1, 1,
	//	1, 0,
	//)
	//solid.SetScale(300, 200)
	//solid.MoveVec(p1)
	//vec := p2.Sub(p1)
	//fmt.Println("cur ag = ", solid.Rotation())
	//
	////solid.Move(0, -solid.ScaleH/2)
	////solid.Rotate(-vec.Angle())
	////vu := vec.Unit()
	////solid.Move(vu.Y()*solid.ScaleH/2, -vu.X()*solid.ScaleH*2)
	//
	//fmt.Println("vec =", vec)
	//solid.Rotate(-vec.Angle())
	//vu := vec.Unit()
	//fmt.Println("sw ,sh = ", solid.ScaleW, solid.ScaleH)
	//solid.Move(vu.Y()*solid.ScaleH/2, -vu.X()*solid.ScaleH/2)
	//fmt.Println("cur ag2 = ", solid.Rotation())

	//扇形/伪扇形
	ro := 12.0
	points := []float64{0, 0, 1, 0}
	ag10 := -math.Pi / 18
	for i := 0.0; i < ro; i++ {
		sin, cos := math.Sincos(ag10 * (i + 1))
		points = append(points, cos, -sin)
	}
	solid := resolv.NewConvexPolygon(0, 0,
		points...,
	)
	solid.Rotate(-ag10 * ro / 2)

	solid.SetScale(200, 200)
	fmt.Println("p1 = ", p1)
	solid.MoveVec(p1)
	fmt.Println("cur ag = ", solid.Rotation())
	vec := p2.Sub(p1)
	fmt.Println("p2 = ", p2)
	vu := vec.Unit()
	fmt.Println(solid.ScaleH)
	fmt.Println("vu = ", vu)
	solid.Move(vec.X()*0.95, vec.Y()*0.95)

	solid.Rotate(-vec.Angle())
	fmt.Println("points = ", solid.Transformed())

	w.Solid = solid
	w.CircleOne = resolv.NewCircle(p1.X(), p1.Y(), 16)
	w.CircleTwo = resolv.NewCircle(p2.X(), p2.Y(), 4)
	w.CircleOther = other
	w.Run = true

}

func (w *World) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		fmt.Println("pos =", float64(x)-600, float64(y)-600, x, y)
		Notify(consts.HandlerWorldMoveToNtf, &protocol.MsgWorldMoveToNtf{
			DstPos: &protocol.Pos{
				X: float64(x) / 100,
				Y: 0,
				Z: float64(y) / 100,
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

	for _, circle := range w.CircleOther {
		vec := vector.Vector{circle.X, circle.Y}
		if w.SolidOther.PointInside(vec) {
			//fmt.Println("vec = ", i, vec)
			DrawCircle(screen, circle, color.RGBA{244, 123, 0, 255})
		} else {
			DrawCircle(screen, circle, color.White)
		}

	}
	if w.SolidOther != nil {
		DrawPolygon(screen, w.SolidOther, color.RGBA{244, 123, 0, 255})
	}

	w.Game.DrawText(screen, 16, 16,
		fmt.Sprintf("%d FPS (frames per second)", int(ebiten.ActualFPS())),
		fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.ActualTPS())),
	)
}
