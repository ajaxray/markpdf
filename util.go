package main

import (
	"fmt"
	"log"
	"os"
)

func debugInfo(message string) {
	if verbose {
		log.Println(message)
	}
}

func fatalIfError(err error, message string) {
	if err != nil {
		fmt.Printf("ERROR: %s \n", message)
		os.Exit(1)
	}
}
