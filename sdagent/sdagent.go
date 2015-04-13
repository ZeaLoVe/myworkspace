package sdagent

import (
	"log"
)

const Version = "0.1"

type SDAgent struct {
	Ser []service `json:"services,omitempty"`

	jobs []Job `json:"-"`
}

func (sda *SDAgent) Run() {

}

func (sda *SDAgent) AddJob() {

}

func (sda *SDAgent) DeleteJob() {

}

func main() {

}
