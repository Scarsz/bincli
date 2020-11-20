package cmd

import (
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/spf13/cobra"
	"strings"
)

func Paste(cmd *cobra.Command, args []string) {
	content := strings.Join(args, " ")
	fmt.Println("Creating bin from text: " + content)

	files := []bin.File{bin.FileFromText("cli.txt", content, "")}

	b := bin.Create(bin.Options{
		Expiration: 1440,
		Files:      files,
	})
	fmt.Println(b.URL())
}
