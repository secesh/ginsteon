Go library for Insteon devices
==============================

This library provides communication with Insteon devices in the Go Language.
http://www.insteon.net/

The goal of this project is to make communication with an Insteon network easy
for go developers.  Insteon is a home automation technology well suited for
adoption by enthusiasts.

This project is in infancy and as yet only demonstrates a proof of concept
communication with a PLM.

Install this package:
---------------------

    go get github.com/tarm/goserial
    go get github.com/secesh/ginsteon

Example:
--------

````go
package main

import(
    "github.com/secesh/ginsteon/plm"
    "log"
    "runtime"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    plm := plm.Plm{Port: "/dev/ttyUSB0"}
    log.Print("about to start")
    plm.Run()
    
    plm.Write("0260")
    plm.Write("0260")
    plm.Write("0260")
    plm.Write("0260")
    
    for {}
}
````