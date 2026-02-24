package window

import (
	"context"
	"strings"

	"github.com/peacock0803sz/mado/internal/ax"
)

// ListOptions holds filter options for the list command.
type ListOptions struct {
	AppFilter    string
	ScreenFilter string
}

// List retrieves all windows and returns them after applying filters.
func List(ctx context.Context, svc ax.WindowService, opts ListOptions) ([]ax.Window, error) {
	windows, err := svc.ListWindows(ctx)
	if err != nil {
		return nil, err
	}

	return filterWindows(windows, opts), nil
}

// filterWindows narrows down the window list based on filter options.
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

// matchScreen filters a window by screen ID or name.
func matchScreen(w ax.Window, filter string) bool {
	// exact match on ID
	if w.ScreenName == filter {
		return true
	}
	// case-insensitive exact match on screen name
	return strings.EqualFold(w.ScreenName, filter)
}
