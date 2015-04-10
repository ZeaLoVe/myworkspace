package service

import (
	"testing"
)

func TestChangeMachine(t *testing.T) {
	var r Register
	var ser Service
	ser.SetDefault()
	r.s = &ser
	newMachine := []string{"http://127.0.0.1:2379"}
	r.ChangeMachines(newMachine)

	var EmptyMachine []string
	r.ChangeMachines(EmptyMachine)
}

func TestUpdateService(t *testing.T) {
	var r Register
	var ser Service
	ser.SetDefault()
	r.s = &ser
	r.UpdateService()
}
