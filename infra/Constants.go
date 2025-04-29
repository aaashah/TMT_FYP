package infra

type ASPDecison int

const (
	SELF_SACRIFICE ASPDecison = iota
	NOT_SELF_SACRIFICE
	INACTION
)

const (
	// ASP weights
	W1  float32 = 0.25
	W2  float32 = 0.25
	W3  float32 = 0.25
	W4  float32 = 0.25
	W5  float32 = 0.25
	W6  float32 = 0.25
	W7  float32 = 0.5
	W8  float32 = 0.25
	W9  float32 = 0.25
	W10 float32 = 0.5
)

const (
	GRID_WIDTH  int = 70
	GRID_HEIGHT int = 30
)
