// Copyright 2010 Petar Maymounkov. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package basicfile

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type (
	bfunc = func(b *testing.B)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randSample() any {
	return rand.Float64()
}

func genRandTestData(n int) avgVar {

	a := avgVar{}
	a.Init()

	for i := 0; i < n; i++ {
		go a.Add(randSample().(float64))
	}

	return a
}

var globalReturn any
var globalInt int
var globalPtr = &globalInt

func bRunWarmup(b *testing.B, name string, f func() any) {
	for i := 0; i < b.N; i++ {
		f()
	}
}

func bRunNoReturn(b *testing.B, name string, f func() any) {
	b.Run(name+"(no return)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			f()
		}
	})
}

func bRunLocal(b *testing.B, name string, f func() any) {
	b.Run(name+"(local)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = f()
		}
	})
}

func bRunGlobal(b *testing.B, name string, f func() any) {
	b.Run(name+"(global)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			globalReturn = f()
		}
	})
}

func bRunVarCheck(b *testing.B, name string, f func() any) {
	bRunWarmup(b, name, f)
	bRunNoReturn(b, name, f)
	bRunLocal(b, name, f)
	bRunGlobal(b, name, f)
}

func bRunWithWarmup(b *testing.B, name string, f func() any) {
	bRunWarmup(b, name, f)
	// bRunNoReturn(b, name, f)
	bRunLocal(b, name, f)
	// bRunGlobal(b, name, f)
}

func bRun(b *testing.B, name string, f func() any) {
	// bRunWarmup(b, name, f)
	// bRunNoReturn(b, name, f)
	bRunLocal(b, name, f)
	// bRunGlobal(b, name, f)
}

func BenchmarkRandSample(b *testing.B) {
	bRun(b, "RandSample", randSample)
}

func BenchmarkGenRandTestData(b *testing.B) {

	for i := 3; i < 12; i++ {
		n := 1 << i

		f := func() any { return genRandTestData(n) }
		bRun(b, fmt.Sprintf("genRandTestData(%d)", n), f)
	}
}

func Test_avgVar_GetStdDev(t *testing.T) {
	type fields struct {
		count int64
		sum   float64
		sumsq float64
	}
	tests := []struct {
		name          string
		fields        fields
		wantPrecision float64
	}{
		{"n=1000", fields{count: 1000, sum: 500, sumsq: 250000}, 0.1},
		{"n=1000", fields{count: 1000, sum: 500, sumsq: 250000}, 0.01},
		{"n=1000", fields{count: 1000, sum: 500, sumsq: 250000}, 0.001},
		{"n=1000", fields{count: 1000, sum: 500, sumsq: 250000}, 0.0001},
		{"n=1000", fields{count: 1000, sum: 500, sumsq: 250000}, 0.00001},

		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.1},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.01},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.001},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.0001},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.00001},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.000001},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.0000001},
		{"n=100000", fields{count: 100000, sum: 50000, sumsq: 2500000000}, 0.00000001},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			av := &avgVar{
				count: tt.fields.count,
				sum:   tt.fields.sum,
				sumsq: tt.fields.sumsq,
			}

			n := tt.fields.count
			ndiv2 := float64(n) / 2.0

			lowerbound := ndiv2 * (1 - tt.wantPrecision)
			upperbound := ndiv2 * (1 + tt.wantPrecision)

			avWant := &avgVar{n, ndiv2, ndiv2 * ndiv2}
			want := avWant.GetStdDev()

			avMin := &avgVar{n, lowerbound, lowerbound * lowerbound}
			avMax := &avgVar{n, upperbound, upperbound * upperbound}
			got := av.GetStdDev()
			wantMin := avMin.GetStdDev()
			wantMax := avMax.GetStdDev()
			if got < wantMin || got > wantMax {
				t.Errorf("avgVar.GetStdDev() outside of spec (precision: %v, min:%v max:%v)= %v, want %v", tt.wantPrecision, wantMin, wantMax, got, want)
			}
		})
	}
}
