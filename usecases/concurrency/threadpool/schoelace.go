package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Point struct {
	x, y int
}

var pointsRegexp = regexp.MustCompile(`\((\d+),(\d+)\)`)

func main() {
	calcByOneThread()
	calcByTheadPool()
}

func calcByTheadPool() {
	const (
		workerNum = 8
	)
	polygonsFile := "polygons.txt"
	bs, _ := os.ReadFile(polygonsFile)
	text := string(bs)
	inputCh := make(chan string, 1000)
	var wg sync.WaitGroup
	wg.Add(workerNum)
	for i := 0; i < workerNum; i++ {
		go calcAreaByCh(inputCh, &wg)
	}
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		if len(line) == 0 {
			continue
		}
		inputCh <- line
	}
	close(inputCh)
	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)
}

func calcByOneThread() {
	polygonsFile := "polygons.txt"
	bs, _ := os.ReadFile(polygonsFile)
	text := string(bs)
	start := time.Now()
	for _, line := range strings.Split(text, "\n") {
		if len(line) == 0 {
			continue
		}
		calcArea(line)
	}
	elapsed := time.Since(start)
	fmt.Printf("Processing took %s \n", elapsed)
}

func calcArea(pointsStr string) {
	var points []Point
	for _, p := range pointsRegexp.FindAllStringSubmatch(pointsStr, -1) {
		x, _ := strconv.Atoi(p[1])
		y, _ := strconv.Atoi(p[2])
		points = append(points, Point{x, y})
	}
	area := 0.0
	prex, prey := 0, 0
	for _, point := range points {
		area += float64(point.x*prey) - float64(point.y*prex)
		prex, prey = point.x, point.y
	}
	// fmt.Println(math.Abs(area) / 2)
}

func calcAreaByCh(ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for pointsStr := range ch {
		calcArea(pointsStr)
	}
}
