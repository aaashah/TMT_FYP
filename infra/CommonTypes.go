package infra

import (
	"math"

	"github.com/google/uuid"
)

type PositionVector struct {
	X int
	Y int
}

func (p1 PositionVector) Dist(p2 PositionVector) float32 {
	deltaX := float64(p1.X - p2.X)
	deltaY := float64(p1.Y - p2.Y)
	dist := math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))
	return float32(dist)
}

type Attachment struct {
	Anxiety   float32
	Avoidance float32
}

type Telomere struct {
	lowerThreshold int
	upperThreshold int
	lifespanDecay  float32
}

func NewTelomere(ageA, ageB int, lifeSpan float32) Telomere {
	return Telomere{ageA, ageB, lifeSpan}
}

func (t Telomere) GetProbabilityOfDeath(age int) float32 {
	if age < t.lowerThreshold {
		return 0.005 * float32(age) // Small increasing probability
	} else if age >= t.upperThreshold {
		return 1.0 // Guaranteed death at AgeB
	} else {
		// Linearly increasing probability from AgeA to AgeB
		return float32(age-t.lowerThreshold) / float32(t.upperThreshold-t.lowerThreshold)
	}
}

type SocialNetwork map[uuid.UUID]float32
