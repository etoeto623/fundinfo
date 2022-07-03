package main

import (
	"flag"
	"fmt"
	"neolong.me/fundinfo/unit_worth"
)

func main() {
	fundCode := flag.String("c", "", "fund code")
	flag.Parse()

	if len(*fundCode) <= 0 {
		fmt.Println("illegal param")
		return
	}

	unit_worth.CrawFundAllWorth(*fundCode)
}