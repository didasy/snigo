// A naive Snikolov trend-detection algorithm implementation.
package snigo

import (
	"math"
)

const (
	DefaultGamma                = 1
	DefaultTheta                = 1
	DefaultDetectionRequirement = 3
)

// Detect source window using positive and negative reference signals for a trend.
// You should loop this function in your worker until detection requirement is met.
// So if your detection requirement is 3, and this function return `true` 3 times in a row,
// then a trend is detected.
func Detect(sourceWindow []float64, positiveReferenceSignals, negativeReferenceSignals [][]float64, gamma float64, theta float64) bool {

	negativeDistances := []float64{}
	positiveDistances := []float64{}

	for _, positiveReference := range positiveReferenceSignals {
		distanceToPositiveReference := DistanceToReference(sourceWindow, positiveReference)
		positiveDistances = append(positiveDistances, distanceToPositiveReference)
	}

	for _, negativeReference := range negativeReferenceSignals {
		distanceToNegativeReference := DistanceToReference(sourceWindow, negativeReference)
		negativeDistances = append(negativeDistances, distanceToNegativeReference)
	}

	ratio := ProbabilityClass(positiveDistances, gamma) / ProbabilityClass(negativeDistances, gamma)
	if ratio > theta {
		return true
	}

	return false
}

// Compute the minimum distance between `source` and all pieces of `reference` of the same length as `source`.
func DistanceToReference(source, reference []float64) float64 {
	nOfObservation := len(source)
	nOfReference := len(reference)
	minimumDistance := math.Inf(1)

	for i := 0; i < (nOfReference - nOfObservation + 1); i++ {
		to := i + nOfObservation - 1
		distance := Distance(reference[i:to], source)
		minimumDistance = math.Min(minimumDistance, distance)
	}

	return minimumDistance
}

// Compute Euclidean distance between two signals `source` and `target` with same length.
func Distance(source, target []float64) float64 {
	distance := 0.0

	for i := 0; i < len(source); i++ {
		distance = distance + math.Pow(source[i]-target[i], 2)
	}

	return distance
}

// Using the distance of an observation to the reference signals of a certain class,
// compute a number proportional to the probability that the observation belongs to that class.
func ProbabilityClass(distances []float64, gamma float64) float64 {
	probability := 0.0

	for _, distance := range distances {
		probability = probability + math.Exp(-1*gamma*distance)
	}

	return probability
}
