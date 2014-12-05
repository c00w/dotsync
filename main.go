package main

import (
	"flag"
	"fmt"
	"github.com/c00w/dotsync/internal"
	"log"
	"os"
)

func main() {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		log.Println("Error: $HOME not set")
	}
	curdir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory: ", err)
	}

	flag.Parse()
	switch flag.Arg(0) {
	case "install":
		internal.Install(curdir, homedir)
	case "save":
		internal.Save(homedir, curdir)
	default:
		fmt.Println("dotsync by Colin L. Rice colin@daedrum.net")
		fmt.Println("\tinstall - install existing files onto computer from the current directory")
		fmt.Println("\tsave - save existing rc files into the current directory")
	}
}
