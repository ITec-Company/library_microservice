package model

import (
	"errors"
	"regexp"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var (
	// regex
	latin       = regexp.MustCompile(`\p{Latin}`)
	cyrillic    = regexp.MustCompile(`[\p{Cyrillic}]`)
	phone       = regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	noSQL       = regexp.MustCompile(`\b(ALTER|CREATE|DELETE|DROP|EXEC(UTE){0,1}|INSERT( +INTO){0,1}|MERGE|SELECT|UPDATE|UNION( +ALL){0,1})\b`)
	onlyLetters = regexp.MustCompile("[^a-zA-Zа-яА-Я]+")

	// errors

	// ErrContainsSQL ...
	ErrContainsSQL = errors.New("no SQL commands allowed to input")
	// ErrInvalidBithDate ...
	ErrInvalidBithDate = errors.New("invalid date of birth. Age must be from 18 to 100. Date format RFC3339")
	// ErrInvalidDevDirection ...
	ErrInvalidDevDirection = errors.New("invalid development direction")
	// ErrInvalidPhoneNumber ...
	ErrInvalidPhoneNumber = errors.New("invalid phone number format")
	// ErrInvalidAlphabet ...
	ErrInvalidAlphabet = errors.New("only latin or cyrillic symblos allowed")
	// ErrInvalidSymbol ...
	ErrInvalidSymbol = errors.New("invalid symbol used. Only space and '-' symbols allowed")
	// ErrInvalidStartDate ...
	ErrInvalidStartDate = errors.New("invalid start date. Start date cannot be before today")
	// ErrInvalidEndDate ...
	ErrInvalidEndDate = errors.New("invalid end date. End date cannot be before today")
	// ErrInvalidID ...
	ErrInvalidID = errors.New("invalid input: id")
)

// IsLetterHyphenSpaces checks if string contains only letter(from simillar alphabet(latin or cyrillic)), hyphen or spaces
// Valid:"Name", "Name name", "Name-name"
// Invalid: "Name123", "NameИмя", "Name@name"
func IsLetterHyphenSpaces(value interface{}) error {
	s := value.(string)
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "-", "", -1)

	err := is.UTFLetter.Validate(s)
	if err != nil {
		return ErrInvalidSymbol
	}
	if cyrillic.MatchString(s) && !latin.MatchString(s) {
		return nil
	} else if latin.MatchString(s) && !cyrillic.MatchString(s) {
		return nil
	}
	return ErrInvalidAlphabet
}

// IsPhone ...
func IsPhone(value interface{}) error {
	s := value.(string)

	if phone.MatchString(s) {
		return nil
	}
	return ErrInvalidPhoneNumber
}

// IsDevDirection checks if string matchs to a Direction
// Frontend Direction = "frontend"
// Backend  Direction = "backend"
// Database Direction = "database"
// Testing  Direction = "testing"
func IsDevDirection(value interface{}) error {
	s := Direction(value.(string))
	if s == Frontend || s == Backend || s == Database || s == Testing {
		return nil
	}
	return ErrInvalidDevDirection
}

// IsDevSubDirection checks if string matchs to a SubDirection
// Java       SubDirection = "frontend"
// Golang     SubDirection = "backend"
// Python     SubDirection = "database"
// JavaScript SubDirection = "testing"
// Postgres   SubDirection = "postgres"
func IsDevSubDirection(value interface{}) error {
	s := SubDirection(value.(string))
	if s == Java || s == Python || s == JavaScript || s == Golang || s == Postgres {
		return nil
	}
	return ErrInvalidDevDirection
}

// IsValidBirthDate ...
func IsValidBirthDate(value interface{}) error {
	t := time.Now()
	d := value.(*time.Time)
	err := validation.Validate(d.Format(time.RFC3339), validation.Date(time.RFC3339).Max(t.AddDate(-18, 0, 0)).Min(t.AddDate(-100, 0, 0)))
	if err != nil {
		return ErrInvalidBithDate
	}
	return nil
}

// IsSQL ...
func IsSQL(value interface{}) error {
	s := value.(string)

	if noSQL.MatchString(strings.ToUpper(s)) {
		return ErrContainsSQL
	}

	str := onlyLetters.ReplaceAllString(s, "")

	if noSQL.MatchString(strings.ToUpper(str)) {
		return ErrContainsSQL
	}
	return nil
}

// IsValidStartDate ...
func IsValidStartDate(value interface{}) error {
	t := time.Now()
	d := value.(*time.Time)
	err := validation.Validate(d.Format(time.RFC3339), validation.Date(time.RFC3339).Min(t.AddDate(0, 0, -1)))
	if err != nil {
		return ErrInvalidStartDate
	}
	return nil
}

// IsValidEndDate ...
func IsValidEndDate(value interface{}) error {
	t := time.Now()
	d := value.(*time.Time)
	err := validation.Validate(d.Format(time.RFC3339), validation.Date(time.RFC3339).Min(t))
	if err != nil {
		return ErrInvalidEndDate
	}
	return nil
}

// IsValidID ...
func IsValidID(value interface{}) error {
	id := value.(int)
	if id < 1 {
		return ErrInvalidID
	}
	return nil
}
