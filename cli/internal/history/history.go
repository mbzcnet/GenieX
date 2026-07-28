// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package history

import geniex_sdk "github.com/qualcomm/GenieX/bindings/go"

// TrimLlmHistory keeps all system messages and the most recent maxTurns
// conversation turns (each turn ≈ one user message + one assistant reply).
// maxTurns <= 0 disables trimming.
func TrimLlmHistory(history []geniex_sdk.LlmChatMessage, maxTurns int) []geniex_sdk.LlmChatMessage {
	if maxTurns <= 0 || len(history) == 0 {
		return history
	}
	var systems, rest []geniex_sdk.LlmChatMessage
	for _, m := range history {
		if m.Role == geniex_sdk.LlmRoleSystem {
			systems = append(systems, m)
		} else {
			rest = append(rest, m)
		}
	}
	keep := maxTurns * 2
	if len(rest) <= keep {
		return history
	}
	rest = rest[len(rest)-keep:]
	// Avoid starting mid-turn on a bare assistant message.
	if len(rest) > 0 && rest[0].Role == geniex_sdk.LlmRoleAssistant {
		rest = rest[1:]
	}
	out := make([]geniex_sdk.LlmChatMessage, 0, len(systems)+len(rest))
	out = append(out, systems...)
	out = append(out, rest...)
	return out
}

// TrimVlmHistory is the VLM equivalent of TrimLlmHistory.
func TrimVlmHistory(history []geniex_sdk.VlmChatMessage, maxTurns int) []geniex_sdk.VlmChatMessage {
	if maxTurns <= 0 || len(history) == 0 {
		return history
	}
	var systems, rest []geniex_sdk.VlmChatMessage
	for _, m := range history {
		if m.Role == geniex_sdk.VlmRoleSystem {
			systems = append(systems, m)
		} else {
			rest = append(rest, m)
		}
	}
	keep := maxTurns * 2
	if len(rest) <= keep {
		return history
	}
	rest = rest[len(rest)-keep:]
	if len(rest) > 0 && rest[0].Role == geniex_sdk.VlmRoleAssistant {
		rest = rest[1:]
	}
	out := make([]geniex_sdk.VlmChatMessage, 0, len(systems)+len(rest))
	out = append(out, systems...)
	out = append(out, rest...)
	return out
}
