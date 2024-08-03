package main

import "time"

type SecurityReport struct {
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`
	Issues    []string  `json:"issues" yaml:"issues"`
}
