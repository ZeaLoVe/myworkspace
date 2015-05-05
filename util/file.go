package util

import (
	"os"
	"path"
	"strconv"
	"time"
)

var MODIFYINTERVAL int

func init() {
	MODIFYINTERVAL = 30
}

func CheckModify(filename string) (bool, error) {
	if fi, err := os.Stat(filename); err != nil {
		return false, err
	} else {
		if time.Since(fi.ModTime()).Minutes() < float64(MODIFYINTERVAL) {
			return true, nil
		} else {
			return false, nil
		}
	}

}

func GenPidFile(filepath, filename string) error {
	var target string
	if filename == "" {
		filename = "sdagent.pid"
	}
	if filepath == "" {
		target = filename
	} else {
		target = path.Join(filepath, filename)
	}
	file, err := os.Create(target)
	if err != nil {
		return err
	} else {
		file.Write([]byte(strconv.Itoa(os.Getpid())))
	}
	return nil
}
