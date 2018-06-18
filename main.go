package main

import (
	"github.com/mswift42/nip/cl"
	"os"
)

func main() {
	// TODO - Add ProgressBars for save method.
	//url := "https://www.bbc.co.uk/iplayer/categories/films/all?sort=atoz&page=1"
	//doc := tv.RemoteDocumentLoader(url)
	//nmc := tv.NewMainCategory(doc)
	//fmt.Println(nmc)
	//progs := nmc.Programmes()
	//for _, i := range progs {
	//	fmt.Println(i)
	//}
	//foodurl := "https://www.bbc.co.uk/iplayer/categories/food/all?sort=atoz&page=1"
	//doc = tv.RemoteDocumentLoader(foodurl)
	//nmc = tv.NewMainCategory(doc)
	//progs = nmc.Programmes()
	//for _, i := range progs {
	//	fmt.Println(i)
	//}
	//tv.SaveDB()
	app := cl.InitCli()
	app.Run(os.Args)
}
