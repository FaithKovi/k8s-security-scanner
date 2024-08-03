package main

import "time"

type SecurityReport struct {
	Timestamp time.Time `yaml:"timestamp"`
	Issues    []string  `yaml:"issues"`
}
