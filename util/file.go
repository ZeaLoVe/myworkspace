package util

import (
	"os"
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
