package ax

import "fmt"

// PermissionError はAccessibility権限がない場合のエラー。
// エラーメッセージに macOS の解決手順を含める (Constitution II)。
type PermissionError struct{}

func (e *PermissionError) Error() string {
	return "Accessibility permission not granted"
}

// Resolution は権限不足の解決手順を返す。
func (e *PermissionError) Resolution() string {
	return `To grant permission:
  1. Open System Settings → Privacy & Security → Accessibility
  2. Enable mado (or your Terminal app) in the list
  3. Re-run the command`
}

// AmbiguousTargetError は複数のウィンドウが一致した場合のエラー。
type AmbiguousTargetError struct {
	Query      string
	Candidates []Window
}

func (e *AmbiguousTargetError) Error() string {
	return fmt.Sprintf("ambiguous target: %d windows match %q", len(e.Candidates), e.Query)
}

// TimeoutError はAX操作タイムアウト時のエラー。
type TimeoutError struct {
	Op string
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("AX operation timed out: %s", e.Op)
}
