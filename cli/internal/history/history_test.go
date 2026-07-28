// Copyright 2024-2026 Qualcomm Technologies, Inc. and/or its subsidiaries.
// SPDX-License-Identifier: BSD-3-Clause

package history

import (
	"testing"

	geniex_sdk "github.com/qualcomm/GenieX/bindings/go"
)

func TestTrimLlmHistory_KeepsSystemAndLastTurns(t *testing.T) {
	history := []geniex_sdk.LlmChatMessage{
		{Role: geniex_sdk.LlmRoleSystem, Content: "sys"},
		{Role: geniex_sdk.LlmRoleUser, Content: "u1"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a1"},
		{Role: geniex_sdk.LlmRoleUser, Content: "u2"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a2"},
		{Role: geniex_sdk.LlmRoleUser, Content: "u3"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a3"},
	}
	got := TrimLlmHistory(history, 2)
	if len(got) != 5 { // system + 4 (2 turns)
		t.Fatalf("len=%d want 5: %+v", len(got), got)
	}
	if got[0].Content != "sys" {
		t.Fatalf("system lost: %+v", got[0])
	}
	if got[1].Content != "u2" || got[4].Content != "a3" {
		t.Fatalf("unexpected window: %+v", got)
	}
}

func TestTrimLlmHistory_NoTrimWhenDisabled(t *testing.T) {
	history := []geniex_sdk.LlmChatMessage{
		{Role: geniex_sdk.LlmRoleUser, Content: "u1"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a1"},
	}
	got := TrimLlmHistory(history, 0)
	if len(got) != 2 {
		t.Fatalf("len=%d want 2", len(got))
	}
}

func TestTrimLlmHistory_DropsLeadingAssistantAfterCut(t *testing.T) {
	// Odd-length rest so cut lands on assistant first.
	history := []geniex_sdk.LlmChatMessage{
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a0"},
		{Role: geniex_sdk.LlmRoleUser, Content: "u1"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a1"},
		{Role: geniex_sdk.LlmRoleUser, Content: "u2"},
		{Role: geniex_sdk.LlmRoleAssistant, Content: "a2"},
	}
	got := TrimLlmHistory(history, 1) // keep 2 msgs from end, may drop leading assistant
	if len(got) == 0 {
		t.Fatal("empty result")
	}
	if got[0].Role == geniex_sdk.LlmRoleAssistant && got[0].Content == "a0" {
		t.Fatalf("should not keep discarded-prefix assistant: %+v", got)
	}
}
