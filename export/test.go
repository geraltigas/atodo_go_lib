package main

import "C"
import (
	"atodo_go/table"
	"atodo_go/task_show"
	"encoding/json"
)

//export InitDB
func InitDB() {
	err := table.InitDB()
	if err != nil {
		return
	}
}

//export Export
func Export() *C.char {
	data, err := task_show.GetShowData()
	if err != nil {
		return C.CString("")
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return C.CString("")
	}
	return C.CString(string(marshal))
}

func main() {}
