package navmesh

import (
	"container/heap"
	"math"
)

const LARGE_NUMBER = math.MaxInt32

// Triangle Heap
type WeightedTriangle struct {
	id     int32 // triangle id
	weight uint32
}

type TriangleHeap struct {
	triangles []WeightedTriangle
	indices   map[int32]int
}

func NewTriangleHeap() *TriangleHeap {
	h := new(TriangleHeap)
	h.indices = make(map[int32]int)
	return h
}

func (th *TriangleHeap) Len() int {
	return len(th.triangles)
}

func (th *TriangleHeap) Less(i, j int) bool {
	return th.triangles[i].weight < th.triangles[j].weight
}

func (th *TriangleHeap) Swap(i, j int) {
	th.triangles[i], th.triangles[j] = th.triangles[j], th.triangles[i]
	th.indices[th.triangles[i].id] = i
	th.indices[th.triangles[j].id] = j
}

func (th *TriangleHeap) Push(x interface{}) {
	th.triangles = append(th.triangles, x.(WeightedTriangle))
	n := len(th.triangles)
	th.indices[th.triangles[n-1].id] = n - 1
}

func (th *TriangleHeap) Pop() interface{} {
	n := len(th.triangles)
	x := th.triangles[n-1]
	th.triangles = th.triangles[:n-1]
	return x
}

func (th *TriangleHeap) DecreaseKey(id int32, weight uint32) {
	if index, ok := th.indices[id]; ok {
		th.triangles[index].weight = weight
		heap.Fix(th, index)
		return
	} else {
		heap.Push(th, WeightedTriangle{id, weight})
	}
}

type Mesh struct {
	Vertices  []Point3   // vertices
	Triangles [][3]int32 // triangles
}

// Dijkstra
type Dijkstra struct {
	Matrix map[int32][]WeightedTriangle // all edge for nodes
}

// create neighbour matrix
func (d *Dijkstra) CreateMatrixFromMesh(mesh Mesh) {
	//fmt.Println("+++++++++++++++++++++++", mesh.Vertices)
	//fmt.Println("-----------------------", mesh.Triangles)
	//fmt.Println("--------------------------------------------------------------")
	d.Matrix = make(map[int32][]WeightedTriangle)
	for i := 0; i < len(mesh.Triangles); i++ {
		for j := 0; j < len(mesh.Triangles); j++ {
			if i == j {
				continue
			}
			//if i == 237 && (j == 236) {
			//	ccc := intersect2(mesh, i, j)
			//	fmt.Println("i2 =", i, mesh.Triangles[i], "j2 = ", j, mesh.Triangles[j], " ilen =", ccc)
			//	for k, v := range mesh.Triangles[i] {
			//		fmt.Println(k, "--", v, "vert = ", mesh.Vertices[v])
			//	}
			//	for k, v := range mesh.Triangles[j] {
			//		fmt.Println(k, "--", v, "vert = ", mesh.Vertices[v])
			//	}
			//}

			//len(intersect(mesh.Triangles[i], mesh.Triangles[j]))
			if intersect2(mesh, i, j) == 2 {
				x1 := (mesh.Vertices[mesh.Triangles[i][0]].X + mesh.Vertices[mesh.Triangles[i][1]].X + mesh.Vertices[mesh.Triangles[i][2]].X) / 3.0
				y1 := (mesh.Vertices[mesh.Triangles[i][0]].Y + mesh.Vertices[mesh.Triangles[i][1]].Y + mesh.Vertices[mesh.Triangles[i][2]].Y) / 3.0
				x2 := (mesh.Vertices[mesh.Triangles[j][0]].X + mesh.Vertices[mesh.Triangles[j][1]].X + mesh.Vertices[mesh.Triangles[j][2]].X) / 3.0
				y2 := (mesh.Vertices[mesh.Triangles[j][0]].Y + mesh.Vertices[mesh.Triangles[j][1]].Y + mesh.Vertices[mesh.Triangles[j][2]].Y) / 3.0
				weight := math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
				d.Matrix[int32(i)] = append(d.Matrix[int32(i)], WeightedTriangle{int32(j), uint32(weight)})
			}
		}
	}

	//fmt.Println("dx = ", d.Matrix)
}

func intersect2(mesh Mesh, i, j int) int {
	var inter int
	t1 := mesh.Triangles[i]
	t2 := mesh.Triangles[j]
	var vertices = make(map[int32]Point3)
	for _, p := range t1 {
		vertices[p] = mesh.Vertices[p]
	}
	for _, p := range t2 {
		vertices[p] = mesh.Vertices[p]
	}

	for _, vi := range t1 {
		for _, vj := range t2 {
			if vertices[vi] == vertices[vj] {
				inter++
			}
		}
	}
	return inter
}

func intersect(a [3]int32, b [3]int32) []int32 {
	var inter []int32
	for i := range a {
		for j := range b {
			if a[i] == b[j] {
				inter = append(inter, a[i])
			}
		}
	}
	return inter
}

func (d *Dijkstra) Run(srcId int32) []int32 {
	// triangle heap
	h := NewTriangleHeap()
	// min distance records
	dist := make([]uint32, len(d.Matrix))
	for i := 0; i < len(dist); i++ {
		dist[i] = LARGE_NUMBER
	}
	// previous
	prev := make([]int32, len(d.Matrix))
	//fmt.Println("prev =", prev, " len =", len(prev))
	for i := 0; i < len(prev); i++ {
		prev[i] = -1
	}
	// visit map
	visited := make([]bool, len(d.Matrix))

	// source vertex, the first vertex in Heap
	dist[srcId] = 0
	heap.Push(h, WeightedTriangle{srcId, 0})

	for h.Len() > 0 { // for every un-visited vertex, try relaxing the path
		// pop the min element
		u := heap.Pop(h).(WeightedTriangle)
		//fmt.Println("curr u = ", u)
		if visited[u.id] {
			continue
		}
		// current known shortest distance to u
		distU := dist[u.id]
		//fmt.Println("dist_u = ", dist_u)
		// mark the vertex as visited.
		visited[u.id] = true
		//fmt.Println(" d.Matrix[u.id] = ", u.id, d.Matrix[u.id])
		// for each neighbor v of u:
		for _, v := range d.Matrix[u.id] {
			//fmt.Println("v of u =", v)
			alt := distU + v.weight // from src->u->v
			if alt < dist[v.id] {
				dist[v.id] = alt
				prev[v.id] = u.id
				if !visited[v.id] {
					h.DecreaseKey(v.id, alt)
				}
				//fmt.Println("u =", u, v.id)
			}
		}
	}
	return prev
}
