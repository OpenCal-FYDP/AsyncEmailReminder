package main

import (
	"fmt"
	"github.com/OpenCal-FYDP/AsyncCalendarOptimizer/internal/emailer"
	"github.com/OpenCal-FYDP/AsyncCalendarOptimizer/internal/storer"
	"log"
	"time"
)

//featureflag
// set to true if we want this not in test mode
const serveEveryone = true

// contains checks if a string is present in a slice
func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// This is a really good starting place https://medium.com/@harlow/processing-kinesis-streams-w-aws-lambda-and-golang-264efc8f979a
func main() {

	emailr, err := emailer.New()
	if err != nil {
		log.Fatalln(err)
	}

	storR := storer.New()

	for {
		events, err := storR.GetEvents()
		if err != nil {
			fmt.Println(err)
		}

		for _, event := range events {
			if contains(event.Attendees, "jspsun@gmail.com") || serveEveryone {
				err := emailr.SendConfirmationEmail("", event.Attendees, event)
				if err != nil {
					fmt.Println(err)
				}

			}
		}
		time.Sleep(time.Minute)
	}
}
