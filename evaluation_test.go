package main

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"testing"

	"github.com/aurelien-rainone/evolve/framework"
)

func checkB(b *testing.B, err error) {
	if err != nil {
		b.Helper()
		b.Fatal("error:", err)
	}
}

// create an imageDNA for testing purposes, with 50 randomly generated polygons
// of the same color.
func createTestCandidate(r, g, b uint8) *imageDNA {
	const numPolys = 50
	rng := rand.New(rand.NewSource(0))
	var img = &imageDNA{
		w:     128,
		h:     128,
		polys: make([]poly, numPolys),
	}
	// add N `numPolys` random polygons
	for i := 0; i < numPolys; i++ {
		img.polys[i] = randomPoly(img, 3, 4, rng)
		img.polys[i].col = color.RGBA{r, g, b, 255}
	}
	return img
}

func BenchmarkCairoFitnessEvaluator(b *testing.B) {
	want := 0.0

	refImageFn := "testdata/red.png"
	cand := monochromeImage(128, 128, color.RGBA{255, 0, 0, 255})

	evaluator := newCairoEvaluator(refImageFn)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		b.StartTimer()
		got := evaluator.Fitness(cand, nil)
		b.StopTimer()

		if got != want {
			b.Fatalf("wrong fitness, want %v, got %v", want, got)
		}
	}
}

func monochromeImage(w, h int, col color.RGBA) *imageDNA {
	return &imageDNA{
		w: w,
		h: h,
		polys: []poly{
			poly{
				col: col,
				pts: []image.Point{
					image.Pt(0, 0),
					image.Pt(w, 0),
					image.Pt(w, h),
					image.Pt(0, h),
				},
			},
		},
	}
}

func loadPNGAsRGBA(fn string) (*image.RGBA, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return convertToRGBA(img), nil
}

func TestCairoEvaluator(t *testing.T) {
	type args struct {
		cand framework.Candidate
		pop  []framework.Candidate
	}
	tests := []struct {
		name   string
		refImg string
		args   args
		want   float64
	}{
		{
			"diff red with red",
			"./testdata/red.png",
			args{monochromeImage(128, 128, color.RGBA{255, 0, 0, 255}), nil},
			0.0,
		},
		{
			"diff green with green",
			"./testdata/green.png",
			args{monochromeImage(128, 128, color.RGBA{0, 255, 0, 255}), nil},
			0.0,
		},
		{
			"diff blue with blue",
			"./testdata/blue.png",
			args{monochromeImage(128, 128, color.RGBA{0, 0, 255, 255}), nil},
			0.0,
		},
		{
			"diff black with white",
			"./testdata/black.png",
			args{monochromeImage(128, 128, color.RGBA{255, 255, 255, 255}), nil},
			3 * 255 * 128 * 128,
		},
		{
			"diff white with black",
			"./testdata/white.png",
			args{monochromeImage(128, 128, color.RGBA{0, 0, 0, 255}), nil},
			3 * 255 * 128 * 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := newCairoEvaluator(tt.refImg)
			if got := ev.Fitness(tt.args.cand, tt.args.pop); got != tt.want {
				t.Errorf("cairoEvaluator.Fitness() = %v, want %v", got, tt.want)
			}
		})
	}
}
