package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/JesusIslam/snigo"
)

type Sample struct {
	Source             []float64   `json:"source"`
	NegativeReferences [][]float64 `json:"negativeReferences"`
	PositiveReferences [][]float64 `json:"positiveReferences"`
}

func main() {
	var samplePath string
	var detectionReq int
	var sourceWindowSize int
	var gamma, theta float64
	flag.StringVar(&samplePath, "sample", "./sample.json", "Set path to sample file in JSON format")
	flag.IntVar(&detectionReq, "req", 3, "Set detection requirement, cannot be less than 1 or will be set to 3")
	flag.IntVar(&sourceWindowSize, "window", 5, "Set source window size, cannot be less than 1 or will be set to 5")
	flag.Float64Var(&gamma, "gamma", snigo.DefaultGamma, "Set gamma, will be defaulted to 1")
	flag.Float64Var(&theta, "theta", snigo.DefaultTheta, "Set theta, will be defaulted to 1")
	flag.Parse()

	if detectionReq < 1 {
		detectionReq = 3
	}
	if sourceWindowSize < 1 {
		sourceWindowSize = 5
	}

	raw, err := ioutil.ReadFile(samplePath)
	if err != nil {
		log.Println("Failed to read sample file:", err)
	}

	sample := &Sample{}
	err = json.Unmarshal(raw, sample)
	if err != nil {
		log.Fatalln("Failed to ")
	}

	window := make([]float64, sourceWindowSize)
	detected := 0
	previousDetection := false
	currentDetection := false
	for from := 0; from < (len(sample.Source) - sourceWindowSize); from++ {
		window = sample.Source[from:(from + sourceWindowSize)]
		currentDetection = snigo.Detect(window, sample.PositiveReferences, sample.NegativeReferences, gamma, theta)
		if currentDetection && previousDetection {
			detected++
		} else {
			detected = 0
		}
		previousDetection = currentDetection

		if detected == detectionReq {
			// trend detected
			log.Println("Trend detected at index %d with window size %d and detection requirement %d", from, sourceWindowSize, detectionReq)
			detected = 0
		}
	}

	log.Println("Detection completed")
}
