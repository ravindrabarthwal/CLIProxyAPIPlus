package config

import "testing"

func TestSanitizeOAuthModelAlias_PreservesForkFlag(t *testing.T) {
	cfg := &Config{
		OAuthModelAlias: map[string][]OAuthModelAlias{
			" CoDeX ": {
				{Name: " gpt-5 ", Alias: " g5 ", Fork: true},
				{Name: "gpt-6", Alias: "g6"},
			},
		},
	}

	cfg.SanitizeOAuthModelAlias()

	aliases := cfg.OAuthModelAlias["codex"]
	if len(aliases) != 2 {
		t.Fatalf("expected 2 sanitized aliases, got %d", len(aliases))
	}
	if aliases[0].Name != "gpt-5" || aliases[0].Alias != "g5" || !aliases[0].Fork {
		t.Fatalf("expected first alias to be gpt-5->g5 fork=true, got name=%q alias=%q fork=%v", aliases[0].Name, aliases[0].Alias, aliases[0].Fork)
	}
	if aliases[1].Name != "gpt-6" || aliases[1].Alias != "g6" || aliases[1].Fork {
		t.Fatalf("expected second alias to be gpt-6->g6 fork=false, got name=%q alias=%q fork=%v", aliases[1].Name, aliases[1].Alias, aliases[1].Fork)
	}
}

func TestSanitizeOAuthModelAlias_AllowsMultipleAliasesForSameName(t *testing.T) {
	cfg := &Config{
		OAuthModelAlias: map[string][]OAuthModelAlias{
			"antigravity": {
				{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5-20251101", Fork: true},
				{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5-20251101-thinking", Fork: true},
				{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5", Fork: true},
			},
		},
	}

	cfg.SanitizeOAuthModelAlias()

	aliases := cfg.OAuthModelAlias["antigravity"]
	expected := []OAuthModelAlias{
		{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5-20251101", Fork: true},
		{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5-20251101-thinking", Fork: true},
		{Name: "gemini-claude-opus-4-5-thinking", Alias: "claude-opus-4-5", Fork: true},
	}
	if len(aliases) != len(expected) {
		t.Fatalf("expected %d sanitized aliases, got %d", len(expected), len(aliases))
	}
	for i, exp := range expected {
		if aliases[i].Name != exp.Name || aliases[i].Alias != exp.Alias || aliases[i].Fork != exp.Fork {
			t.Fatalf("expected alias %d to be name=%q alias=%q fork=%v, got name=%q alias=%q fork=%v", i, exp.Name, exp.Alias, exp.Fork, aliases[i].Name, aliases[i].Alias, aliases[i].Fork)
		}
	}
}

func TestSanitizeOAuthModelAlias_InjectsDefaultGitHubCopilotAliases(t *testing.T) {
	cfg := &Config{
		OAuthModelAlias: map[string][]OAuthModelAlias{
			"codex": {
				{Name: "gpt-5", Alias: "g5"},
			},
		},
	}

	cfg.SanitizeOAuthModelAlias()

	copilotAliases := cfg.OAuthModelAlias["github-copilot"]
	if len(copilotAliases) == 0 {
		t.Fatal("expected default github-copilot aliases to be injected")
	}

	aliasSet := make(map[string]bool, len(copilotAliases))
	for _, a := range copilotAliases {
		aliasSet[a.Alias] = true
		if !a.Fork {
			t.Fatalf("expected all default github-copilot aliases to have fork=true, got fork=false for %q", a.Alias)
		}
	}
	expectedAliases := []string{
		"claude-haiku-4-5",
		"claude-opus-4-1",
		"claude-opus-4-5",
		"claude-opus-4-6",
		"claude-sonnet-4-5",
		"claude-sonnet-4-6",
	}
	for _, expected := range expectedAliases {
		if !aliasSet[expected] {
			t.Fatalf("expected default github-copilot alias %q to be present", expected)
		}
	}
}

func TestSanitizeOAuthModelAlias_DoesNotOverrideUserGitHubCopilotAliases(t *testing.T) {
	cfg := &Config{
		OAuthModelAlias: map[string][]OAuthModelAlias{
			"github-copilot": {
				{Name: "claude-opus-4.6", Alias: "my-opus", Fork: true},
			},
		},
	}

	cfg.SanitizeOAuthModelAlias()

	copilotAliases := cfg.OAuthModelAlias["github-copilot"]
	if len(copilotAliases) != 1 {
		t.Fatalf("expected 1 user-configured github-copilot alias, got %d", len(copilotAliases))
	}
	if copilotAliases[0].Alias != "my-opus" {
		t.Fatalf("expected user alias to be preserved, got %q", copilotAliases[0].Alias)
	}
}

func TestSanitizeOAuthModelAlias_GitHubCopilotDoesNotReinjectAfterExplicitDeletion(t *testing.T) {
	cfg := &Config{
		OAuthModelAlias: map[string][]OAuthModelAlias{
			"github-copilot": nil, // explicitly deleted
		},
	}

	cfg.SanitizeOAuthModelAlias()

	copilotAliases := cfg.OAuthModelAlias["github-copilot"]
	if len(copilotAliases) != 0 {
		t.Fatalf("expected github-copilot aliases to remain empty after explicit deletion, got %d aliases", len(copilotAliases))
	}
	if _, exists := cfg.OAuthModelAlias["github-copilot"]; !exists {
		t.Fatal("expected github-copilot key to be preserved as nil marker after sanitization")
	}
}
