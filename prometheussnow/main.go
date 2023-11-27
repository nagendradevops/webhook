package main

import . "github.com/nagendradevops/webhook/prometheussnow/server"

func min() {
	snow, err := snowConfig()
	if err != nil {
		panic(err)
	}
	snow.snowService()
}

/*package main

import . "github.com/gatblau/onix/prometheus/ses/server"

func main() {
	ses, err := NewSeS()
	if err != nil {
		panic(err)
	}
	ses.Start()
}*/
