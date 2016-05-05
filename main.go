package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ffloyd/evergrid-go/simenv"
)

func main() {
	log.SetLevel(log.DebugLevel)

	dummyAgent1 := simenv.DummyAgent{"Agent 1"}
	dummyAgent2 := simenv.DummyAgent{"Agent 2"}

	gt := new(simenv.GlobalTimer)
	gt.Init()

	gt.AddAgent(dummyAgent1)
	gt.AddAgent(dummyAgent2)

	gt.Run()
}
