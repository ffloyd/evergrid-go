package simenv

import "testing"

func TestSimEnvWithFewSimpleAgents(t *testing.T) {
	agent1 := NewSimpleAgent("1")
	agent2 := NewSimpleAgent("2")
	agent3 := NewSimpleAgent("3")

	se := NewSimEnv()
	se.Add(agent1, agent2, agent3)
	se.Run()

	if agent1.MessageCount != 6 {
		t.Fail()
	}

	if agent2.MessageCount != 6 {
		t.Fail()
	}

	if agent3.MessageCount != 6 {
		t.Fail()
	}
}
