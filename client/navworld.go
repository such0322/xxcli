package client

import (
	"encoding/json"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
	"image/color"
	"log"
	"os"
	"xiuxianclient/client/navmesh"
)

type CurList struct {
	Vertices []navmesh.Point3
	Indices  []int32
}

var list navmesh.TriangleList
var cList CurList
var vertices []navmesh.Point3
var triangles [][3]int32
var dijkstra navmesh.Dijkstra

type NavWorld struct {
	Game        *Game
	sf          float32
	ScaleFactor float32
	xx          float32
	yy          float32
	Px          float32
	Py          float32
	pp          float32

	Solids []*resolv.ConvexPolygon
	roads  []road

	path *navmesh.Path
}

func NewNavWorld(game *Game) *NavWorld {
	return &NavWorld{
		Game: game,
		path: &navmesh.Path{},
	}
}

func (w *NavWorld) Init() {
	//w.LoadOriginMesh()
	w.LoadCurrentMesh("./meshdata/10001.json")
	//fmt.Println("vertices = ", vertices)
	//fmt.Println("triangles = ", triangles)

	dijkstra.CreateMatrixFromMesh(navmesh.Mesh{Vertices: vertices, Triangles: triangles})

	startId := getTriangleId(w.roads[0].start)
	endId := getTriangleId(w.roads[0].end)
	if startId == -1 || endId == -1 {
		return
	}
	fmt.Println("sid eid =", startId, endId, w.roads[0].start, w.roads[0].end)
	path := route(startId, endId, w.roads[0].start, w.roads[0].end)
	if path != nil {
		w.path = path
	}

	//w.Solid = resolv.NewConvexPolygon(
	//	100, 100,
	//	250, 80,
	//	80, 150,
	//)

	//fmt.Println(w.Solid.Points)
	if w.path != nil {
		w.path.Line = append([]navmesh.Point3{w.roads[0].start}, w.path.Line...)
		w.path.Line = append(w.path.Line, w.roads[0].end)
	} else {
		w.path.Line = append([]navmesh.Point3{w.roads[0].start}, w.roads[0].end)
	}

	//fmt.Println("path = ", w.path)
}

func (w *NavWorld) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyJ) {
		w.ScaleFactor += w.sf
	}
	if ebiten.IsKeyPressed(ebiten.KeyK) {
		w.ScaleFactor -= w.sf
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		w.yy++
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		w.yy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		w.xx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		w.xx--
	}
	w.Py = w.yy * w.ScaleFactor * w.pp
	w.Px = w.xx * w.ScaleFactor * w.pp

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		posX, PosY := (float32(x)-w.Px)/w.ScaleFactor, (float32(y)-w.Py)/w.ScaleFactor
		posId := getTriangleId(navmesh.Point3{
			X: posX,
			Y: PosY,
		})
		fmt.Println("pos =", posX, PosY, "pid =", posId)
	}

	w.Solids = []*resolv.ConvexPolygon{}
	for k, _ := range triangles {
		//fmt.Println(v)
		var points []float64
		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][0]].X+w.Px))
		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][0]].Y+w.Py))

		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][1]].X+w.Px))
		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][1]].Y+w.Py))

		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][2]].X+w.Px))
		points = append(points, float64(w.ScaleFactor*vertices[triangles[k][2]].Y+w.Py))
		//fmt.Println(points)
		w.Solids = append(w.Solids, resolv.NewConvexPolygon(points...))
	}

}

func (w *NavWorld) Draw(screen *ebiten.Image) {
	//screen.Clear()

	for _, solid := range w.Solids {
		DrawPolygon(screen, solid, color.White)

		x := (solid.Points[0].X() + solid.Points[1].X() + solid.Points[2].X()) / 3
		y := (solid.Points[0].Y() + solid.Points[1].Y() + solid.Points[2].Y()) / 3

		np := navmesh.Point3{X: (float32(x) - w.Px) / w.ScaleFactor, Y: (float32(y) - w.Py) / w.ScaleFactor}
		id := getTriangleId(np)
		//fmt.Println("sanjiao dian id = ", id, np)
		w.Game.DrawText(screen, int(x), int(y),
			fmt.Sprintf("%d", id),
		)
	}
	controllingColor := color.RGBA{0, 255, 80, 255}
	//fmt.Println("draw path = ", path)

	if w.path != nil && len(w.path.Line) > 1 {
		for i := 0; i < len(w.path.Line)-1; i++ {
			cline := w.path.Line[i]
			nline := w.path.Line[i+1]
			ebitenutil.DrawLine(screen,
				float64(w.ScaleFactor*cline.X+w.Px),
				float64(w.ScaleFactor*cline.Y+w.Py),
				float64(w.ScaleFactor*nline.X+w.Px),
				float64(w.ScaleFactor*nline.Y+w.Py),
				controllingColor)
		}
	}

}

