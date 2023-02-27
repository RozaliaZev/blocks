package main

import (
	"bufio"
	"bytes"
	"fmt"
	svg "github.com/ajstarks/svgo"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Color string

const (
	ColorUnspecified = ""
	ColorGreen       = "green"
	ColorBlue        = "blue"
	ColorRed         = "red"
	ColorPink        = "pink"
	ColorYellow      = "yellow"
	ColorBlueviolet  = "blueviolet"
	ColorDimgray     = "dimgray"
	ColorOrangered   = "orangered"
	ColorSpringgreen = "springgreen"
	ColorYellowgreen = "yellowgreen"
)

type Coord struct {
	i, j int
}

type Cell struct {
	color Color
	mark  bool
}

type Field [][]Cell

func main() {
	f, err := os.Open("temp/file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	wr := bytes.Buffer{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		wr.WriteString(sc.Text())
	}
	fmt.Println(wr.String())

	arr := strings.Split(wr.String(), " ")
	jSize, _ := strconv.Atoi(arr[1])
	iSize, _ := strconv.Atoi(arr[3])
	numColor, _ := strconv.Atoi(arr[5])

	field := generateField(iSize, jSize, numColor)

	maxGroup := field.findMaxSizeColorGroup()
	for _, coord := range maxGroup {
		field.setMark(coord, true)
	}

	http.Handle("/rect", http.HandlerFunc(field.rect))
	err = http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func generateField(iSize, jSize, numColor int) Field {
	setOfColors := []Color{ColorGreen, ColorBlue, ColorRed, ColorPink, ColorBlueviolet, ColorDimgray, ColorOrangered, ColorSpringgreen, ColorYellow, ColorYellowgreen}
	rand.Seed(time.Now().UnixNano())

	newField := make(Field, iSize)
	for i := range newField {
		newField[i] = make([]Cell, jSize)
		for j := range newField[i] {
			color := setOfColors[rand.Intn(numColor)]
			newField[i][j] = Cell{color: color}
		}
	}

	return newField
}

func (f Field) findMaxSizeColorGroup() []Coord {
	findField := f.Clone()
	var maxGroup []Coord

	for i, row := range findField {
		for j, cell := range row {
			color := cell.color
			if color == ColorUnspecified {
				continue
			}

			coord := Coord{i: i, j: j}
			group := findField.extractGroupByColor(coord, color)
			if len(group) > len(maxGroup) {
				maxGroup = group
			}
		}
	}

	return maxGroup
}

func (f Field) Clone() Field {
	if f == nil {
		return nil
	}

	newField := make(Field, len(f))
	for i := range newField {
		newField[i] = make([]Cell, len(f[i]))
		for j := range newField[i] {
			newField[i][j] = f[i][j]
		}
	}

	return newField
}

func (f Field) setMark(c Coord, value bool) {
	f[c.i][c.j].mark = value
}

func (f *Field) extractGroupByColor(coord Coord, targetColor Color) []Coord {
	if !f.inRange(coord) {
		return nil
	}

	color := f.getCell(coord).color
	if color == ColorUnspecified || color != targetColor {
		return nil
	}

	f.setColor(coord, ColorUnspecified)

	result := []Coord{coord}
	result = append(result, f.extractGroupByColor(coord.Add(Coord{-1, 0}), targetColor)...)
	result = append(result, f.extractGroupByColor(coord.Add(Coord{+1, 0}), targetColor)...)
	result = append(result, f.extractGroupByColor(coord.Add(Coord{0, -1}), targetColor)...)
	result = append(result, f.extractGroupByColor(coord.Add(Coord{0, +1}), targetColor)...)

	return result
}

func (f Field) getCell(c Coord) Cell {
	return f[c.i][c.j]
}

func (f Field) setColor(c Coord, value Color) {
	f[c.i][c.j].color = value
}

func (f Field) inRange(coord Coord) bool {
	inRange := (coord.i >= 0) && (coord.i < len(f)) &&
		(coord.j >= 0) && (coord.j < len(f[coord.i]))
	return inRange
}

func (c Coord) Add(coord Coord) Coord {
	c.i += coord.i
	c.j += coord.j
	return c
}

func (f Field) rect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	s := svg.New(w)
	s.Start(20*40, 20*40)
	for i, row := range f {
		for j, cell := range row {
			cellCode := fmt.Sprintf("stroke:black;fill:%s", cell.color)
			if cell.mark {
				cellCode += ";stroke-width:4"
			}
			s.Rect(i*20, j*12, 20, 12, cellCode)
		}
	}
	s.End()
}

