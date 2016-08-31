package main

import (
	"log"

	"github.com/asdine/brazier/cli"
	"github.com/asdine/brazier/store/boltdb"
)

func main() {
	s := boltdb.NewStore("brazier.db")
	defer s.Close()

	cmd := cli.New(s)
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
