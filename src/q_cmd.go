package src

import (
	"fmt"
	"log"
	"strings"

	"sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
)

func init() {
	_, err := CLI.AddCommand("q",
		"query",
		"The q command queries for definitions.",
		&queryCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type QueryCmd struct{}

var queryCmd QueryCmd

func (c *QueryCmd) Execute(args []string) error {
	query := strings.Join(args, " ")
	// defs, _, err := apiclient.Defs.List(&sourcegraph.DefListOptions{
	// 	Query:    query,
	// 	Doc:      false,
	// 	Exported: true,
	// })
	res, _, err := apiclient.Search.Search(&sourcegraph.SearchOptions{
		Query:       query,
		Defs:        true,
		ListOptions: sourcegraph.ListOptions{PerPage: 1},
	})
	if err != nil {
		return err
	}
	defs := res.Defs

	if len(defs) == 0 {
		fmt.Printf("No results for %q.\n", query)
		return nil
	}
	for _, def := range defs {
		fmt.Printf("%s: %s\n", def.Repo, def.Name)

		// Fetch docs and stats.
		def, _, err = apiclient.Defs.Get(def.DefSpec(), &sourcegraph.DefGetOptions{Doc: true})
		if err != nil {
			return err
		}
		fmt.Println(stripHTML(def.DocHTML))
	}
	return nil
}

func stripHTML(html string) string {
	return strings.Replace(strings.Replace(html, "<p>", "", -1), "</p>", "", -1)
}
