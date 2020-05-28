package beacon

import (
  "fmt"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "log"
  "os"
)


// new captain call this to find spinner to join
func (b *beacon) newCaptain() func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    responseBody, err := json.Marshal(map[string]interface{} {
      "Number":20,
      "Text":"asdf",
    })
    if err != nil {
  		fmt.Println(err)
      return
  	}
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
  }
}

// spinner first call this (create new overlay)
func (b *beacon) newSpinner() func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    // read: new spinner id (name)
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    var req struct {
      SpinnerId string
    }
    err = json.Unmarshal(body, &req)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    fmt.Fprintf(os.Stdout, "New Spinner Joined! SpinnerId: %s\n",req.SpinnerId)
    // create overlay network for this spinner id_spinner
    err = b.state.BeaconCreateSpinnerOverlay(req.SpinnerId+"_overlay")
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    // assemble response body
    responseBody, err := json.Marshal(map[string]string{
      "SwarmToken": b.swarm_token,
      "BeaconIp": b.swarm_ip,
      "BeaconOverlay": b.overlay_name,
      "BeaconName": b.container_name,
      "SpinnerOverlay": req.SpinnerId+"_overlay",
    })
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(responseBody)
  }
}

// user call this to find a spinner to submit the task
func (b *beacon) newTask() func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    return
  }
}

// spinner periodically ping this to notify alive
func (b *beacon) register() func(http.ResponseWriter, *http.Request) {
  return func(w http.ResponseWriter, r *http.Request) {
    // read: spinner_id(name),
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    var req SpinnerInfo
    err = json.Unmarshal(body, &req)
    if err != nil {
      log.Println(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
    }
    // update spinner into SpinnerTable
    b.spinners.mux.Lock()
    b.spinners.table[req.Id] = req
    b.spinners.mux.Unlock()
  }
}
