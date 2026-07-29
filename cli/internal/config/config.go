// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	// Global settings
	DataDir string

	// Server settings
	Host      string // Server host and port (default: "127.0.0.1:18181")
	Origins   string // Allowed CORS origins (default: "*")
	KeepAlive int64  // Connection keep-alive timeout in seconds (default: 300)
	// Model-load defaults applied when a request omits them (llama_cpp only;
	// per-request body fields still override). Compute is the alias resolved by
	// the SDK (sdk/src/device.cpp); empty means the SDK's own default.
	NCtx       int32  // Default context window size (default: 16384)
	Ngl        int32  // Default GPU/NPU layers to offload, -1 = all (default: -1)
	Compute    string // Default compute unit; empty = SDK default (hybrid@llama_cpp, npu@qairt)
	KvCache    string // Convenience: default type for both K and V (empty = auto)
	CacheTypeK string // Default KV cache type for K (empty = fall back to KvCache / auto)
	CacheTypeV string // Default KV cache type for V (empty = fall back to KvCache / auto)
	// CpuPercent / GpuPercent: 0 = leave to SDK (no override). 1-100 caps threads
	// or scales an explicit positive --ngl. See cli/internal/resource.
	CpuPercent int32
	GpuPercent int32

	// HTTPS / TLS settings
	HTTPS    bool   // Whether to serve over HTTPS (default: false)
	CertFile string // TLS certificate file path
	KeyFile  string // TLS private key file path

	// Env only params
	HFToken string
	Log     string
}

// init sets up viper defaults and env binding. Runs once at package load.
func init() {
	// ENV only param need to set default here
	viper.SetDefault("hftoken", "") // Default empty token
	viper.SetDefault("log", "none") // Default log level

	viper.SetEnvPrefix("geniex")
	viper.AutomaticEnv()
}

// Get returns a fresh snapshot of the current viper state. Unmarshalling on
// every call is intentional: cobra populates subcommand flags in stages, so
// callers at different points in startup observe different values. Sharing a
// cached *Config via sync.Once would freeze whichever snapshot was observed
// first — usually before the subcommand's own flags are visible.
func Get() *Config {
	c := &Config{}
	viper.Unmarshal(c)
	c.HFToken = resolveHFToken(c.HFToken)
	return c
}

// ResolvedCacheTypes merges KvCache into CacheTypeK/V when either side is empty.
// Empty result means the SDK auto policy (q8_0 when n_ctx >= 8192).
func (c *Config) ResolvedCacheTypes() (k, v string) {
	k, v = c.CacheTypeK, c.CacheTypeV
	if c.KvCache != "" {
		if k == "" {
			k = c.KvCache
		}
		if v == "" {
			v = c.KvCache
		}
	}
	return k, v
}

func resolveHFToken(geniexToken string) string {
	if geniexToken != "" {
		return geniexToken
	}

	if token := os.Getenv("HF_TOKEN"); token != "" {
		return token
	}

	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		return ""
	}

	tokenPath := filepath.Join(homeDir, ".cache", "huggingface", "token")
	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
