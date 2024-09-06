package api

import "strings"

func ParsePipelineMessage(message string) (string, string) {
	errorString := message
	messageLine := errorString

	if strings.HasPrefix(errorString, "jobs:") {
		messageLine = errorString[5:]
	}

	parts := strings.SplitN(messageLine, " ", 2)
	word := parts[0]
	job := word

	if strings.Index(word, ":") == -1 {
		job = word
	} else {
		job = word[:strings.Index(word, ":")]
	}

	errorMessage := messageLine[len(job)+1:]

	return job, errorMessage
}