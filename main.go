package main

import (
	"github.com/jramirezm00/preview-task-acquia/util"
)

func main() {
	err := util.CreateResponse()
	if err != nil {
		panic(err)
	}
}
