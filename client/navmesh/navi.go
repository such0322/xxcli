package navmesh

import (
	"errors"
	"fmt"
	"github.com/kvartborg/vector"
)

var (
	ERROR_TRIANGLELIST_ILLEGAL = errors.New("triangle list illegal")
)

type TriangleList struct {
	Vertices  []Point3
	Triangles [][3]int32 // triangles
}

type BorderList struct {
	Indices []int32 // 2pt as border
}

type Path struct {
	Line []Point3
}

type NavMesh struct{}

func (nm *NavMesh) Route(list TriangleList, start, end Point3) (*Path, error) {
	r := Path{}
	// 计算临边
	fmt.Println("Vertices = ", list.Vertices)
	border := nm.createBorder(list.Triangles)
	fmt.Println("border = ", border)
	// 目标点
	vertices := append(list.Vertices, end)
	border = append(border, int32(len(vertices))-1, int32(len(vertices))-1)

	// 第一个可视区域
	lineStart := start
	lastVisLeft, lastVisRight, lastPLeft, lastPRight := nm.updateVis(start, vertices, border, 0, 1)
	var res vector.Vector
	for k := 2; k <= len(border)-2; k += 2 {
		curVisLeft, curVisRight, pLeft, pRight := nm.updateVis(lineStart, vertices, border, k, k+1)
		//V3Cross(&res, lastVisLeft, curVisRight)
		res, _ = lastVisLeft.Cross(curVisRight)
		if res.Z() > 0 { // 左拐点
			lineStart = vertices[border[lastPLeft]]
			r.Line = append(r.Line, lineStart)
			// 找到一条不共点的边作为可视区域
			i := 2 * (lastPLeft/2 + 1)
			for ; i <= len(border)-2; i += 2 {
				if border[lastPLeft] != border[i] && border[lastPLeft] != border[i+1] {
					lastVisLeft, lastVisRight, lastPLeft, lastPRight = nm.updateVis(lineStart, vertices, border, i, i+1)
					break
				}
			}

			k = i
			continue
		}

		//V3Cross(&res, lastVisRight, curVisLeft)
		res, _ = lastVisRight.Cross(curVisLeft)
		if res.Z() < 0 { // 右拐点
			lineStart = vertices[border[lastPRight]]
			r.Line = append(r.Line, lineStart)
			// 找到一条不共点的边
			i := 2 * (lastPRight/2 + 1)
			for ; i <= len(border)-2; i += 2 {
				if border[lastPRight] != border[i] && border[lastPRight] != border[i+1] {
					lastVisLeft, lastVisRight, lastPLeft, lastPRight = nm.updateVis(lineStart, vertices, border, i, i+1)
					break
				}
			}

			k = i
			continue
		}

		//V3Cross(&res, lastVisLeft, curVisLeft)
		res, _ = lastVisLeft.Cross(curVisLeft)
		if res.Z() < 0 {
			lastVisLeft = curVisLeft
			lastPLeft = pLeft
		}

		//V3Cross(&res, lastVisRight, curVisRight)
		res, _ = lastVisRight.Cross(curVisRight)
		if res.Z() > 0 {
			lastVisRight = curVisRight
			lastPRight = pRight
		}
	}
	fmt.Println("routePath = ", r)
	return &r, nil
}

func (nm *NavMesh) createBorder(list [][3]int32) []int32 {
	var border []int32
	for k := 0; k < len(list)-1; k++ {
		for _, i := range list[k] {
			for _, j := range list[k+1] {
				if i == j {
					border = append(border, i)
				}
			}
		}
	}
	return border
}

func (nm *NavMesh) updateVis(v0 Point3, vertices []Point3, indices []int32, i1, i2 int) (l, r vector.Vector, left, right int) {
	leftVec := P3Sub(vertices[indices[i1]], v0)
	rightVec := P3Sub(vertices[indices[i2]], v0)
	res, _ := leftVec.Cross(rightVec)
	if res.Z() > 0 {
		return rightVec, leftVec, i2, i1
	} else {
		return leftVec, rightVec, i1, i2
	}
}

type Point3 struct {
	X, Y, Z, _ float64
}

func P3Sub(pnt0, pnt1 Point3) vector.Vector {
	return vector.Vector{
		pnt0.X - pnt1.X,
		pnt0.Y - pnt1.Y,
		pnt0.Z - pnt1.Z,
	}
}
