package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

func Purge(cmd *cobra.Command, args []string) {
	dir, err := ioutil.ReadDir(os.TempDir())
	if err != nil {
		panic(err)
	}

	var counter int
	for _, d := range dir {
		if !d.IsDir() {
			continue
		}
		if !strings.HasPrefix(d.Name(), "bin-") {
			continue
		}
		err := os.RemoveAll(d.Name())
		if err != nil {
			fmt.Println("ERROR: Failed to delete " + d.Name() + ": " + err.Error())
		} else {
			counter++
		}
	}

	fmt.Println("Deleted", counter, "bins")
}
