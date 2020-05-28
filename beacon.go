package beacon

import (
  "github.com/gorilla/mux"
  "github.com/armadanet/captain/dockercntrl"
  "fmt"
  "strconv"
  "net/http"
  "sync"
  "log"
  "time"
  "os"
)

type Beacon interface {
  Run(port int)
}

type beacon struct {
  router_public       *mux.Router
  router_internal     *mux.Router
  spinners            *SpinnerTable
  state               *dockercntrl.State
  container_name      string
  overlay_name        string
  swarm_token         string
  swarm_ip            string
}

type SpinnerTable struct {
  table         map[string]SpinnerInfo
  mux           sync.Mutex
}

type SpinnerInfo struct {
  Id              string      `json:"Id"`           // unique id of this spinner
  OverlayName     string      `json:"OverlayName"`  // name of the overlay network for this spinner
  LastUpdate      time.Time   `json:"LastUpdate"`   // last timestamp of update
}

func New(containerName string, overlayName string) (*beacon, error) {
  // initiate spinner table
  spinnerTable := SpinnerTable {
    table:make(map[string]SpinnerInfo),
  }
  // initiate docker control state
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  // public server router
  router_public := mux.NewRouter().StrictSlash(true)
  // internal server router
  router_internal := mux.NewRouter().StrictSlash(true)
  // beacon instance
  b := beacon {
    router_public: router_public,
    router_internal: router_internal,
    spinners: &spinnerTable,
    state: state,
    container_name: containerName,
    overlay_name: overlayName,
  }
  // set up router handler
  router_public.HandleFunc("/newCaptain", b.newCaptain()).Name("NewCaptain")
  router_public.HandleFunc("/newSpinner", b.newSpinner()).Name("NewSpinner")
  router_public.HandleFunc("/newTask", b.newTask()).Name("NewTask")
  router_internal.HandleFunc("/register", b.register()).Name("Register")
  // return the beacon instance
  return &b, nil
}

// public_port, containerName, overlayName
func (b *beacon) Run(port int) {
  // beacon create overlay network
  token, ip, err := b.state.BeaconCreateOverlay(b.container_name, b.overlay_name)
  if err != nil {
    log.Println(err)
    return
  }
  b.swarm_token = token
  b.swarm_ip = ip

  // start monitor routine monitoring the spinner table
  go b.monitorSpinnerTable()

  // start the internal beacon server
  go startInternalServer(b.router_internal)

  // start the public beacon server (restful Api)
  fmt.Println("Public server listening on port "+strconv.Itoa(port)+"...")
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), b.router_public))
}

// monitor all spinners
func (b *beacon) monitorSpinnerTable() {
  for {
    b.spinners.mux.Lock()
    now := time.Now()
    for key, spinner := range b.spinners.table {
      temp := spinner.LastUpdate.Add(6*time.Second)
      if temp.Before(now) {
        fmt.Fprintf(os.Stdout, "Spinner %s left system\n",key)
        delete(b.spinners.table, key)
      }
    }
    b.spinners.mux.Unlock()
    time.Sleep(4*time.Second)
  }
}

// new routine start internal server
func startInternalServer(r *mux.Router) {
  // hard code internal server port
  fmt.Println("Internal server listening on port 8787...")
  log.Fatal(http.ListenAndServe(":8787", r))
}
