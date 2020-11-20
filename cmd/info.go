package cmd

import (
	"fmt"
	"github.com/Scarsz/bincli/bin"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"regexp"
)

func Info(cmd *cobra.Command, args []string) {
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
			fmt.Println("WARN: Bin " + id.String() + " doesn't have a key, only minimal information is known")
		}

		b, err := bin.Retrieve(id, result["key"])
		if err != nil {
			print("ERROR: " + err.Error())
			continue
		}

		fmt.Println("Bin", id.String())
		if b.Description != "" { fmt.Println("-> Description:", b.Description) }
		fmt.Println("Hits:", b.Hits)
		fmt.Println("File count:", len(b.Files))
	}
}
