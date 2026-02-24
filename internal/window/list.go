package window

import (
	"context"
	"strings"

	"github.com/peacock0803sz/mado/internal/ax"
)

// ListOptions はlist コマンドのフィルタオプション。
type ListOptions struct {
	AppFilter    string
	ScreenFilter string
}

// List はウィンドウ一覧を取得してフィルタリングして返す。
func List(ctx context.Context, svc ax.WindowService, opts ListOptions) ([]ax.Window, error) {
	windows, err := svc.ListWindows(ctx)
	if err != nil {
		return nil, err
	}

	return filterWindows(windows, opts), nil
}

// filterWindows はフィルタオプションに基づいてウィンドウ一覧を絞り込む。
func filterWindows(windows []ax.Window, opts ListOptions) []ax.Window {
	result := make([]ax.Window, 0, len(windows))
	for _, w := range windows {
		if opts.AppFilter != "" && !strings.EqualFold(w.AppName, opts.AppFilter) {
			continue
		}
		if opts.ScreenFilter != "" && !matchScreen(w, opts.ScreenFilter) {
			continue
		}
		result = append(result, w)
	}
	return result
}

// matchScreen はスクリーンIDまたは名前でウィンドウをフィルタする。
func matchScreen(w ax.Window, filter string) bool {
	// ID での完全一致
	if w.ScreenName == filter {
		return true
	}
	// スクリーン名での完全一致
	return strings.EqualFold(w.ScreenName, filter)
}
