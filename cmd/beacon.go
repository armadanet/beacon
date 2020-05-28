package main

import (
  "github.com/armadanet/beacon"
  //"os"
)

func main() {
  // 9898 is the open port to outside
  bea, err := beacon.New("beacon", "beacon_overlay")
  if err != nil {panic(err)}
  bea.Run(9898)
}
