package main

import (
	"atodo_go/table"
	"fmt"
)

func main() {
	fmt.Println("Hello, World!")
	err := table.InitDB()
	if err != nil {
		return
	}
	table.InitTaskTable()
}
