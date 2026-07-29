// Copyright (c) 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

// Single source of truth for the user-facing device alias table
// (cpu / gpu / npu / hybrid → concrete device_id + n_gpu_layers).
// Language bindings (Go CLI, Python, Android/JNI) call through to this.

#include <algorithm>
#include <cctype>
#include <cstdlib>
#include <string>

#include "geniex.h"
#include "logging.h"

#if defined(_WIN32)
#define portable_strdup _strdup
#else
#define portable_strdup strdup
#endif

namespace {

constexpr const char* kPluginQairt = "qairt";

constexpr const char* kAliasCPU    = "cpu";
constexpr const char* kAliasGPU    = "gpu";
constexpr const char* kAliasNPU    = "npu";
constexpr const char* kAliasHybrid = "hybrid";
constexpr const char* kAliasAuto   = "auto";

constexpr const char* kDeviceHTP0      = "HTP0";
#if defined(__linux__) && defined(__aarch64__) && !defined(__ANDROID__)
// On Linux ARM64 (Snapdragon X Elite), Vulkan is the primary GPU compute
// backend (Mesa Turnip / freedreno). OpenCL via rusticl may not enumerate
// any devices on this platform, so we map the "gpu" alias to the ggml-vulkan
// device name instead of "GPUOpenCL". Android uses OpenCL/Adreno and stays
// with "GPUOpenCL".
constexpr const char* kDeviceGPU = "Vulkan0";
#else
constexpr const char* kDeviceGPU = "GPUOpenCL";
#endif

constexpr const char* kDeviceQairtNPU  = "NPU";

std::string to_lower(const char* s) {
    if (!s) return {};
    std::string out(s);
    std::transform(
        out.begin(), out.end(), out.begin(), [](unsigned char c) { return static_cast<char>(std::tolower(c)); });
    return out;
}

std::string to_lower_trim(const char* s) {
    std::string out   = to_lower(s);
    size_t      start = 0;
    while (start < out.size() && std::isspace(static_cast<unsigned char>(out[start]))) ++start;
    size_t end = out.size();
    while (end > start && std::isspace(static_cast<unsigned char>(out[end - 1]))) --end;
    return out.substr(start, end - start);
}

bool is_known_alias(const std::string& s) {
    return s == kAliasCPU || s == kAliasGPU || s == kAliasNPU || s == kAliasHybrid;
}

}  // namespace

int32_t geniex_resolve_device(const geniex_ResolveDeviceInput* input, geniex_ResolveDeviceOutput* output) {
    if (!input || !output) {
        GENIEX_LOG_ERROR("geniex_resolve_device: input/output is null");
        return GENIEX_ERROR_COMMON_INVALID_INPUT;
    }

    // Initialise output so partial failures leave a sane state.
    output->device_id = nullptr;
    output->ngl       = input->ngl_default;
    output->warning   = nullptr;

    if (!input->plugin_id) {
        GENIEX_LOG_ERROR("geniex_resolve_device: plugin_id is null");
        return GENIEX_ERROR_COMMON_INVALID_INPUT;
    }

    const std::string plugin = input->plugin_id;
    std::string       alias  = to_lower_trim(input->mode);

    if (!alias.empty() && alias != kAliasAuto && !is_known_alias(alias)) {
        GENIEX_LOG_ERROR("geniex_resolve_device: invalid device mode '{}'", alias);
        return GENIEX_ERROR_COMMON_INVALID_DEVICE;
    }

    // Empty / "auto" → plugin-specific product default for every surface
    // (CLI, Python, Android, server). Bindings must not re-default.
    //
    //   llama_cpp → hybrid (empty device_id): per-tensor HTP+CPU scheduler —
    //     the general Snapdragon fast path when Hexagon is present; ops the
    //     NPU cannot run stay on CPU without pinning a single HTP session.
    //   qairt    → npu: QAIRT only has the Hexagon path.
    //
    // Explicit aliases (cpu / gpu / npu / hybrid) always win when supplied.
    if (alias.empty() || alias == kAliasAuto) {
        alias = (plugin == kPluginQairt) ? kAliasNPU : kAliasHybrid;
    }

    // QAIRT is NPU-only and rejects any non-zero n_gpu_layers, so force
    // ngl to 0. Non-npu aliases are coerced with a warning, not an error.
    if (plugin == kPluginQairt) {
        if (alias != kAliasNPU) {
            std::string msg =
                "qairt plugin only supports NPU inference; ignoring device='" + alias + "' and running on NPU";
            output->warning = portable_strdup(msg.c_str());
        }
        output->device_id = portable_strdup(kDeviceQairtNPU);
        output->ngl       = 0;
        return GENIEX_SUCCESS;
    }

    // llama_cpp: ngl passes through unchanged (-1 means "all layers" to
    // llama.cpp). Only cpu forces it to 0. hybrid leaves device_id empty so
    // llama.cpp's multi-backend scheduler picks HTP/CPU (and GPU if registered)
    // per tensor — the same path on Windows / Linux / Android Snapdragon.
    if (alias == kAliasCPU) {
        output->ngl = 0;
    } else if (alias == kAliasGPU) {
        output->device_id = portable_strdup(kDeviceGPU);
    } else if (alias == kAliasNPU) {
        output->device_id = portable_strdup(kDeviceHTP0);
    }
    // hybrid / residual auto: device_id stays null → hybrid scheduler.
    return GENIEX_SUCCESS;
}
