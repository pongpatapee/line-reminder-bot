package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Reminder struct {
	ID               string
	Text             string
	ReminderDateTime time.Time
}

// Reminder formats:
//
// New reminders: !remindme new <reminder message>.
//
// List all scheduled reminders: !remindme list.
//
// Edit reminder: !remindme edit {reminder ID} <new reminder message>.
//
// Delete reminder: !remindme delete {reminder ID}.
func parseReminderCommand(commandString string) {
	commands := strings.Fields(commandString)

	if len(commands) == 0 || commands[0] != "!remindme" {
		fmt.Println(commands)
		fmt.Println("returning")
		return
	}

	validCommands := map[string]bool{
		"new":    true,
		"list":   true,
		"edit":   true,
		"delete": true,
	}

	if len(commands) <= 1 {
		log.Println("No command given...")
		return
	}

	commandType := commands[1]

	if !validCommands[strings.ToLower(commandType)] {
		log.Printf("Invalid command: %v\n", commandType)
		return
	}

	switch commandType {
	case "new":
		reminderMessage := strings.Join(commands[2:], " ")
		handleNewReminder(reminderMessage)
	case "list":
		handleListReminders()
	case "edit":
		if len(commands) <= 3 {
			log.Fatal("No reminderID or no reminder message recieved")
		}

		reminderId, err := strconv.Atoi(commands[2])
		if err != nil {
			log.Fatalf("Invalid reminder ID: %v\n", commands[2])
			log.Fatal(err)
		}

		newReminderMessage := strings.Join(commands[3:], " ")
		handleEditReminders(reminderId, newReminderMessage)

	case "delete":
		if len(commands) <= 2 {
			log.Fatal("No reminderID received")
		}

		reminderId, err := strconv.Atoi(commands[2])
		if err != nil {
			log.Fatal(err)
		}

		handleDeleteReminders(reminderId)
	}
}

func handleNewReminder(reminderMessage string) {
	fmt.Printf("Adding new reminder: %v\n", reminderMessage)
}

func handleListReminders() {
	fmt.Println("Listing all reminders: ")
}

func handleEditReminders(reminderId int, newReminderMessage string) {
	fmt.Printf("Editing reminderID %v with new reminder message: %v\n", reminderId, newReminderMessage)
}

func handleDeleteReminders(reminderId int) {
	fmt.Printf("Deleting reminderID: %v\n", reminderId)
}

// main func for testing only remove after
func main() {
	parseReminderCommand("!remindme delete 12")
}
