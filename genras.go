package main

import (
	"log"
	"modi/src/utils"
)

func main() {
	err := utils.GenRSAPubAndPri(1024, "./resources/pem")
	if err != nil {
		log.Fatal(err)
	}
}
