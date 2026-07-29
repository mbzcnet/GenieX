# Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
# SPDX-License-Identifier: BSD-3-Clause

"""``geniex.resolve_device_map`` aliases — mirrors ``geniex_resolve_device``.

Source of truth lives in ``sdk/src/device.cpp``. Any change to the alias
table there must update these tests in the same PR.

Platform-specific: on Linux ARM64 the ``gpu`` alias returns ``Vulkan0``;
on other platforms (Windows, Android) it returns ``GPUOpenCL``.
"""

from __future__ import annotations

import geniex


def test_auto_resolves_to_a_known_runtime(geniex_session):
    runtime, device_id, ngl = geniex.resolve_device_map('auto')
    assert runtime in geniex.get_runtime_list()
    assert device_id is None or isinstance(device_id, str)
    assert ngl is None or isinstance(ngl, int)


def test_cpu_alias_zeroes_gpu_layers(geniex_session):
    runtime, _, ngl = geniex.resolve_device_map('cpu')
    assert runtime == 'llama_cpp'
    assert ngl == 0


def test_hybrid_alias_offloads_all_layers(geniex_session):
    # resolve_device_map passes no explicit ngl, so the resolver returns -1
    # (all layers) which surfaces as None (no override).
    runtime, _, ngl = geniex.resolve_device_map('hybrid')
    assert runtime == 'llama_cpp'
    assert ngl is None


def test_llama_cpp_sdk_default_is_hybrid(geniex_session):
    """Bare runtime name uses SDK empty/auto default: hybrid for llama_cpp."""
    runtime, device_id, ngl = geniex.resolve_device_map('llama_cpp')
    assert runtime == 'llama_cpp'
    # hybrid → empty device_id → llama.cpp multi-backend scheduler
    assert device_id is None or device_id == ''
    assert ngl is None


def test_auto_applies_sdk_default_for_selected_runtime(geniex_session):
    """'auto' picks a registered runtime; compute unit is the SDK default."""
    runtime, device_id, ngl = geniex.resolve_device_map('auto')
    assert runtime in geniex.get_runtime_list()
    if runtime == 'llama_cpp':
        assert device_id is None or device_id == ''
        assert ngl is None
    elif runtime == 'qairt':
        assert device_id  # NPU
        assert ngl == 0


def test_llama_cpp_npu_alias_pins_htp0(geniex_session):
    runtime, device_id, ngl = geniex.resolve_device_map('llama_cpp:npu')
    assert runtime == 'llama_cpp'
    assert device_id == 'HTP0'
    assert ngl is None


def test_qairt_npu_alias_resolves_to_qairt(geniex_session):
    runtime, device_id, _ = geniex.resolve_device_map('qairt:npu')
    assert runtime == 'qairt'
    assert isinstance(device_id, str) and device_id


def test_qairt_default_is_npu(geniex_session):
    """Bare 'qairt' runtime uses SDK default: npu."""
    runtime, device_id, ngl = geniex.resolve_device_map('qairt')
    assert runtime == 'qairt'
    assert isinstance(device_id, str) and device_id
    assert ngl == 0


def test_gpu_alias_resolves_to_device(geniex_session):
    """gpu alias returns a non-null, non-empty device_id (platform-dependent)."""
    runtime, device_id, ngl = geniex.resolve_device_map('gpu')
    assert runtime == 'llama_cpp'
    assert isinstance(device_id, str) and device_id
    # Platform: "GPUOpenCL" (Win/Android) or "Vulkan0" (Linux ARM64)
    assert device_id in ('GPUOpenCL', 'Vulkan0')
    assert ngl is None
