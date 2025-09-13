package main

import "api-gin/cmd"

func main() {
	if err := cmd.Excute(); err != nil {
		panic(err)
	}
}
