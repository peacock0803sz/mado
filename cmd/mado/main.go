package main

import (
	"fmt"
	"os"

	"github.com/peacock0803sz/mado/internal/ax"
	"github.com/peacock0803sz/mado/internal/cli"
)

// version はビルド時に -ldflags で埋め込まれる。
var version = "dev"

func main() {
	svc := ax.NewWindowService()
	cmd := cli.NewRootCmd(svc)

	// バージョン情報をルートコマンドに注入
	cmd.Version = version

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
