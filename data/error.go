package data

import (
	"strings"
)

// Error objects are returned if construction of a data set fails.
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
