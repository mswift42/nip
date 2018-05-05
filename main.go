package main

import (
	"github.com/mswift42/nip/tv"
	"fmt"
)

func main() {
	url := "https://www.bbc.co.uk/iplayer/categories/films/all?sort=atoz&page=1"
	doc := tv.RemoteDocumentLoader(url)
	nmc := tv.NewMainCategory(doc)
	progs := nmc.Programmes()
	for _, i := range progs {
		fmt.Println(i.Title)
	}
}
