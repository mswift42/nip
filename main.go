package main

import (
	"github.com/mswift42/nip/tv"
	"fmt"
	tv2 "github.com/mswift42/ipn/tv"
)

func main() {
	url := "https://www.bbc.co.uk/iplayer/categories/films/all?sort=atoz&page=1"
	bu := tv.BeebURL(url)
	fmt.Println(bu)
	c := make(chan *tv2.IplayerDocumentResult)
	go bu
}
