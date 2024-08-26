package main

import (
	"fx-vote-server/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       fx vote server application
// @version                     1.0.0
// @description                 vote AppMain
func main() {
	initialize.InitEnv()
}
