package stats

func BayesianShrinkByVarianceMatching(rate []float64, n []float64) (res []float64, k float64, mean float64) {
	mean = Mean(rate)
	k = VarianceMatching(rate, n, mean)
	for i, p := range rate {
		tmp := bayesianShrink(p, n[i], mean, k)
		res = append(res, tmp)
	}
	return res, k, mean
}

func VarianceMatching(p []float64, n []float64, mean float64) (k float64) {
	avgSamplingNoise := AvgSamplingNoise(p, n)
	totalVariance := Variance(p, mean)
	skillSpread := totalVariance - avgSamplingNoise

	if skillSpread <= 0 {
		return 1e9
	}
	k = mean*(1-mean)/skillSpread - 1
	return
}

func bayesianShrink(p float64, n float64, mean, k float64) (adjustedP float64) {
	adjustedP = (p*n + mean*k) / (n + k)
	return
}
