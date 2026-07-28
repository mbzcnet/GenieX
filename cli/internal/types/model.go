// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package types

// ModelParam holds the model-load knobs the keep-alive cache keys instances on.
// NCtx / NGpuLayers / NThreads / CacheTypeK / CacheTypeV are llama_cpp-only;
// DeviceID is the compute unit resolved by the SDK (empty = the SDK's own default).
// NThreads 0 = SDK device matrix (offload ~6, pure CPU = hardware_concurrency).
type ModelParam struct {
	NCtx       int32
	NGpuLayers int32
	NThreads   int32 // 0 = SDK auto; set via CLI --cpu-percent
	DeviceID   string
	// CacheTypeK / CacheTypeV: empty = SDK auto (q8_0 when n_ctx >= 8192).
	CacheTypeK string
	CacheTypeV string
}
