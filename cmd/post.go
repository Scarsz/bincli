package cmd

import (
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/spf13/cobra"
	"strings"
)

func Post(cmd *cobra.Command, args []string) {
	fmt.Println("Creating bin from files: " + strings.Join(args, ", "))

	var files []bin.File
	for _, fileName := range args {
		files = append(files, bin.FileFromFileName(fileName))
	}

	b := bin.Create(bin.Options{
		Expiration: 1440,
		Files:      files,
	})
	fmt.Println(b.URL())
}
