package main

import (
	"fmt"
	poker "learn-go-with-tests/project"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, closeFunc, err := poker.FileSystemPlayerStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer closeFunc()

	fmt.Println("Let's play poker")
	fmt.Println("Type `{Name} wins` to record a win")

	poker.NewCLI(store, os.Stdin).PlayPoker()

}
