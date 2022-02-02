package main

import (
	"log"
)

func main() {
	boggart, err := NewBoggart("/etc/boggart/config.yml")
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(boggart.Serve(":2222"))
}
