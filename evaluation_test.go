package main

import (
	"math/rand"
	"testing"
)

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Helper()
		b.Fatal("error:", err)
	}
}

func BenchmarkFitnessEvaluator(b *testing.B) {
	checkB(b, readConfig())
	rng := rand.New(rand.NewSource(99))
	want := 1.711917e+06
	failStr := "if rng has not been changed and is still rand.New(rand.NewSource(99)), want fitness %v, got %v"

	// generate random reference and candidate images
	factory, err := newImageDNAfactory(128, 128)
	ref := factory.GenerateRandomCandidate(rng).(*imageDNA).render()
	cand := factory.GenerateRandomCandidate(rng).(*imageDNA)

	// create the fitness evaluator
	evaluator := fitnessEvaluator{img: ref}
	checkB(b, err)

	b.ResetTimer()
	fitnesses := make([]float64, b.N)
	for n := 0; n < b.N; n++ {
		fitnesses[n] = evaluator.FitnessWorking(cand, nil)
	}
	b.StopTimer()
	got := fitnesses[len(fitnesses)-1]
	if got != want {
		b.Fatalf(failStr, want, got)
	}
}
