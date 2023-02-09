package main

import (
	"log"
	"os"

	"github.com/gauravsarma1992/csvrestapi"
)

func main() {
	var (
		svr *csvrestapi.Server
		err error
	)

	if svr, err = csvrestapi.New(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}

	if err = svr.Run(); err != nil {
		log.Println(err)
		os.Exit(-1)
	}
	log.Println(svr)
}
