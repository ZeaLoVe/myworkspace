package sdagent

import (
	. "myworkspace/service"
)

const Version = "0.1"

type SDAgent struct {
	Ser []Service `json:"services,omitempty"`

	jobs []Job `json:"-"`
}

func (sda *SDAgent) Run() {

}

func (sda *SDAgent) LoadConfig() {

}

func (sda *SDAgent) Start() {

}

func (sda *SDAgent) StopAll() {

}

func (sda *SDAgent) StopJob() {

}

func (sda *SDAgent) startJob() {

}

func main() {

}
