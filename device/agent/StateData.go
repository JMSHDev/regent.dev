package main

import "fmt"

type StateData struct {
	AgentStatus   string
	ProgramStatus string
}

func (s *StateData) ToJson() string {
	return fmt.Sprintf("{\"agentStatus\": \"%v\", \"programStatus\": \"%v\"}", s.AgentStatus, s.ProgramStatus)
}
