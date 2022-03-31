package main

import (
	"fmt"
	"time"
)

type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
	WhenStr   string
}

func (m message) String() string {
	return fmt.Sprintf("Name: %v; When: %v; Message: %v\n", m.Name, m.WhenStr, m.Message)
}

func (m message) WhenString() string {
	y, mo, d := m.When.Date()
	h, M, s := m.When.Clock()
	return fmt.Sprintf("%v-%02d-%02d %02d:%02d:%02d", y, int(mo), d, h, M, s)
}
