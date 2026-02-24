package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// newVersionCmd はversionサブコマンドを生成する (T050)。
// Accessibility権限不要。
func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "バージョン情報を表示する",
		Run: func(cmd *cobra.Command, _ []string) {
			// cmd.Root().Version はmain.goで埋め込まれる
			fmt.Fprintf(cmd.OutOrStdout(), "mado version %s\n", cmd.Root().Version)
		},
	}
}
