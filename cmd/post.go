package cmd

import (
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func Post(cmd *cobra.Command, args []string) {
	var files []bin.File
	for _, fileName := range args {
		files = append(files, bin.FileFromFileName(fileName))
	}

	b := bin.Create(bin.Options{
		Expiration: 1440,
		Files:      files,
	})
	fmt.Println(b.URL())
	_ = clipboard.WriteAll(b.URL())
}
