package main

import (
	"github.com/mswift42/nip/cl"
	"os"
	"log"
)

func main() {
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
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
