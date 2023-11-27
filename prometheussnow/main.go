package main

import . "github.com/nagendradevops/webhook/prometheussnow/server"

func main() {
	snow, err := snowConfig()
	if err != nil {
		panic(err)
	}
	snow.snowService()
}
