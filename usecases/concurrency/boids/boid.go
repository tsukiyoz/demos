package main

import (
	"math"
	"math/rand"
	"time"
)

type Boid struct {
	id       int
	position Vector
	velocity Vector
}

func (b *Boid) run() {
	for {
		b.moveOne()
		time.Sleep(5 * time.Millisecond)
	}
}

func (b *Boid) calcAcceleration() Vector {
	avgAcceleration := Vector{x: 0, y: 0}
	avgPosition := Vector{x: 0, y: 0}
	separation := Vector{x: 0, y: 0}
	count := 0

	lock.RLock()
	for i := math.Max(b.position.x-viewRadius, float64(0)); i < math.Min(b.position.x+viewRadius, float64(screenWidth)); i++ {
		for j := math.Max(b.position.y-viewRadius, float64(0)); j < math.Min(b.position.y+viewRadius, float64(screenHeight)); j++ {
			if otherId := boidMap[int(i)][int(j)]; otherId != -1 && otherId != b.id {
				dist := boids[otherId].position.Distance(b.position)
				if dist < viewRadius {
					avgAcceleration = avgAcceleration.Add(boids[otherId].velocity)
					avgPosition = avgPosition.Add(boids[otherId].position)
					separation = separation.Add(b.position.Subtract(boids[otherId].position).DivideV(dist))
					count++
				}
			}
		}
	}
	lock.RUnlock()

	res := Vector{x: b.borderBounce(b.position.x, screenWidth), y: b.borderBounce(b.position.y, screenHeight)}
	if count > 0 {
		alignment := avgAcceleration.
			DivideV(float64(count)).
			Subtract(b.velocity).
			MultiplyV(adjRate)

		cohesion := avgPosition.
			DivideV(float64(count)).
			Subtract(b.position).
			MultiplyV(adjRate)

		separation = separation.
			MultiplyV(adjRate)

		res = res.Add(alignment).Add(cohesion).Add(separation)
	}
	return res
}

func (b *Boid) borderBounce(pos float64, border float64) float64 {
	if pos < viewRadius {
		return 1 / pos
	} else if pos > border-viewRadius {
		return 1 / (pos - border)
	}
	return 0
}

func (b *Boid) moveOne() {
	acceleration := b.calcAcceleration()

	lock.Lock()
	b.velocity = b.velocity.Add(acceleration).limit(-vMax, vMax)
	boidMap[int(b.position.x)][int(b.position.y)] = -1
	b.position = b.position.Add(b.velocity)
	if b.position.x < 0 {
		b.position.x = -b.position.x
		b.velocity.x = -b.velocity.x
	}
	if b.position.x >= screenWidth {
		b.position.x = 2*screenWidth - b.position.x
		b.velocity.x = -b.velocity.x
	}
	if b.position.y < 0 {
		b.position.y = -b.position.y
		b.velocity.y = -b.velocity.y
	}
	if b.position.y >= screenHeight {
		b.position.y = 2*screenHeight - b.position.y
		b.velocity.y = -b.velocity.y
	}
	boidMap[int(b.position.x)][int(b.position.y)] = b.id
	lock.Unlock()
}

func createBoid(id int) {
	boid := &Boid{
		id:       id,
		position: Vector{x: rand.Float64() * screenWidth, y: rand.Float64() * screenHeight},
		velocity: Vector{x: (rand.Float64() * vMax * 2) - vMax, y: (rand.Float64() * vMax * 2) - vMax},
	}

	if boid.position.x+boid.velocity.x < 0 {
		boid.velocity.x = -boid.velocity.x
	}
	if boid.position.x+boid.velocity.x > screenWidth {
		boid.velocity.x = screenWidth - boid.position.x
	}
	if boid.position.y+boid.velocity.y < 0 {
		boid.velocity.y = -boid.velocity.y
	}
	if boid.position.y+boid.velocity.y > screenHeight {
		boid.velocity.y = screenHeight - boid.position.y
	}

	boids[id] = boid
	boidMap[int(boid.position.x)][int(boid.position.y)] = id
	go boid.run()
}
