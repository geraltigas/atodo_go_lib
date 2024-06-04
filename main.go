package main

import (
	"atodo_go/table"
	"atodo_go/web"
)

func main() {
	err := table.InitDB()
	if err != nil {
		return
	}
	web.RunWebServer(web.InitWebInterface())
}
