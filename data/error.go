package data

import (
	"strings"
)

type Error struct {
	DataSetName string
	Message     string
	Err         error
}

func (err *Error) Error() string {
	msg := []string{
		err.DataSetName,
		err.Message,
	}
	if err.Err != nil {
		msg = append(msg, err.Err.Error())
	}
	return strings.Join(msg, ": ")
}
