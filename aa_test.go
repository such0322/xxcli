package main

import (
	"fmt"
	"github.com/kvartborg/vector"
	"github.com/solarlune/resolv"
	"testing"
)

func TestLine(t *testing.T) {
	s := 10.0
	l1 := resolv.NewLine(1.5238, -1.9883, 0.1498, 0.6785)
	l1.SetScale(s, s)
	vec := vector.Vector{0.3629, -0.1683}
	vec.Scale(s)
	b := l1.PointInside(vec)
	fmt.Println("bool = ", b)
}
