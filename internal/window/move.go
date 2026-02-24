package window

import (
	"context"
	"strings"

	"github.com/peacock0803sz/mado/internal/ax"
)

// Point は座標を表す。
type Point struct {
	X, Y int
}

// Size はサイズを表す。
type Size struct {
	W, H int
}

// MoveOptions はmove コマンドのオプション。
type MoveOptions struct {
	AppFilter    string
	TitleFilter  string
	ScreenFilter string
	Position     *Point
	Size         *Size
	All          bool
}

// Move は対象ウィンドウを移動またはリサイズする。
// 複数一致かつ --all なしの場合は AmbiguousTargetError を返す。
func Move(ctx context.Context, svc ax.WindowService, opts MoveOptions) ([]ax.Window, error) {
	windows, err := svc.ListWindows(ctx)
	if err != nil {
		return nil, err
	}

	targets := filterForMove(windows, opts)

	if len(targets) == 0 {
		return nil, &ax.AmbiguousTargetError{
			Query:      buildQuery(opts),
			Candidates: nil,
		}
	}

	if len(targets) > 1 && !opts.All {
		return nil, &ax.AmbiguousTargetError{
			Query:      buildQuery(opts),
			Candidates: targets,
		}
	}

	var affected []ax.Window
	for _, w := range targets {
		// フルスクリーン状態のウィンドウは操作不可 (exit 5)
		if w.State == ax.StateFullscreen {
			return nil, &fullscreenError{window: w}
		}

		if opts.Position != nil {
			if err := svc.MoveWindow(ctx, w.PID, w.Title, opts.Position.X, opts.Position.Y); err != nil {
				return affected, err
			}
			w.X = opts.Position.X
			w.Y = opts.Position.Y
		}

		if opts.Size != nil {
			if err := svc.ResizeWindow(ctx, w.PID, w.Title, opts.Size.W, opts.Size.H); err != nil {
				return affected, err
			}
			w.Width = opts.Size.W
			w.Height = opts.Size.H
		}

		affected = append(affected, w)
	}

	return affected, nil
}

// filterForMove はmoveコマンド用のウィンドウフィルタ。
func filterForMove(windows []ax.Window, opts MoveOptions) []ax.Window {
	result := make([]ax.Window, 0)
	for _, w := range windows {
		if opts.AppFilter != "" && !strings.EqualFold(w.AppName, opts.AppFilter) {
			continue
		}
		if opts.TitleFilter != "" && !strings.Contains(strings.ToLower(w.Title), strings.ToLower(opts.TitleFilter)) {
			continue
		}
		if opts.ScreenFilter != "" && !matchScreen(w, opts.ScreenFilter) {
			continue
		}
		result = append(result, w)
	}
	return result
}

func buildQuery(opts MoveOptions) string {
	parts := make([]string, 0)
	if opts.AppFilter != "" {
		parts = append(parts, `--app "`+opts.AppFilter+`"`)
	}
	if opts.TitleFilter != "" {
		parts = append(parts, `--title "`+opts.TitleFilter+`"`)
	}
	if len(parts) == 0 {
		return "(no filter)"
	}
	return strings.Join(parts, " ")
}

// fullscreenError はフルスクリーン状態のウィンドウへの操作エラー。
type fullscreenError struct {
	window ax.Window
}

func (e *fullscreenError) Error() string {
	return `cannot move fullscreen window: "` + e.window.Title + `"`
}
