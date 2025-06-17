package stats

import (
	"math"
)

func Mean(data []float64) (mean float64) {
	sum := 0.0
	for _, num := range data {
		sum = sum + num
	}
	mean = sum / float64(len(data))
	return
}

func AvgSamplingNoise(p []float64, n []float64) (avgSamplingNoise float64) {
	var sum float64
	for i, prob := range p {
		sum = sum + samplingNoise(prob, n[i])
	}
	avgSamplingNoise = sum / float64(len(p))
	return
}

func samplingNoise(p float64, n float64) (samplingNoise float64) {
	samplingNoise = p * (1.0 - p) / float64(n)
	return
}

func Variance(p []float64, mean float64) (variance float64) {
	var sum float64
	for _, prob := range p {
		sum = sum + math.Pow(prob-mean, 2)
	}
	variance = sum / float64(len(p))
	return
}
