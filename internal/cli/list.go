package cli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/output"
	"github.com/peacock0803sz/mado/internal/window"
)

// newListCmd はlist サブコマンドを生成する (T023)。
func newListCmd(svc ax.WindowService, root *RootFlags) *cobra.Command {
	var appFilter string
	var screenFilter string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "開いているウィンドウ一覧を表示する",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), root.Timeout)
			defer cancel()

			f := output.New(newOutputFormat(root.Format), os.Stdout, os.Stderr)

			if err := svc.CheckPermission(); err != nil {
				_ = f.PrintError(2, err.Error(), nil)
				os.Exit(2)
			}

			opts := window.ListOptions{
				AppFilter:    appFilter,
				ScreenFilter: screenFilter,
			}

			windows, err := window.List(ctx, svc, opts)
			if err != nil {
				return err
			}

			return f.PrintWindows(windows)
		},
	}

	cmd.Flags().StringVar(&appFilter, "app", "", "アプリ名でフィルタ（大文字小文字無視、完全一致）")
	cmd.Flags().StringVar(&screenFilter, "screen", "", "スクリーンIDまたは名前でフィルタ（完全一致）")

	return cmd
}
