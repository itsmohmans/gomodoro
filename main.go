package main

import (
	"errors"
	"flag"
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

type SessionType struct {
	Name     string
	Duration time.Duration
}

var (
	Work      = SessionType{Name: "work", Duration: 25 * time.Minute}
	Break     = SessionType{Name: "break", Duration: 5 * time.Minute}
	LongBreak = SessionType{Name: "longbreak", Duration: 15 * time.Minute}
)

var (
	maxSessions    int
	sessionNumber  int
	currentSession SessionType
)

func init() {
	flag.DurationVar(&Work.Duration, "work", 25*time.Minute, "Specifies the work duration in minutes (e.g. 25m).")
	flag.DurationVar(&Break.Duration, "break", 5*time.Minute, "Specifies the break duration in minutes (e.g. 5m).")
	flag.DurationVar(&LongBreak.Duration, "longbreak", 15*time.Minute, "Specifies the long break duration in minutes (e.g. 10m).")
	flag.IntVar(&maxSessions, "sessions", 3, "Specifies the number of work sessions before a long break.")
	flag.Parse()
	currentSession = Work
	sessionNumber = 1
}

func setSession(session SessionType) error {
	switch session {
	case Work, Break, LongBreak:
		currentSession = session
		return nil
	default:
		return errors.New("invalid session type")
	}
}

func switchSession() {
	switch currentSession {
	case Break, LongBreak:
		err := setSession(Work)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		return
	default:
		if sessionNumber == maxSessions {
			err := setSession(LongBreak)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			sessionNumber = 1 // reset
		} else {
			err := setSession(Break)
			if err != nil {
				fmt.Println("Error: ", err)
			}
			sessionNumber++
		}
	}
}

func startTimer() {
	// fmt.Printf("Starting timer:\n [1/%d] a %v %s session.\n", maxSessions, currentSession.Duration, currentSession.Name)
	var bar *progressbar.ProgressBar
	secs := 0

	timer := time.NewTimer(currentSession.Duration)
	ticker := time.NewTicker(1 * time.Second)
	bar = progressbar.Default(int64(currentSession.Duration.Seconds()), fmt.Sprintf("%s: [%d/%d]", currentSession.Name, sessionNumber, maxSessions))
	for {
		select {
		case <-timer.C:
			// fmt.Printf("Just finished a %s session, the next session is ", currentSession.Name)
			switchSession()
			// fmt.Printf("a %v %s.\n", currentSession.Duration, currentSession.Name)
			timer.Reset(currentSession.Duration)
			secs = 0
			// bar = progressbar.Default(int64(currentSession.Duration.Seconds()), fmt.Sprintf("%s: [%d/%d]", currentSession.Name, sessionNumber, maxSessions))
			bar.Clear()
			bar.Reset()
			bar.Describe(fmt.Sprintf("%s: [%d/%d]", currentSession.Name, sessionNumber, maxSessions))
			// return
		case <-ticker.C:
			secs++
			bar.Add(1)
			// fmt.Printf("%d/%v\n", secs, currentSession.Duration.Seconds())
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	startTimer()
}
