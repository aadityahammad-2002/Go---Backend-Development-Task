package service

import (
	"testing"
	"time"
)

func TestCalculateAgeWhenBirthdayIsToday(t *testing.T) {
	today := time.Now()
	dob := time.Date(today.Year()-30, today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)

	age := CalculateAge(dob)

	if age != 30 {
		t.Errorf("expected age 30, got %d", age)
	}
}

func TestCalculateAgeWhenBirthdayAlreadyPassed(t *testing.T) {
	today := time.Now()
	dob := time.Date(today.Year()-25, time.January, 1, 0, 0, 0, 0, time.UTC)

	age := CalculateAge(dob)

	if age != 25 {
		t.Errorf("expected age 25, got %d", age)
	}
}

func TestCalculateAgeWhenBirthdayNotYetThisYear(t *testing.T) {
	today := time.Now()
	nextMonth := today.AddDate(0, 1, 0)
	dob := time.Date(today.Year()-22, nextMonth.Month(), nextMonth.Day(), 0, 0, 0, 0, time.UTC)

	age := CalculateAge(dob)

	if age != 21 {
		t.Errorf("expected age 21, got %d", age)
	}
}

func TestCalculateAgeWithLeapYearBirthday(t *testing.T) {
	today := time.Now()
	dobYear := today.Year() - 28

	dob := time.Date(dobYear, time.February, 29, 0, 0, 0, 0, time.UTC)
	if !isLeapYear(dobYear) {
		dob = time.Date(dobYear, time.March, 1, 0, 0, 0, 0, time.UTC)
	}

	age := CalculateAge(dob)

	if age < 27 || age > 28 {
		t.Errorf("expected age around 27-28, got %d", age)
	}
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
