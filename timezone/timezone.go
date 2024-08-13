package timezone

import (
	"fmt"
	"time"
)

type Timezone struct {
	Identifier string
	Offset     string
}

//
// If timezone is incorrect, it will return an error and a default timezone of UTC
//

var UTC = &Timezone{
	Identifier: "UTC",
	Offset:     "+00:00",
}

func GetTimezoneData(identifier string) (*Timezone, error) {
	// load the timezone location
	location, err := time.LoadLocation(identifier)
	if err != nil {
		return UTC, fmt.Errorf("error loading location: %v", err)
	}

	// get the current time in UTC
	now := time.Now().UTC()

	// convert the current time to the specified timezone
	timeInLocation := now.In(location)

	// get the offset in seconds
	_, offset := timeInLocation.Zone()

	// convert offset to hours and minutes
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60

	// format the offset as a string
	offsetString := fmt.Sprintf("%02d:%02d", offsetHours, offsetMinutes)

	return &Timezone{
		Identifier: identifier,
		Offset:     offsetString,
	}, nil
}
