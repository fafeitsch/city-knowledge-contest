package contest

import (
	"fmt"
	"math"
	"math/rand"
)

type polygon []Coordinate

type triangle [3]Coordinate

type triangulation []triangle

func (p polygon) computeTriangulation() triangulation {
	if len(p) < 3 {
		panic(fmt.Sprintf("the passed polygon has only %d vertices, but must have at least 3", len(p)))
	}
	current := polygon(append([]Coordinate{}, p...))
	result := make([]triangle, 0, len(p)-2)
	for len(current) > 3 {
		index := 0
		for index < len(current) && !current.isEar(index) {
			index = index + 1
		}
		if index == len(current) {
			for i, j := 0, len(current)-1; i < j; i, j = i+1, j-1 {
				current[i], current[j] = current[j], current[i]
			}
			continue
		}
		previousIndex := index - 1
		if previousIndex < 0 {
			previousIndex = len(current) - 1
		}
		result = append(result, triangle{current[index], current[previousIndex], current[(index+1)%len(current)]})
		current = append(current[:index], current[index+1:]...)
	}
	return append(result, triangle{current[0], current[1], current[2]})
}

func (t triangulation) randomPoint(random *rand.Rand) Coordinate {
	areas := make([]float64, len(t))
	sum := 0.0
	for index, tria := range t {
		sum = sum + tria.area()
		areas[index] = sum
	}
	die := random.Float64() * sum
	index := 0
	for areas[index] < die {
		index = index + 1
	}
	return t[index].randomPoint(random)
}

func (p polygon) isEar(index int) bool {
	nextIndex := (index + 1) % len(p)
	nextVertex := p[nextIndex]
	previousIndex := index - 1
	if previousIndex < 0 {
		previousIndex = len(p) - 1
	}
	previousVertex := p[previousIndex]
	vertex := p[index]
	index = (nextIndex + 1) % len(p)
	possibleEar := triangle([3]Coordinate{previousVertex, vertex, nextVertex})
	if !isLeftOfLine(possibleEar[2], possibleEar[0], possibleEar[1]) {
		return false
	}
	isEar := true
	for index != previousIndex && isEar {
		isEar = isEar && !possibleEar.containsPoint(p[index])
		index = (index + 1) % len(p)
	}
	return isEar
}

func isLeftOfLine(toTest Coordinate, start Coordinate, end Coordinate) bool {
	d := (toTest.Lng-start.Lng)*(end.Lat-start.Lat) - (toTest.Lat-start.Lat)*(end.Lng-start.Lng)
	return d < 0
}

func (t triangle) containsPoint(coordinate Coordinate) bool {
	d1 := isLeftOfLine(t[0], t[1], t[2]) == isLeftOfLine(coordinate, t[1], t[2])
	d2 := isLeftOfLine(t[1], t[2], t[0]) == isLeftOfLine(coordinate, t[2], t[0])
	d3 := isLeftOfLine(t[2], t[0], t[1]) == isLeftOfLine(coordinate, t[0], t[1])
	return d1 && d2 && d3
}

func (t triangle) area() float64 {
	m1 := (t[0].Lat - t[1].Lat) / (t[0].Lng - t[1].Lng)
	m2 := (t[0].Lat - t[2].Lat) / (t[0].Lng - t[2].Lng)
	tan := (m1 - m2) / (1 + m1*m2)
	angle := math.Abs(math.Atan(tan))
	a := math.Sqrt(math.Pow(t[0].Lat-t[1].Lat, 2) + math.Pow(t[0].Lng-t[1].Lng, 2))
	b := math.Sqrt(math.Pow(t[0].Lat-t[2].Lat, 2) + math.Pow(t[0].Lng-t[2].Lng, 2))
	return 0.5 * a * b * math.Sin(angle)
}

func (t triangle) randomPoint(random *rand.Rand) Coordinate {
	a := random.Float64()
	b := random.Float64()
	if a+b > 1 {
		a = 1 - a
		b = 1 - b
	}
	x1, y1 := t[1].Lng-t[0].Lng, t[1].Lat-t[0].Lat
	x2, y2 := t[2].Lng-t[0].Lng, t[2].Lat-t[0].Lat
	px, py := a*x1+b*x2, a*y1+b*y2
	point := Coordinate{Lng: t[0].Lng + px, Lat: t[0].Lat + py}
	if !t.containsPoint(point) {
		// panic("p is not contained in triangle")
	}
	return point
}
