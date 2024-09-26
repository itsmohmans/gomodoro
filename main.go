package main

import (
	"errors"
	"flag"
	"fmt"
	"os/exec"
	"strings"
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
	sound          bool
	autostart      bool
	currentSession SessionType
)

func init() {
	// FIXME: check and invalidate negative value or invalid inputs
	flag.DurationVar(&Work.Duration, "work", 25*time.Minute, "Specifies the work duration in minutes (e.g. 25m).")
	flag.DurationVar(&Break.Duration, "break", 5*time.Minute, "Specifies the break duration in minutes (e.g. 5m).")
	flag.DurationVar(&LongBreak.Duration, "longbreak", 15*time.Minute, "Specifies the long break duration in minutes (e.g. 10m).")
	flag.IntVar(&maxSessions, "sessions", 3, "Specifies the number of work sessions before a long break.")
	flag.BoolVar(&sound, "sound", false, "Specifies whether the timer should play a beep sound or not when the time's up.")
	flag.BoolVar(&autostart, "auto", false, "Specifies whether the timer should automatically start the next session.")
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
func startNextSession(timer *time.Timer, bar *progressbar.ProgressBar) {
	timer.Reset(currentSession.Duration)
	bar.Clear()
	// FIXME: the progress bar should take the new session's duration value as well
	bar.Reset()
	bar.Describe(fmt.Sprintf("%s: [%d/%d]", currentSession.Name, sessionNumber, maxSessions))
}

func playSound(filePath string) {
	// For linux only ig for now, make sure `aplay` is installed
	// TODO: use https://github.com/gopxl/beep to handle this instead
	err := exec.Command("aplay", filePath).Run()
	if err != nil {
		fmt.Println("Error playing sound:", err)
	}
}

func startTimer() {
	timer := time.NewTimer(currentSession.Duration)
	ticker := time.NewTicker(1 * time.Second)
	bar := progressbar.Default(int64(currentSession.Duration.Seconds()), fmt.Sprintf("%s: [%d/%d]", currentSession.Name, sessionNumber, maxSessions))
	for {
		select {
		case <-timer.C:
			if sound {
				playSound("./sounds/chime.wav")
			}
			switchSession()
			if !autostart {
				var startNext rune
				fmt.Printf("\nStart next session: %v %s (Y/n -- exit)? ", currentSession.Duration, currentSession.Name)
				_, err := fmt.Scanf("%c", &startNext)
				if err != nil {
					fmt.Println("Error reading input: ", err)
				}
				if strings.ToLower(string(startNext)) == "y" || startNext == '\n' {
					startNextSession(timer, bar)
				} else {
					timer.Stop()
					ticker.Stop()
					bar.Close()
					return
				}
			} else {
				startNextSession(timer, bar)
			}
		case <-ticker.C:
			bar.Add(1)
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func main() {
	startTimer()
}
