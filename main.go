// Copyright © 2024- Luka Ivanović
// This code is licensed under the terms of the MIT Licence (see LICENCE for details).
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Highlights struct {
	Entries []struct {
		Text    string `json:"text"`
		Chapter string `json:"chapter"`
		Page    int    `json:"page"`
	} `json:"entries"`
	NumberOfPages int `json:"number_of_pages"`
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: format-highlights <path> [<path>...]")
		os.Exit(1)
	}

	for _, arg := range args {
		file, err := os.ReadFile(arg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		var highlights Highlights
		err = json.Unmarshal(file, &highlights)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}

		f, err := os.OpenFile(arg+".end", os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create `%s` file. Using stdout for output file.\n", arg+".end")
			f = os.Stdout
		}

		for _, entry := range highlights.Entries {
			text := entry.Text
			text = strings.ReplaceAll(text, "‘", "'")
			text = strings.ReplaceAll(text, "’", "'")
			text = strings.ReplaceAll(text, "“", "\"")
			text = strings.ReplaceAll(text, "”", "\"")
			text = strings.ReplaceAll(text, "—", " - ")
			lines := strings.Split(text, "\n")
			// TODO: handle errors
			f.WriteString("> ")
			f.WriteString(strings.Join(lines, "\n> "))
			f.WriteString("\n> \n> ")
			f.WriteString(entry.Chapter)
			f.WriteString(fmt.Sprintf(" (page %d/%d)\n\n", entry.Page, highlights.NumberOfPages))
		}
		if f == os.Stdout {
			f.WriteString("---\n\n")
		}
	}
}
