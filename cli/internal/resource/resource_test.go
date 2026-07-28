// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package resource

import (
	"runtime"
	"testing"
)

func TestResolveNThreads_ZeroMeansSDKAuto(t *testing.T) {
	if got := ResolveNThreads(0); got != 0 {
		t.Errorf("ResolveNThreads(0) = %d, want 0 (SDK auto)", got)
	}
	if got := ResolveNThreads(-1); got != 0 {
		t.Errorf("ResolveNThreads(-1) = %d, want 0", got)
	}
}

func TestResolveNThreads_Percent(t *testing.T) {
	nCPU := runtime.NumCPU()
	got := ResolveNThreads(50)
	want := int32(max(1, nCPU*50/100))
	if got != want {
		t.Errorf("ResolveNThreads(50) = %d, want %d (50%% of %d cores)", got, want, nCPU)
	}
	if got := ResolveNThreads(100); got != int32(nCPU) && !(nCPU == 0 && got == 1) {
		// NumCPU can be 0 in exotic environments; Resolve clamps to >= 1 only when pct>0.
		if nCPU > 0 && got != int32(nCPU) {
			t.Errorf("ResolveNThreads(100) = %d, want %d", got, nCPU)
		}
	}
	if got := ResolveNThreads(200); got != int32(nCPU) && nCPU > 0 {
		t.Errorf("ResolveNThreads(200) = %d, want %d (clamped to 100%%)", got, nCPU)
	}
}

func TestScaleGpuLayers(t *testing.T) {
	cases := []struct {
		ngl, pct, want int32
	}{
		{-1, 90, -1}, // all layers: percent not applied
		{0, 90, 0},   // CPU-only: percent not applied
		{20, 0, 20},  // percent unset
		{20, 100, 20},
		{20, 90, 18},
		{20, 50, 10},
		{1, 10, 1}, // floor at 1
		{5, 10, 1},
	}
	for _, tc := range cases {
		if got := ScaleGpuLayers(tc.ngl, tc.pct); got != tc.want {
			t.Errorf("ScaleGpuLayers(%d, %d) = %d, want %d", tc.ngl, tc.pct, got, tc.want)
		}
	}
}
