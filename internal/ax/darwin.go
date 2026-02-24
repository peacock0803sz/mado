//go:build darwin

package ax

/*
#cgo LDFLAGS: -framework ApplicationServices -framework CoreGraphics

#include <ApplicationServices/ApplicationServices.h>
#include <CoreGraphics/CoreGraphics.h>

int ax_is_trusted() {
    return AXIsProcessTrusted() ? 1 : 0;
}
*/
import "C"

import (
	"context"
	"fmt"
)

// darwinService はAX APIを使うWindowServiceのdarwin実装。
type darwinService struct{}

// NewWindowService は実際のAX APIを使うWindowServiceを返す。
func NewWindowService() WindowService {
	return &darwinService{}
}

func (s *darwinService) CheckPermission() error {
	// T015で実装
	if C.ax_is_trusted() == 0 {
		return &PermissionError{}
	}
	return nil
}

func (s *darwinService) ListWindows(ctx context.Context) ([]Window, error) {
	// T018で実装
	return nil, fmt.Errorf("not implemented")
}

func (s *darwinService) ListScreens(ctx context.Context) ([]Screen, error) {
	// T058で実装
	return nil, fmt.Errorf("not implemented")
}

func (s *darwinService) MoveWindow(ctx context.Context, pid uint32, title string, x, y int) error {
	// T025で実装
	return fmt.Errorf("not implemented")
}

func (s *darwinService) ResizeWindow(ctx context.Context, pid uint32, title string, w, h int) error {
	// T025で実装
	return fmt.Errorf("not implemented")
}
