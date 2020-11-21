package cmd

import (
	"bufio"
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func Paste(cmd *cobra.Command, args []string) {
	var content string

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		// no piped input, use args instead
		if len(args) == 0 {
			fmt.Println("ERROR: No text given to paste")
			return
		} else {
			content = strings.Join(args, " ")
		}
	} else {
		// have piped input
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			content += scanner.Text() + "\n"
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	files := []bin.File{bin.FileFromText("cli.txt", content, "")}

	b := bin.Create(bin.Options{
		Expiration: 1440,
		Files:      files,
	})
	fmt.Println(b.URL())
	_ = clipboard.WriteAll(b.URL())
}
