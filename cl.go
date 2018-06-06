package main

import (
	"github.com/urfave/cli"
	"github.com/mswift42/nip/tv"
	"fmt"
)

func initCli() *cli.App {
	db, err := tv.RestoreProgrammeDB("tv/mockdb.json")
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Setup()
	app.Name = "nip"

	app.Commands = []cli.Command{
		{
			Name: "list",
			Aliases: []string{"l"},
			Usage: "List all available categories.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListAvailableCategories())
				return nil
			},
		},
		{
			Name: "category",
			Aliases: []string{"c"},
			Usage: "List all programmes for a category.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListCategory(c.Args().Get(0)))
				return nil
			},
		},
		{
			Name: "search",
			Aliases: []string{"s"},
			Usage: "Search for a programme.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.FindTitle(c.Args().Get(0)))
				return nil
			},

		},
	}
	return app
}
