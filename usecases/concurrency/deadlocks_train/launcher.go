package main

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tsukiyoz/demos/usecases/concurrency/deadlocks_train/common"
	"github.com/tsukiyoz/demos/usecases/concurrency/deadlocks_train/hierarchy"
)

var (
	trains        [4]*common.Train
	intersections [4]*common.Intersection
)

const trainLength = 70

var _ ebiten.Game = (*Game)(nil)

type Game struct{}

func (g *Game) Draw(screen *ebiten.Image) {
	DrawTracks(screen)
	DrawIntersections(screen)
	DrawTrains(screen)
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return 320, 320
}

func (g *Game) Update() error {
	return nil
}

func main() {
	for i := 0; i < 4; i++ {
		trains[i] = &common.Train{
			ID:     i,
			Length: trainLength,
			Front:  0,
		}
	}

	for i := 0; i < 4; i++ {
		intersections[i] = &common.Intersection{
			ID:       i,
			Mutex:    sync.Mutex{},
			LockedBy: -1,
		}
	}

	go hierarchy.MoveTrain(trains[0], 300, []*common.Crossing{
		{Position: 125, Intersection: intersections[0]},
		{Position: 175, Intersection: intersections[1]},
	})

	go hierarchy.MoveTrain(trains[1], 300, []*common.Crossing{
		{Position: 125, Intersection: intersections[1]},
		{Position: 175, Intersection: intersections[2]},
	})

	go hierarchy.MoveTrain(trains[2], 300, []*common.Crossing{
		{Position: 125, Intersection: intersections[2]},
		{Position: 175, Intersection: intersections[3]},
	})

	go hierarchy.MoveTrain(trains[3], 300, []*common.Crossing{
		{Position: 125, Intersection: intersections[3]},
		{Position: 175, Intersection: intersections[0]},
	})

	ebiten.SetWindowSize(320*3, 320*3)
	ebiten.SetWindowTitle("Trains in a box")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