func (w *NavWorld) LoadOriginMesh() {
	w.pp = 1
	w.sf = 0.05
	w.ScaleFactor = 0.15
	w.xx = 10
	w.yy = 10
	w.Px = w.ScaleFactor * w.xx * w.pp
	w.Py = w.ScaleFactor * w.yy * w.pp

	f, err := os.Open("./meshdata/mesh.json")
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		log.Fatal(err)
	}
	vertices = list.Vertices
	triangles = list.Triangles
	//dijkstra.CreateMatrixFromMesh(navmesh.Mesh{vertices, triangles})

	w.roads = append(w.roads,
		road{
			start: navmesh.Point3{
				X: 1876,
				Y: 296,
				Z: 0,
			},
			end: navmesh.Point3{
				X: 3923,
				Y: 1116,
				Z: 0,
			},
		},
	)

}

func (w *NavWorld) LoadCurrentMesh(src string) {
	//0.5 10 40 40
	w.pp = 1000
	w.sf = 0.5 / w.pp
	w.ScaleFactor = 10 / w.pp
	w.xx = 40
	w.yy = 40
	w.Px = w.ScaleFactor * w.xx * w.pp
	w.Py = w.ScaleFactor * w.yy * w.pp

	f, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	if err := json.NewDecoder(f).Decode(&cList); err != nil {
		log.Fatal(err)
	}

	var nList CurList
	for i, vertex := range cList.Vertices {
		fmt.Println(i, "-----", vertex)
		none := true
		for _, cv := range nList.Vertices {
			if vertex == cv {
				none = false
			}
		}
		if none {
			nList.Vertices = append(nList.Vertices, vertex)
			for _, v := range cList.Indices {
				if i {

				}
				fmt.Println("v = ", v)
			}
		}
	}
	//fmt.Println("-----------------", nList.Vertices)
	fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++")
	for _, vertex := range cList.Vertices {
		list.Vertices = append(list.Vertices, navmesh.Point3{
			X: vertex.X * w.pp,
			Y: vertex.Z * w.pp,
		})
	}

	for i := 0; i < len(cList.Indices)/3; i++ {
		var t [3]int32
		t[0] = cList.Indices[i*3]
		t[1] = cList.Indices[i*3+1]
		t[2] = cList.Indices[i*3+2]
		list.Triangles = append(list.Triangles, t)
	}

	vertices = list.Vertices
	triangles = list.Triangles

	w.roads = append(w.roads,
		road{
			start: navmesh.Point3{
				X: 4945.112,
				Y: -5610.972,
				Z: 0,
			},
			end: navmesh.Point3{
				X: 7800.0034,
				Y: -999.99695,
				Z: 0,
			},
		},
	)
}

func route(srcId, dstId int32, src, dest navmesh.Point3) *navmesh.Path {
	djPath := dijkstra.Run(srcId)
	//fmt.Println("djPath =", djPath)

	pathTriangle := [][3]int32{triangles[dstId]}
	routeId := []int32{dstId}
	prevId := dstId
	fmt.Println("prevId =", prevId)
	for {
		curId := djPath[prevId]
		if curId == -1 {
			return nil
		}
		pathTriangle = append([][3]int32{triangles[curId]}, pathTriangle...)
		routeId = append([]int32{curId}, routeId...)
		if curId == srcId {
			break
		}
		prevId = curId
		fmt.Println("prevId =", prevId)
	}
	fmt.Println("routeId = ", routeId)
	fmt.Println("pathTriangle =", pathTriangle)
	nm := navmesh.NavMesh{}
	triList := navmesh.TriangleList{Vertices: vertices, Triangles: pathTriangle}
	path, err := nm.Route(triList, &src, &dest)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return path
}

func sign(p1, p2, p3 navmesh.Point3) float32 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Y-p3.Y)
}

func inside(pt, v1, v2, v3 navmesh.Point3) bool {
	b1 := sign(pt, v1, v2) <= 0
	b2 := sign(pt, v2, v3) <= 0
	b3 := sign(pt, v3, v1) <= 0
	return ((b1 == b2) && (b2 == b3))
}

func getTriangleId(pt navmesh.Point3) (id int32) {
	for k := 0; k < len(triangles); k++ {
		if inside(pt,
			navmesh.Point3{X: vertices[triangles[k][0]].X, Y: vertices[triangles[k][0]].Y},
			navmesh.Point3{X: vertices[triangles[k][1]].X, Y: vertices[triangles[k][1]].Y},
			navmesh.Point3{X: vertices[triangles[k][2]].X, Y: vertices[triangles[k][2]].Y}) {
			return int32(k)
		}
	}
	return -1
}

type road struct {
	start navmesh.Point3
	end   navmesh.Point3
}
