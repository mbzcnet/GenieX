// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolvedCacheTypes_MergesKvCache(t *testing.T) {
	c := &Config{KvCache: "q8_0"}
	k, v := c.ResolvedCacheTypes()
	if k != "q8_0" || v != "q8_0" {
		t.Fatalf("got %q/%q, want q8_0/q8_0", k, v)
	}

	c = &Config{KvCache: "q8_0", CacheTypeK: "f16"}
	k, v = c.ResolvedCacheTypes()
	if k != "f16" || v != "q8_0" {
		t.Fatalf("got %q/%q, want f16/q8_0 (K override)", k, v)
	}

	c = &Config{}
	k, v = c.ResolvedCacheTypes()
	if k != "" || v != "" {
		t.Fatalf("got %q/%q, want empty auto", k, v)
	}
}

func TestResolveHFToken_GenieXPreferredOverHFToken(t *testing.T) {
	t.Setenv("HF_TOKEN", "hf_token")

	got := resolveHFToken("geniex_token")
	if got != "geniex_token" {
		t.Fatalf("expected GENIEX token, got %q", got)
	}
}

func TestResolveHFToken_UsesHFTokenWhenGenieXMissing(t *testing.T) {
	t.Setenv("HF_TOKEN", "hf_token")

	got := resolveHFToken("")
	if got != "hf_token" {
		t.Fatalf("expected HF_TOKEN value, got %q", got)
	}
}

func TestResolveHFToken_UsesDefaultTokenFileWhenEnvMissing(t *testing.T) {
	home := t.TempDir()
	tokenPath := filepath.Join(home, ".cache", "huggingface", "token")
	if err := os.MkdirAll(filepath.Dir(tokenPath), 0o755); err != nil {
		t.Fatalf("failed to create token dir: %v", err)
	}
	if err := os.WriteFile(tokenPath, []byte("  file_token\n"), 0o644); err != nil {
		t.Fatalf("failed to write token file: %v", err)
	}

	t.Setenv("HF_TOKEN", "")
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	got := resolveHFToken("")
	if got != "file_token" {
		t.Fatalf("expected token from file fallback, got %q", got)
	}
}

func TestResolveHFToken_EmptyWhenNoSourceAvailable(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HF_TOKEN", "")
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)

	got := resolveHFToken("")
	if got != "" {
		t.Fatalf("expected empty token when no source is available, got %q", got)
	}
}

// TestGet_GenieXHFTokenBeatsHFToken verifies the full precedence chain through
// Get(): GENIEX_HFTOKEN (via viper) must beat HF_TOKEN. This guards against
// viper config drift silently breaking the documented order.
func TestGet_GenieXHFTokenBeatsHFToken(t *testing.T) {
	t.Setenv("GENIEX_HFTOKEN", "geniex_value")
	t.Setenv("HF_TOKEN", "hf_value")

	if got := Get().HFToken; got != "geniex_value" {
		t.Fatalf("expected geniex_value (GENIEX_HFTOKEN should win), got %q", got)
	}
}

// TestGet_HFTokenUsedWhenGenieXUnset verifies that Get() falls back to
// HF_TOKEN when GENIEX_HFTOKEN is unset, exercising the env path through viper.
func TestGet_HFTokenUsedWhenGenieXUnset(t *testing.T) {
	t.Setenv("GENIEX_HFTOKEN", "")
	t.Setenv("HF_TOKEN", "hf_value")

	if got := Get().HFToken; got != "hf_value" {
		t.Fatalf("expected hf_value via HF_TOKEN fallback, got %q", got)
	}
}
