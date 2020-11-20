package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func Get(cmd *cobra.Command, args []string) {
	fmt.Println("Downloading bin " + args[0])

	//b := bin.Retrieve(uuid.MustParse("4c72731b-3e41-4c4a-b758-0898311fda7c"))
	//b.Files[0].ContentText()

}
