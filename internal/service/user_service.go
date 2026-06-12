package service

import (
	"time"
)

func CalculateAge(dob time.Time) int {
	today := time.Now()

	age := today.Year() - dob.Year()

	if today.Month() < dob.Month() || (today.Month() == dob.Month() && today.Day() < dob.Day()) {
		age--
	}

	return age
}
