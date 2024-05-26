package main

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// Reminder formats:
//
// New reminders: !remindme [time string] <reminder message>.
func ParseReminderCommand(commandString string) (time.Duration, string, error) {
	if !strings.HasPrefix(commandString, "!remindme") {
		return 0, "", errors.New("not a reminder")
	}

	commands := strings.Fields(commandString)

	if len(commands) <= 2 {
		return 0, "", errors.New("usage: !remindme [time string] <reminder message>")
	}

	reminderMessage := strings.Join(commands[2:], " ")

	reminderTime, err := time.Parse("3:04pm", commands[1])
	if err != nil {
		return 0, "", err
	}

	currentTime := time.Now()

	reminderTime = time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		reminderTime.Hour(),
		reminderTime.Minute(),
		reminderTime.Second(),
		0,
		currentTime.Location(),
	)

	if reminderTime.Before(currentTime) {
		reminderTime = reminderTime.Add(24 * time.Hour)
	}

	timeTilReminder := reminderTime.Sub(currentTime)

	fmt.Println("Reminder time: ", reminderTime)
	fmt.Println("Time until reminder: ", timeTilReminder)

	return timeTilReminder, reminderMessage, nil
}

// // main func for testing only remove after
// func main() {
// 	duration, reminderMessage, err := parseReminderCommand("!remindme 9:00pm please go do stuff")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Printf("Duration: %v, reminder message: %v\n", duration, reminderMessage)
// }
