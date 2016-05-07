package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simulation"
)

func main() {
	log.SetLevel(log.DebugLevel)

	sim := simulation.New()
	sim.Run()
}
