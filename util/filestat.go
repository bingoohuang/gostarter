package util

import "os"

type ExistsEnum int

const (
	Exists ExistsEnum = iota
	NotExists
	Unknown
)

func FileStat(name string) (ExistsEnum, error) {
	_, err := os.Stat(name)
	if err == nil {
		return Exists, nil
	}

	if os.IsNotExist(err) {
		return NotExists, nil
	}

	return Unknown, err
}
