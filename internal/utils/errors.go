package utils

import "os"

func WriteToStderr(err error) {
	os.Stderr.WriteString(err.Error())
}
