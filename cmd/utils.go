package main

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

// parse
func parseInput(inp string) (string, error) {
	if len(inp) >= 200 {
		return "", errors.New("the exercism name can't be that long")
	}
	inp = strings.TrimSpace(inp)
	inp = strings.ToLower(inp)

	ExercismId := inp
	re := regexp.MustCompile("exercism download --exercise=([a-zA-Z]+-[a-zA-Z]+)+ --track=go")
	isLink := re.MatchString(ExercismId)

	if isLink {
		IdFound := re.FindStringSubmatch(ExercismId)
		if len(IdFound) < 1 || len(strings.TrimSpace(IdFound[1])) <= 0 {
			return "", errors.New("cannot exctract exercism name from this link")
		}
		ExercismId = IdFound[1]
	} else {
		ExercismId = strings.ReplaceAll(ExercismId, " ", "-")
	}

	for _, ch := range ExercismId {
		if !unicode.IsLetter(ch) && ch != '-' {
			return "", errors.New("invalid exercism name")
		}
	}

	return ExercismId, nil
}
