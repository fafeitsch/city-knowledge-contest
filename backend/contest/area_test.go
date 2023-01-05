package contest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
)

func Test_computeTriangulation(t *testing.T) {
	t.Run("square", func(t *testing.T) {
		p := polygon{
			Coordinate{Lng: 1, Lat: 2},
			Coordinate{Lng: 4, Lat: 2},
			Coordinate{Lng: 4, Lat: 5},
			Coordinate{Lng: 1, Lat: 5},
		}
		result := p.computeTriangulation()
		expected := triangulation([]triangle{{p[0], p[3], p[1]}, {p[1], p[2], p[3]}})
		assert.Equal(t, expected, result)
	})
	t.Run("simple polygon", func(t *testing.T) {
		p := polygon{
			Coordinate{Lng: 1, Lat: 4},
			Coordinate{Lng: -1, Lat: 1},
			Coordinate{Lng: -2, Lat: 0},
			Coordinate{Lng: 1, Lat: 0},
			Coordinate{Lng: 0, Lat: 2},
			Coordinate{Lng: 3, Lat: 4},
			Coordinate{Lng: 0, Lat: 5},
			Coordinate{Lng: 1, Lat: 7},
			Coordinate{Lng: -5, Lat: 5},
		}
		result := p.computeTriangulation()
		assert.Equal(t, triangle{p[2], p[1], p[3]}, result[0])
		assert.Equal(t, triangle{p[3], p[1], p[4]}, result[1])
		assert.Equal(t, triangle{p[1], p[0], p[4]}, result[2])
		assert.Equal(t, triangle{p[4], p[0], p[5]}, result[3])
		assert.Equal(t, triangle{p[0], p[8], p[5]}, result[4])
		assert.Equal(t, triangle{p[5], p[8], p[6]}, result[5])
		assert.Equal(t, triangle{p[6], p[7], p[8]}, result[6])
		assert.Equal(t, 7, len(result))
	})
}

func Test_triangulation_randomPoint(t *testing.T) {
	p := polygon{
		Coordinate{Lng: 1, Lat: 4},
		Coordinate{Lng: -1, Lat: 1},
		Coordinate{Lng: -2, Lat: 0},
		Coordinate{Lng: 1, Lat: 0},
		Coordinate{Lng: 0, Lat: 2},
		Coordinate{Lng: 3, Lat: 4},
		Coordinate{Lng: 0, Lat: 5},
		Coordinate{Lng: 1, Lat: 7},
		Coordinate{Lng: -5, Lat: 5},
	}
	source := rand.NewSource(567)
	random := rand.New(source)
	result := make([]string, 5)
	for index, _ := range result {
		point := p.computeTriangulation().randomPoint(random)
		result[index] = fmt.Sprintf("%.2f,%.2f", point.Lng, point.Lat)
	}
	assert.Equal(t, "1.17,3.60;1.41,3.13;-1.51,4.71;-0.12,5.48;0.47,6.75", strings.Join(result, ";"))
}

func Test_triangle_area(t *testing.T) {
	tr := triangle([3]Coordinate{
		{
			Lat: 3,
			Lng: 2,
		},
		{Lng: 9, Lat: -2},
		{Lng: 7, Lat: 2},
	})
	assert.InDelta(t, 9, tr.area(), 0.1)
}
