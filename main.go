package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	workDuration      time.Duration
	breakDuration     time.Duration
	longBreakDuration time.Duration
	sessions          int
)

func init() {
	flag.DurationVar(&workDuration, "work", 25*time.Minute, "Specifies the work duration in minutes (e.g. 25m).")
	flag.DurationVar(&breakDuration, "break", 5*time.Minute, "Specifies the break duration in minutes (e.g. 5m).")
	flag.DurationVar(&longBreakDuration, "longbreak", 15*time.Minute, "Specifies the long break duration in minutes (e.g. 10m).")
	flag.IntVar(&sessions, "sessions", 3, "Specifies the number of work sessions before a long break.")
	flag.Parse()
}

func main() {

	fmt.Printf("Starting a %v work session..", workDuration)
	time.Sleep(workDuration)
	fmt.Printf("Go take a %v break", breakDuration)
}
