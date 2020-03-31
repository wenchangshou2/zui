package main

import (
	"fmt"
	"github.com/wenchangshou2/zui/pkg/computer"
)

var (
	G_Backend  computer.Backend
)

func InitComputerControl(){
	var (
		backendName string
		backend computer.Backend
	)
	for _,backendInfo:=range computer.Backends{
		backendName=backendInfo.Name
		var err error
		backend,err=backendInfo.Init()
		if err==nil{
			break
		}
	}
	fmt.Printf("backendName",backendName)
	G_Backend=backend

}
