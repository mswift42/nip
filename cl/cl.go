package cl

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"

	"bufio"

	"strings"

	"os"

	"github.com/mswift42/nip/tv"
	"github.com/urfave/cli"
)

func extractIndex(c *cli.Context) (int, error) {
	if len(c.Args()) < 1 {
		fmt.Println("Please enter valid index number.")
	}
	ind := c.Args().Get(0)
	index, err := strconv.ParseInt(ind, 10, 64)
	if err != nil {
		fmt.Println("Please enter valid index number.")
		return 0, err
	}
	return int(index), nil
}

// TODO - set folder for storing and reading of db.
// TODO - split SaveDb into more functions for saving of db and refreshing .

// InitCli loads the ProgrammeDB into memory
// and sets up the command line commands.
func InitCli() *cli.App {
	db, err := tv.RestoreProgrammeDB("mockdb.json")
	if err != nil {
		panic(err)
	}
	app := cli.NewApp()
	app.Setup()
	app.Name = "nip"
	app.Usage = "search for iplayer tv programmes."

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
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				url, err := db.FindURL(index)
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
		{
			Name:    "url",
			Aliases: []string{"u"},
			Usage:   "print programme's url",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				url, err := db.FindURL(index)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(url)
				}
				return nil
			},
		},
		{
			Name:    "synopsis",
			Aliases: []string{"syn"},
			Usage:   "print programme's synopsis",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				prog, err := db.FindProgramme(index)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println(prog.String() + "\n" + prog.Synopsis)
				}
				return nil
			},
		},
		{
			Name:    "links",
			Aliases: []string{"lnk"},
			Usage:   "show related links for a programme with index n",
			Action: func(c *cli.Context) error {
				index, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
				}
				rl, err := db.FindRelatedLinks(index)
				if err != nil {
					fmt.Println(err)
				}
				if len(rl) == 0 {
					fmt.Println("Sorry, no related links were found.")
				} else {
					for _, i := range rl {
						fmt.Println(i.Title, " : ", i.Url)
					}
				}
				return nil
			},
		},
		{
			Name:    "download",
			Aliases: []string{"g", "d", "get"},
			Usage:   "use youtube-dl to download programme with index n",
			Action: func(c *cli.Context) error {
				var format string
				if len(c.Args()) == 2 {
					format = c.Args().Get(1)
				}
				ind, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
					return nil
				}
				prog, err := db.FindProgramme(ind)
				if err != nil {
					fmt.Println("Could not find Programme with index ", ind)
					return nil
				}
				fmt.Println("Downloading Programme \n", prog.String())
				u := tv.BBCPrefix + prog.URL
				var cmd *exec.Cmd
				if format != "" {
					cmd = exec.Command("/bin/sh", "-c", "youtube-dl -f "+format+" "+u)
				} else {
					cmd = exec.Command("/bin/sh", "-c", "youtube-dl -f best "+u)
				}
				outpipe, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				scanner := bufio.NewScanner(outpipe)
				scanner.Split(bufio.ScanRunes)
				var target string
				for scanner.Scan() {
					fmt.Print(scanner.Text())
					target += scanner.Text()
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println(err)
				}
				split := strings.Split(target, "\n")
				for _, i := range split {
					if strings.Contains(i, "Destination:") {
						fmt.Println("Found it: ", i[24:])
						cwd, err := os.Getwd()
						if err != nil {
							fmt.Println(err)
						}
						db.MarkSaved(cwd + string(os.PathSeparator) + i[24:])
						fmt.Println(db.SavedProgrammes)
					}
				}
				return nil
			},
		},
		{
			Name:    "formats",
			Aliases: []string{"f"},
			Usage:   "list youtube-dl formats for programme with index n",
			Action: func(c *cli.Context) error {
				ind, err := extractIndex(c)
				if err != nil {
					fmt.Println("Please enter valid index number.")
					return nil
				}
				prog, err := db.FindProgramme(ind)
				if err != nil {
					fmt.Println("could not find Programme with index ", ind)
				}
				fmt.Println("Listing Formats for Programme \n", prog.String())
				u := tv.BBCPrefix + prog.URL
				cmd := exec.Command("/bin/sh", "-c", "youtube-dl -F "+u)
				if err != nil {
					fmt.Println(err)
				}
				outpipe, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
				}
				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
				}
				scanner := bufio.NewScanner(outpipe)
				for scanner.Scan() {
					fmt.Println(scanner.Text())
				}
				err = cmd.Wait()
				if err != nil {
					fmt.Println(err)
				}
				return nil
			},
		},
	}
	return app

}
