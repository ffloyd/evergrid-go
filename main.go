package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
)

func main() {
	log.SetLevel(log.DebugLevel)

	simulation := simenv.New()
	simulation.Run()
}
