package tests

import (
	"testing"

	"github.com/aaashah/TMT_FYP/infra"
)

// same functionality as GetASMDecision in agents/ExtendedAgent.go (feed in ms, wv, rv, threshold to test in isolation)
func GetASMDecision(ms float32, wv float32, rv float32, threshold float32) infra.ASMDecison {
	thresholdScore := 0.0

	sum := 0
	for _, score := range []float32{ms, wv, rv} {
		if threshold > 0 {
			thresholdScore += min(float64(score/threshold), 1)
		} else {
			thresholdScore += 1
		}
		if score > threshold {
			sum += 1
		} else {
			sum -= 1
		}
	}

	if sum > 0 {
		return infra.SELF_SACRIFICE // Self-sacrifice
	} else if sum < 0 {
		return infra.NOT_SELF_SACRIFICE // Reject self-sacrifice
	} else {
		return infra.INACTION // No action
	}
}

func TestGetASMDecision(t *testing.T) {
	tests := []struct {
		name      string
		ms, wv, rv float32
		threshold  float32
		expected   infra.ASMDecison
	}{
		{
			name: "All above threshold -> self-sacrifice",
			ms: 0.8, wv: 0.9, rv: 0.85,
			threshold: 0.5,
			expected: infra.SELF_SACRIFICE,
		},
		{
			name: "All below threshold -> not self-sacrifice",
			ms: 0.3, wv: 0.2, rv: 0.1,
			threshold: 0.5,
			expected: infra.NOT_SELF_SACRIFICE,
		},
		{
			name: "Zero threshold -> should treat thresholdScore correctly and rely on sign of sum",
			ms: 0.4, wv: 0.6, rv: 0.8,
			threshold: 0.0,
			expected: infra.SELF_SACRIFICE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision := GetASMDecision(tt.ms, tt.wv, tt.rv, tt.threshold)
			if decision != tt.expected {
				t.Errorf("Expected %v, got %v (ms=%.2f, wv=%.2f, rv=%.2f, threshold=%.2f)", tt.expected, decision, tt.ms, tt.wv, tt.rv, tt.threshold)
			}
		})
	}
}