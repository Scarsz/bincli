package main

import (
	"fmt"
	"github.com/Scarsz/bincli/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{Use: "bin"}

	cmdInfo := &cobra.Command{
		Use: "info <id/url>",
		Short: "Get information about a bin",
		Args: cobra.MinimumNArgs(1),
		Run: cmd.Info,
	}

	// uploading commands
	cmdGet := &cobra.Command{
		Use: "get <url>",
		Short: "Download a bin",
		Long: "Download a bin from the given URL",
		Args: cobra.MinimumNArgs(1),
		Run: cmd.Get,
	}
	cmdPost := &cobra.Command{
		Use: "post <file>...",
		Short: "Create a new bin from the given files",
		Args: cobra.MinimumNArgs(1),
		Run: cmd.Post,
	}
	cmdPaste := &cobra.Command{
		Use: "paste <text>",
		Short: "Create a new bin with the given text",
		Args: cobra.MinimumNArgs(1),
		Run: cmd.Paste,
	}

	rootCmd.AddCommand(cmdInfo, cmdGet, cmdPost, cmdPaste)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}
