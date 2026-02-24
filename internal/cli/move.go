package cli

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/output"
	"github.com/peacock0803sz/mado/internal/window"
)

// newMoveCmd はmove サブコマンドを生成する (T029)。
func newMoveCmd(svc ax.WindowService, root *RootFlags) *cobra.Command {
	var (
		appFilter    string
		titleFilter  string
		screenFilter string
		positionStr  string
		sizeStr      string
		all          bool
	)

	cmd := &cobra.Command{
		Use:   "move",
		Short: "ウィンドウを移動またはリサイズする",
		RunE: func(cmd *cobra.Command, _ []string) error {
			// T030: --position も --size も未指定の場合 exit 3
			if positionStr == "" && sizeStr == "" {
				f := output.New(newOutputFormat(root.Format), os.Stdout, os.Stderr)
				_ = f.PrintError(3, "--position または --size のいずれかが必要です", nil)
				os.Exit(3)
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), root.Timeout)
			defer cancel()

			f := output.New(newOutputFormat(root.Format), os.Stdout, os.Stderr)

			if err := svc.CheckPermission(); err != nil {
				_ = f.PrintError(2, err.Error(), nil)
				os.Exit(2)
			}

			opts := window.MoveOptions{
				AppFilter:    appFilter,
				TitleFilter:  titleFilter,
				ScreenFilter: screenFilter,
				All:          all,
			}

			if positionStr != "" {
				x, y, err := parseCoords(positionStr)
				if err != nil {
					_ = f.PrintError(3, fmt.Sprintf("--position の値が不正です: %v", err), nil)
					os.Exit(3)
				}
				opts.Position = &window.Point{X: x, Y: y}
			}

			if sizeStr != "" {
				w, h, err := parseCoords(sizeStr)
				if err != nil {
					_ = f.PrintError(3, fmt.Sprintf("--size の値が不正です: %v", err), nil)
					os.Exit(3)
				}
				if w <= 0 || h <= 0 {
					_ = f.PrintError(3, "--size の幅と高さは正の整数が必要です", nil)
					os.Exit(3)
				}
				opts.Size = &window.Size{W: w, H: h}
			}

			affected, err := window.Move(ctx, svc, opts)
			if err != nil {
				var exitCode int
				switch err.(type) {
				case *ax.AmbiguousTargetError:
					exitCode = 4
					ambigErr := err.(*ax.AmbiguousTargetError)
					_ = f.PrintError(4, ambigErr.Error(), ambigErr.Candidates)
					os.Exit(4)
				default:
					return err
				}
				_ = exitCode
			}

			return f.PrintMoveResult(affected)
		},
	}

	cmd.Flags().StringVar(&appFilter, "app", "", "アプリ名でフィルタ（大文字小文字無視、完全一致）")
	cmd.Flags().StringVar(&titleFilter, "title", "", "タイトルでフィルタ（大文字小文字無視、部分一致）")
	cmd.Flags().StringVar(&screenFilter, "screen", "", "スクリーンIDまたは名前でフィルタ")
	cmd.Flags().StringVar(&positionStr, "position", "", "移動先座標 x,y（グローバル座標）")
	cmd.Flags().StringVar(&sizeStr, "size", "", "変更後サイズ width,height")
	cmd.Flags().BoolVar(&all, "all", false, "複数一致時に全ウィンドウへ適用")

	return cmd
}

// parseCoords は "x,y" 形式の文字列を2つの整数にパースする。
func parseCoords(s string) (int, int, error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("カンマ区切りで2値を指定してください (例: 100,200)")
	}
	a, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, fmt.Errorf("最初の値が整数ではありません: %q", parts[0])
	}
	b, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, fmt.Errorf("2番目の値が整数ではありません: %q", parts[1])
	}
	return a, b, nil
}
