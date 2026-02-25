package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/cli"
)

// executePresetCmd は preset サブコマンドを実行してエラーを返す
func executePresetCmd(t *testing.T, svc ax.WindowService, configContent string, args ...string) error {
	t.Helper()

	dir := t.TempDir()
	cfgFile := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgFile, []byte(configContent), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MADO_CONFIG", cfgFile)

	cmd := cli.NewRootCmd(svc)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)

	return cmd.Execute()
}

var validPresetConfig = `presets:
  - name: coding
    description: "Editor left, terminal right"
    rules:
      - app: Code
        position: [0, 0]
        size: [960, 1080]
      - app: Terminal
        position: [960, 0]
        size: [960, 1080]
`

func TestPresetList_Empty(t *testing.T) {
	svc := &ax.MockWindowService{}
	err := executePresetCmd(t, svc, "format: text\n", "preset", "list")
	if err != nil {
		t.Fatalf("expected no error for empty preset list, got: %v", err)
	}
}

func TestPresetList_WithPresets(t *testing.T) {
	svc := &ax.MockWindowService{}
	err := executePresetCmd(t, svc, validPresetConfig, "preset", "list")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestPresetShow_MissingArgs(t *testing.T) {
	svc := &ax.MockWindowService{}
	cmd := cli.NewRootCmd(svc)
	cmd.SetArgs([]string{"preset", "show"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	// Cobra handles ExactArgs(1) → error
	if err == nil {
		t.Fatal("expected error for missing args, got nil")
	}
}

func TestPresetApply_MissingArgs(t *testing.T) {
	svc := &ax.MockWindowService{}
	cmd := cli.NewRootCmd(svc)
	cmd.SetArgs([]string{"preset", "apply"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	// Cobra handles ExactArgs(1) → error
	if err == nil {
		t.Fatal("expected error for missing args, got nil")
	}
}

func TestPresetRec_MissingArgs(t *testing.T) {
	svc := &ax.MockWindowService{}
	cmd := cli.NewRootCmd(svc)
	cmd.SetArgs([]string{"preset", "rec"})
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	// Cobra handles RangeArgs(1, 2) → error
	if err == nil {
		t.Fatal("expected error for missing args, got nil")
	}
}

func TestPresetRec_Stdout(t *testing.T) {
	svc := &ax.MockWindowService{
		Windows: []ax.Window{
			{AppName: "Code", Title: "main.go", PID: 1, X: 0, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
			{AppName: "Terminal", Title: "zsh", PID: 2, X: 960, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
		},
	}

	cmd := cli.NewRootCmd(svc)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"preset", "rec", "coding"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := out.String()
	if !strings.Contains(output, "name: coding") {
		t.Errorf("output missing preset name, got:\n%s", output)
	}
	if !strings.Contains(output, "app: Code") {
		t.Errorf("output missing Code app, got:\n%s", output)
	}
	if !strings.Contains(output, "app: Terminal") {
		t.Errorf("output missing Terminal app, got:\n%s", output)
	}
}

func TestPresetRec_ToFile(t *testing.T) {
	svc := &ax.MockWindowService{
		Windows: []ax.Window{
			{AppName: "Code", Title: "main.go", PID: 1, X: 0, Y: 0, Width: 960, Height: 1080, State: ax.StateNormal},
		},
	}

	dir := t.TempDir()
	outFile := filepath.Join(dir, "preset.yaml")

	cmd := cli.NewRootCmd(svc)
	var out bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs([]string{"preset", "rec", "coding", outFile})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "name: coding") {
		t.Errorf("file missing preset name, got:\n%s", content)
	}
	if !strings.Contains(content, "app: Code") {
		t.Errorf("file missing Code app, got:\n%s", content)
	}
}
