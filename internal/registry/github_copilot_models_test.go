package registry

import "testing"

func TestGitHubCopilotGPT55MatchesGPT54Capabilities(t *testing.T) {
	models := GetGitHubCopilotModels()

	var gpt54, gpt55 *ModelInfo
	for _, model := range models {
		if model == nil {
			continue
		}
		switch model.ID {
		case "gpt-5.4":
			gpt54 = model
		case "gpt-5.5":
			gpt55 = model
		}
	}

	if gpt54 == nil {
		t.Fatal("expected gpt-5.4 Copilot model to be registered")
	}
	if gpt55 == nil {
		t.Fatal("expected gpt-5.5 Copilot model to be registered")
	}

	if gpt55.OwnedBy != "github-copilot" || gpt55.Type != "github-copilot" {
		t.Fatalf("expected gpt-5.5 to be a GitHub Copilot model, got owned_by=%q type=%q", gpt55.OwnedBy, gpt55.Type)
	}
	if gpt55.ContextLength != gpt54.ContextLength {
		t.Fatalf("expected gpt-5.5 context length to match gpt-5.4: got %d want %d", gpt55.ContextLength, gpt54.ContextLength)
	}
	if gpt55.MaxCompletionTokens != gpt54.MaxCompletionTokens {
		t.Fatalf("expected gpt-5.5 max completion tokens to match gpt-5.4: got %d want %d", gpt55.MaxCompletionTokens, gpt54.MaxCompletionTokens)
	}
	if !sameStrings(gpt55.SupportedEndpoints, []string{"/responses"}) {
		t.Fatalf("expected gpt-5.5 to use Responses API only so chat requests can be translated, got %+v", gpt55.SupportedEndpoints)
	}
	if gpt55.Thinking == nil || gpt54.Thinking == nil || !sameStrings(gpt55.Thinking.Levels, gpt54.Thinking.Levels) {
		t.Fatalf("expected gpt-5.5 thinking levels to match gpt-5.4: got %+v want %+v", gpt55.Thinking, gpt54.Thinking)
	}
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
