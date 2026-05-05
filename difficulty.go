// ABOUTME: DifficultyConfig and retargeting compute the proof-of-work
// ABOUTME: difficulty for each block based on observed block times.
package quark

type DifficultyConfig struct {
	InitialDifficulty int32
	RetargetInterval  int
	TargetBlockTime   int64
	MaxAdjustFactor   int64
}

func DefaultDifficultyConfig() *DifficultyConfig {
	return &DifficultyConfig{
		InitialDifficulty: 8,
		RetargetInterval:  0,
		TargetBlockTime:   1,
		MaxAdjustFactor:   4,
	}
}

func adjustDifficulty(current int32, target, actual, maxFactor int64) int32 {
	if actual <= 0 {
		actual = 1
	}
	low := target / maxFactor
	if low < 1 {
		low = 1
	}
	high := target * maxFactor
	if actual < low {
		actual = low
	}
	if actual > high {
		actual = high
	}
	if actual == target {
		return current
	}
	if actual < target {
		ratio := target / actual
		var delta int32
		for ratio >= 2 {
			delta++
			ratio /= 2
		}
		return current + delta
	}
	ratio := actual / target
	var delta int32
	for ratio >= 2 {
		delta++
		ratio /= 2
	}
	out := current - delta
	if out < 1 {
		out = 1
	}
	return out
}
