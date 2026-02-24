package cli

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/output"
)

// RootFlags はルートコマンドのグローバルフラグを保持する。
type RootFlags struct {
	Format  string
	Timeout time.Duration
}

// NewRootCmd はルートコマンドを生成する。
// グローバル変数を使わないコンストラクタパターンでテスト可能にする。
func NewRootCmd(svc ax.WindowService) *cobra.Command {
	flags := &RootFlags{}

	root := &cobra.Command{
		Use:   "mado",
		Short: "macOS window management CLI",
		Long: `mado (窓) — macOS のウィンドウを操作するCLIツール。

Accessibility権限が必要なコマンド: list, move
権限不要なコマンド: help, version, completion`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// グローバルフラグ
	root.PersistentFlags().StringVar(&flags.Format, "format", "text", "出力フォーマット (text|json)")
	root.PersistentFlags().DurationVar(&flags.Timeout, "timeout", 5*time.Second, "AX操作タイムアウト")

	// サブコマンド登録 (T024, T031で追加)
	root.AddCommand(newListCmd(svc, flags))
	root.AddCommand(newMoveCmd(svc, flags))
	root.AddCommand(newVersionCmd())
	root.AddCommand(newCompletionCmd(root))

	return root
}

// newOutputFormat はフラグ文字列からoutput.Formatに変換する。
func newOutputFormat(s string) output.Format {
	if s == "json" {
		return output.FormatJSON
	}
	return output.FormatText
}
