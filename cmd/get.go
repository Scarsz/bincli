package cmd

import (
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/google/uuid"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"regexp"
)

func Get(cmd *cobra.Command, args []string) {
	expression := regexp.MustCompile("(?P<uuid>[a-f0-9-]{36})(#(?P<key>[A-z0-9]{32}))?")

	for _, raw := range args {
		match := expression.FindStringSubmatch(raw)
		result := make(map[string]string)
		for i, name := range expression.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		if result["uuid"] == "" {
			fmt.Println("ERROR: UUID not found in argument " + raw)
			continue
		}
		id := uuid.MustParse(result["uuid"])
		key := result["key"]

		if key == "" {
			fmt.Println("WARN: Bin " + id.String() + " doesn't have a valid key, can't decrypt")
			continue
		}

		b, err := bin.Retrieve(id, result["key"])
		if err != nil {
			print("ERROR: " + err.Error())
			continue
		}

		dir := b.SaveToTemp()
		if len(b.Files) > 1 {
			err = open.Run(dir)
		} else if len(b.Files) == 1 {
			err = open.Run(dir + "/" + b.Files[0].Name)
		}
		if err != nil {
			print("ERROR: " + err.Error())
			continue
		}
	}
}
