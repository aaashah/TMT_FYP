package infra

import (
	"math"

	"github.com/google/uuid"
)

type PositionVector struct {
	X int
	Y int
}

func (p1 PositionVector) Dist(p2 PositionVector) float64 {
	deltaX := float64(p1.X - p2.X)
	deltaY := float64(p1.Y - p2.Y)
	return math.Sqrt((deltaX * deltaX) + (deltaY * deltaY))
}

type Attachment struct {
	Anxiety   float32
	Avoidance float32
}

type Telomere struct {
	age            int
	lowerThreshold int
	upperThreshold int
	lifespanDecay  float32
}

type PTSParams struct {
	CheckProb    float32 // prob of sending a wellbeing check
	ReplyProb    float32 // prob of replying to a check
	Alpha        float32 // reinforcement param
	Beta         float32 // reinforcement param
}

func NewTelomere(age, ageA, ageB int, lifeSpan float32) *Telomere {
	return &Telomere{age, ageA, ageB, lifeSpan}
}

func (t *Telomere) GetAge() int {
	return t.age
}

func (t *Telomere) IncrementAge() {
	t.age++
}

func (t *Telomere) GetProbabilityOfDeath() float32 {
	if t.age < t.lowerThreshold {
		return 0.005 * float32(t.age) // Small increasing probability
	} else if t.age >= t.upperThreshold {
		return 1.0 // Guaranteed death at AgeB
	} else {
		// Linearly increasing probability from AgeA to AgeB
		return float32(t.age-t.lowerThreshold) / float32(t.upperThreshold-t.lowerThreshold)
	}
}

type SocialNetwork map[uuid.UUID]float32

type Ysterofimia struct {
	SelfSacrificeCount		    int // number of times agent volunteered self-sacrifices
	SelfSacrificeEsteems        float32 // total esteems from agents who self-sacrificed
	OtherEliminationCount       int // number of times agent was eliminated by other than self-sacrifice
	OtherEliminationsEsteems    float32 // total esteems from agents who eliminated other than self-sacrifice
}

func NewYsterofimia() *Ysterofimia {
	return &Ysterofimia{
		SelfSacrificeCount:         0,
		SelfSacrificeEsteems:       0.0,
		OtherEliminationCount:      0,
		OtherEliminationsEsteems:   0.0,
	}
}
func (y *Ysterofimia) IncrementSelfSacrificeCount() {
	y.SelfSacrificeCount++
}

func (y *Ysterofimia) AddSelfSacrificeEsteems(esteem float32) {
	y.SelfSacrificeEsteems += esteem
}
func (y *Ysterofimia) IncrementOtherEliminationCount() {
	y.OtherEliminationCount++
}
func (y *Ysterofimia) AddOtherEliminationsEsteems(esteem float32) {
	y.OtherEliminationsEsteems += esteem
}

func (y *Ysterofimia) ComputeYsterofimia() float32 {
	totalEliminations := y.SelfSacrificeCount + y.OtherEliminationCount
	totalEsteem := y.SelfSacrificeEsteems + y.OtherEliminationsEsteems

	if totalEliminations == 0 || totalEsteem == 0 {
        return 0.5 // no eliminations or no esteem data
    }

	//esteemRatio := float32(y.SelfSacrificeEsteems) / float32(totalEsteem)
	sacrificeEsteemRatio := float32(y.SelfSacrificeEsteems) / float32(y.SelfSacrificeCount)
	otherEsteemRatio := float32(y.OtherEliminationsEsteems) / float32(y.OtherEliminationCount)

	if sacrificeEsteemRatio > otherEsteemRatio {
		return float32(y.SelfSacrificeCount) / float32(totalEliminations)
	} else {
		return float32(y.OtherEliminationCount) / float32(totalEliminations)
	}
}

type ProximityArray []float32

func (pa ProximityArray) MapToRelativeProximities() ProximityArray {
	var totInvProx float32 = 0.0
	for _, prox := range pa {
		totInvProx += 1 / prox
	}
	for idx, prox := range pa {
		invProx := 1 / prox
		pa[idx] = invProx / totInvProx
	}
	return pa
}