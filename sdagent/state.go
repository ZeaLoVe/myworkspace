package sdagent

import (
	"encoding/json"
	"net/http"
)

type Statistic struct {
	Running   int        `json:"RunningService,omitempty"`
	Ready     int        `json:"ReadyService,omitempty"`
	Prepare   int        `json:"PrepareService,omitempty"`
	Sum       int        `json:"RegisterService,omitempty"`
	JobStates []JobState `json:"JobStates,omitempty"`
}

func (sda *SDAgent) RegisterHandle(w http.ResponseWriter, r *http.Request) {
	res, err := json.Marshal(sda)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(res)
	}
}

func (sda *SDAgent) ServiceHandle(w http.ResponseWriter, r *http.Request) {
	var res string
	for i, _ := range sda.Jobs {
		tmp, _ := json.Marshal(sda.Jobs[i].S)
		res = res + string(tmp)
	}
	w.Write([]byte(res))
}

func (sda *SDAgent) JobHandle(w http.ResponseWriter, r *http.Request) {
	var res string
	for i, _ := range sda.Jobs {
		tmp, _ := json.Marshal(sda.Jobs[i].state)
		res = res + string(tmp)
	}
	w.Write([]byte(res))
}

func (sda *SDAgent) StatisticHandle(w http.ResponseWriter, r *http.Request) {
	var running, ready, prepare, sum int
	sum = len(sda.Jobs)
	for i, _ := range sda.Jobs {
		switch sda.Jobs[i].config.JOBSTATE {
		case RUNNING:
			running += 1
		case PREPARE:
			prepare += 1
		case READY:
			ready += 1
		}
	}
	var state Statistic
	state.Running = running
	state.Ready = ready
	state.Prepare = prepare
	state.Sum = sum
	res, err := json.Marshal(state)
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(res)
	}
}
