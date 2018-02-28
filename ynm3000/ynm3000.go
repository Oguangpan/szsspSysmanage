/*
2018.02.28 by Deadpig panndora
The value of this program is no value. Szssp belongs to the company
Well, its name is called YNM3000.
Property value should be:
   - HDD ID
   - Device model
   - mac address
   - user
   - Department
The method should be:
   - Get the hard disk serial number
   - Get mac address
   - search xlsx file mac address
   - Add or modify xlsx file contents
   - Set the ip address
*/
package ynm3000

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type Ynm3k interface {
	GetHardNumber()
	GetMacAddress()
	//SearchXlsx()
	//SetIpaddress()
	//AmXlsx()
}

type Computer struct {
	hddId      []string
	model      string
	macAddress []intfInfo
	user       string
	department string
}

// Get the hard disk physical serial number information
func (cp *Computer) GetHardNumber() {
	ids_byte, _ := exec.Command("cmd", "/C", "wmic diskdrive get serialnumber").Output()
	ids := string(ids_byte)
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		cp.hddId = append(cp.hddId, j)
	}
	fmt.Println(cp.hddId)
}

// Get the Netcard information
type intfInfo struct {
	Name string
	Mac  net.HardwareAddr
}

func (cp *Computer) GetMacAddress() {
	intf, err := net.Interfaces()
	if err != nil {
		return
	}
	var tmp intfInfo
	for _, v := range intf {
		tmp.Name = v.Name
		tmp.Mac = v.HardwareAddr
		cp.macAddress = append(cp.macAddress, tmp)
	}
	// test code
	for _, v := range cp.macAddress {
		fmt.Println(v)
	}
}
