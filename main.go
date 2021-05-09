package main

import (
	"fmt"
    
    // TODO: undefined package import
    "block"
    "node"
    "process"
    "dcml"
    "utils"

	log "github.com/sirupsen/logrus"
)

func init() {
    fmt.Println("Usage: go run main.go")
}

func main() {
LOOP:
	Consensus()

	// run the comparison algorith
	fmt.Println("\n----------Would you like to run comparison algorithm?----------\n")
	fmt.Println("Enter y or Y to run comparison algorithm, n or N to say googbye!")

	fmt.Scanln(&input)
	if input == "y" || input == "Y" {
		ComparisonDPoS()
		ComparisonFCSW()
	} else {
		log.Warnf("You will not run comparsion algorithm!")
		return
	}

	// run the next loop
	fmt.Println("\n---------------------------Next loop?----------------------------\n")
	fmt.Println("Current consensus round have done, would you like to start next round? y to continue, n to stop:")

	// interaction in the end
	fmt.Scanln(&input)
	if input == "y" || input == "Y" {
		goto LOOP
	} else {
		log.Infof("Consensus endup, see you next time!")
		return
	}
}
