package cl

import (
	"fmt"
	"github.com/mswift42/nip/tv"
	"github.com/urfave/cli"
	"os/exec"
	"runtime"
	"strconv"
)

func InitCli() *cli.App {
	db, err := tv.RestoreProgrammeDB("mockdb.json")
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Setup()
	app.Name = "nip"

	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "List all available categories.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListAvailableCategories())
				return nil
			},
		},
		{
			Name:    "category",
			Aliases: []string{"c"},
			Usage:   "List all programmes for a category.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.ListCategory(c.Args().Get(0)))
				return nil
			},
		},
		{
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "Search for a programme.",
			Action: func(c *cli.Context) error {
				fmt.Println(db.FindTitle(c.Args().Get(0)))
				return nil
			},
		},
		{
			Name:    "show",
			Aliases: []string{"sh"},
			Usage:   "Open Programmes homepage.",
			Action: func(c *cli.Context) error {
				ind := c.Args().Get(0)
				index, err := strconv.ParseInt(ind, 10, 0)
				if err != nil {
					fmt.Println("Please enter valid index number.")
				}
				url, err := db.FindURL(int(index))
				if err != nil {
					fmt.Println(err)
				}
				switch runtime.GOOS {
				case "linux":
					err = exec.Command("xdg-open", url).Start()
				case "darwin":
					err = exec.Command("open", url).Start()
				case "windows":
					err = exec.Command("cmd", "/c", url).Start()
				default:
					fmt.Println("Unsupported platform.")
				}
				if err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
	}
	return app
}
