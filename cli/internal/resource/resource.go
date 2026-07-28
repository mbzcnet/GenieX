// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

// Package resource resolves CLI resource-limit flags into concrete SDK knobs.
//
// Resolution stays in the CLI (not bindings/go fillC, not the native SDK) so:
//   - Default 0 leaves n_threads / n_gpu_layers unset and the SDK's device-aware
//     matrices apply (offload pins ~6 threads; pure CPU uses hardware_concurrency).
//   - Explicit percents become concrete values before they cross the C ABI —
//     no binding-only policy that Python/Android would miss.
package resource

import "runtime"

// ResolveNThreads maps cpu_percent → n_threads for geniex_ModelConfig.
//
//	percent <= 0  → 0 (SDK resolve_n_threads: device matrix)
//	percent 1-100 → max(1, NumCPU * percent / 100)
//	percent > 100 → treated as 100
func ResolveNThreads(cpuPercent int32) int32 {
	if cpuPercent <= 0 {
		return 0
	}
	pct := cpuPercent
	if pct > 100 {
		pct = 100
	}
	n := runtime.NumCPU() * int(pct) / 100
	if n < 1 {
		n = 1
	}
	return int32(n)
}

// ScaleGpuLayers scales a positive layer count by gpu_percent.
//
//	ngl <= 0 (-1 = all layers, 0 = CPU-only): unchanged — total layer count is
//	  only known after model load, so percent cannot scale "all layers".
//	percent <= 0 or >= 100: unchanged
//	ngl > 0 and percent 1-99: max(1, ngl * percent / 100)
func ScaleGpuLayers(ngl, gpuPercent int32) int32 {
	if ngl <= 0 || gpuPercent <= 0 || gpuPercent >= 100 {
		return ngl
	}
	scaled := int(ngl) * int(gpuPercent) / 100
	if scaled < 1 {
		scaled = 1
	}
	return int32(scaled)
}
