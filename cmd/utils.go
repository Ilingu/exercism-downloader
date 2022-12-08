package main

import (
	"errors"
	subresult "exercism-cli/cmd/subs/result"
	"regexp"
	"strings"
	"unicode"
)

func PushMessage(msg string, isErr bool) {
	NewResponse := subresult.ResultSubProgram{Model: subresult.ResultModel{Result: msg, IsError: isErr}}
	go NewResponse.SpawnResultSubProgram() // Render result to ui
}

// parse
func parseInput(inp string) (string, error) {
	if len(inp) >= 200 {
		return "", errors.New("the exercism name can't be that long")
	}
	inp = strings.TrimSpace(inp)
	inp = strings.ToLower(inp)

	ExercismId := inp
	re := regexp.MustCompile("exercism download --exercise=(?:([a-zA-Z]+)|([a-zA-Z]+(-[a-zA-Z]+)+)) --track=go")
	isLink := re.MatchString(ExercismId)

	if isLink {
		IdFound := re.FindStringSubmatch(ExercismId)
		if len(IdFound) < 4 {
			return "", errors.New("cannot exctract exercism name from this link")
		}

		if len(strings.TrimSpace(IdFound[1])) > 0 {
			ExercismId = IdFound[1]
		} else if len(strings.TrimSpace(IdFound[2])) > 0 {
			ExercismId = IdFound[2]
		} else {
			return "", errors.New("cannot exctract exercism name from this link")
		}
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
