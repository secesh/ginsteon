//  Copyright 2012 ChaseFox (Matthew R Chase)
//  
//  This file is part of ginsteon, a go library for communicating with
//  Insteon devices.  http://www.insteon.net/
//  
//  ginsteon is free software: you can redistribute it and/or modify
//  it under the terms of the GNU General Public License as published
//  by the Free Software Foundation, either version 3 of the License,
//  or (at your option) any later version.
//  
//  ginsteon is distributed in the hope that it will be useful, but
//  WITHOUT ANY WARRANTY; without even the implied warranty of 
//  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//  GNU General Public License for more details.
//  
//  You should have received a copy of the GNU General Public License
//  along with ginsteon.  If not, see <http://www.gnu.org/licenses/>.

// Package plm provides methods to communicate with an Insteon PLM.
// http://www.insteon.net/
package plm

import (
	"github.com/tarm/goserial"
	"log"
	"encoding/hex"
	"io"
	"time"
)
const(
	plmDelay = 300 * time.Millisecond //a delay used when communicating with the PLM.
)
type Plm struct{
	Port string
	isOpen bool
	port io.ReadWriteCloser
	
	readChannel   chan *readResult
	busyReceiving bool
    writeQueue    []string
}
type readResult struct{
    b []byte
    err error
}

func (plm *Plm) Run(){
	plm.open()
	go plm.listen()
	go plm.masterControl()
}
func (plm *Plm) listen(){
    for{
        buf := make([]byte, 128)
        log.Print("Reading from serial port")
        n, err := plm.port.Read(buf)
        plm.readChannel <- &readResult{buf[:n], nil}
        if err != nil { 
            log.Print("received an error")
            return 
        }
    }
}
func (plm *Plm) masterControl() {
	//woot tron 1.
    received := []byte{}
    
    log.Print("MasterControl has started.")
    for {
    	
    	if(!plm.busyReceiving && len(plm.writeQueue) > 0){
    		plm.write(plm.writeQueue[0])
    		plm.writeQueue = plm.writeQueue[1:]
    		time.Sleep(plmDelay) //give the PLM a chance to start responding.
    	}
    	
    	
        timeout := time.NewTicker(plmDelay) //if we hit this with no data received, we 
                                            //can assume the PLM has finished sending data.
        defer timeout.Stop() // 1) Is this necessary? 2) is this local to the for{}?
        select{
        case got := <-plm.readChannel:
            //log.Print("got a result")
            switch{
            case got.err != nil:
                //Catching an EOF error here can indicate the port was disconnected.
                // -- if using a USB to serial port, and the device is unplugged 
                //    while being read, we'll receive an EOF.
                log.Fatal("  error:" + got.err.Error())
                plm.busyReceiving = false
            default:
                plm.busyReceiving = true
                received = append(received, got.b...)
                //log.Print(got.b)
            }
        case <-timeout.C:
        	
            //stop waiting for the reader to send something on channel rc
            plm.busyReceiving = false
            if(len(received)>0){
                log.Print("Received something from PLM")
                log.Print(received)
                received = []byte{}
            }
        }
    }

}
func (plm *Plm) open() (b bool){
	log.Print("Opening serial port connection to PLM")
	plm.readChannel   = make(chan *readResult)
	plm.busyReceiving = false

	var err error
    c := &serial.Config{Name: plm.Port, Baud: 19200}
    plm.port, err = serial.OpenPort(c)
    if err != nil { 
    	log.Fatal(err) 
   	}else{
   		plm.isOpen = true
   	}
   	return plm.isOpen
}
func (plm *Plm) write(s string){
	// if(!plm.isOpen && !plm.open()){
	// 	//won't get here because there's a fatal in open()
	// 	return
	// }
	var err error
	log.Print("Writing to PLM: " + s)
    hex, _ := hex.DecodeString(s)
    _, err = plm.port.Write(hex)
    if err != nil { 
    	plm.isOpen = false
        log.Fatal("Could not write to PLM: " + err.Error())
   	}
}
func (plm *Plm) Write(s string){
	plm.writeQueue = append(plm.writeQueue, s)
}
