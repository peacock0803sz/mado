package preset

import (
	"context"
	"testing"

	"github.com/peacock0803sz/mado/internal/ax"
)

func TestRecord_TwoWindows(t *testing.T) {
	svc := &ax.MockWindowService{
		Windows: []ax.Window{
			{AppName: "Code", Title: "main.go", PID: 1, X: 0, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
			{AppName: "Terminal", Title: "zsh", PID: 2, X: 960, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
		},
	}

	p, err := Record(context.Background(), svc, "coding")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if p.Name != "coding" {
		t.Errorf("name = %q, want %q", p.Name, "coding")
	}
	if len(p.Rules) != 2 {
		t.Fatalf("len(rules) = %d, want 2", len(p.Rules))
	}

	// Different apps → no title
	for i, r := range p.Rules {
		if r.Title != "" {
			t.Errorf("rules[%d].Title = %q, want empty (different apps)", i, r.Title)
		}
	}

	r0 := p.Rules[0]
	if r0.App != "Code" {
		t.Errorf("rules[0].App = %q, want %q", r0.App, "Code")
	}
	if r0.Position[0] != 0 || r0.Position[1] != 0 {
		t.Errorf("rules[0].Position = %v, want [0, 0]", r0.Position)
	}
	if r0.Size[0] != 960 || r0.Size[1] != 1080 {
		t.Errorf("rules[0].Size = %v, want [960, 1080]", r0.Size)
	}
}

func TestRecord_SameAppMultipleWindows(t *testing.T) {
	svc := &ax.MockWindowService{
		Windows: []ax.Window{
			{AppName: "Code", Title: "main.go", PID: 1, X: 0, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
			{AppName: "Code", Title: "test.go", PID: 1, X: 960, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
		},
	}

	p, err := Record(context.Background(), svc, "dual-editor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Rules) != 2 {
		t.Fatalf("len(rules) = %d, want 2", len(p.Rules))
	}

	// Same app → title should be populated
	if p.Rules[0].Title != "main.go" {
		t.Errorf("rules[0].Title = %q, want %q", p.Rules[0].Title, "main.go")
	}
	if p.Rules[1].Title != "test.go" {
		t.Errorf("rules[1].Title = %q, want %q", p.Rules[1].Title, "test.go")
	}
}

func TestRecord_FiltersNonNormal(t *testing.T) {
	svc := &ax.MockWindowService{
		Windows: []ax.Window{
			{AppName: "Code", Title: "main.go", PID: 1, X: 0, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
			{AppName: "Finder", Title: "Downloads", PID: 3, State: ax.StateMinimized},
			{AppName: "Safari", Title: "Google", PID: 4, State: ax.StateFullscreen},
			{AppName: "Mail", Title: "Inbox", PID: 5, State: ax.StateHidden},
		},
	}

	p, err := Record(context.Background(), svc, "filtered")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Rules) != 1 {
		t.Fatalf("len(rules) = %d, want 1 (only normal windows)", len(p.Rules))
	}
	if p.Rules[0].App != "Code" {
		t.Errorf("rules[0].App = %q, want %q", p.Rules[0].App, "Code")
	}
}

func TestRecord_NoWindows(t *testing.T) {
	svc := &ax.MockWindowService{}

	p, err := Record(context.Background(), svc, "empty")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(p.Rules) != 0 {
		t.Errorf("len(rules) = %d, want 0", len(p.Rules))
	}
}

func TestRecord_InvalidName(t *testing.T) {
	svc := &ax.MockWindowService{}

	_, err := Record(context.Background(), svc, "invalid name!")
	if err == nil {
		t.Fatal("expected error for invalid preset name, got nil")
	}
}
