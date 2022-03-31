package main

import (
	"fmt"
	"time"
)

type message struct {
	Name    string
	Message string
	When    time.Time
}

func (m message) String() string {
	return fmt.Sprintf("Name: %v; When: %v; Message: %v\n", m.Name, m.When, m.Message)
}
