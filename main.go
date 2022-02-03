package main

import (
	"log"
	"os"
)

func main() {
	configfile := "/etc/boggart/config.yml"
	if len(os.Args) > 1 {
		configfile = os.Args[1]
	}

	boggart, err := NewBoggart(configfile)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(boggart.Serve(":2222"))
}
