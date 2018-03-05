/*
2018.02.28 by Deadpig panndora
The value of this program is no value. Szssp belongs to the company
Well, its name is called YNM3000.
Property value should be:
   - HDD ID
   - Ip
   - mac address
   - User
   - Department
The method should be:
   - Get the hard disk serial number
   - Get mac address
   - search xlsx file mac address
   - Add or modify xlsx file contents
   - Set the Ip address
*/
package ynm3000

import (
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/tealeg/xlsx"
)

const xlsxFilePath = "//33.66.96.14/public/2018Taizhang.xlsx"

type Ynm3k interface {
	GetHardNumber()
	GetMacAddress()
	SearchXlsx() (int, string)
	//SetIpaddress()
	//AmXlsx()
}

type Computer struct {
	HddIds     []string            // All hard disks read on this device
	HddId      string              // The hard disk in the form file
	Ip         string              //The IP address stored on the form
	MacAddress []map[string]string //This machine all the Ip address
	MacAddres  string
	User       string
	Department string
}

// Get the hard disk physical serial number information
func (cp *Computer) GetHardNumber() {
	ids_byte, _ := exec.Command("cmd", "/C", "wmic diskdrive get serialnumber").Output()
	ids := string(ids_byte)
	var slicel []string = strings.Fields(ids)[1:]
	for _, j := range slicel {
		cp.HddIds = append(cp.HddIds, j)
	}

}

// Get the Netcard information

func (cp *Computer) GetMacAddress() {
	intf, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, v := range intf {
		tmp := make(map[string]string)
		tmp["Name"] = v.Name
		tmp["Mac"] = v.HardwareAddr.String()
		if tmp["Mac"] != "" {
			cp.MacAddress = append(cp.MacAddress, tmp)
		}
	}
	//fmt.Println(cp.MacAddress)
}

// Query  mac address from xlsx related data
// Return the number of rows, device data, matching the name of the network card
func (cp *Computer) SearchXlsx() (int, string) {
	xlFile, err := xlsx.OpenFile(xlsxFilePath)
	if err != nil {
		fmt.Println(err)
		return 0, ""
	}

	for _, sheet := range xlFile.Sheets {
		if sheet.Name == "计算机" {
			for column, row := range sheet.Rows {
				for _, j := range cp.MacAddress {
					for _, cell := range row.Cells {
						if cell.String() == j["Mac"] {

							/*The following switch to determine the basis of the column in the xlsx table information,
							so you want to directly use the package file you also need the corresponding xlsx file,
							or need to be changed as required.*/
							for i, cell := range row.Cells {
								switch i {
								case 2:
									cp.Department = cell.String()
								case 3:
									cp.User = cell.String()
								case 9:
									cp.HddId = cell.String()
								case 10:
									cp.Ip = cell.String()
								case 11:
									cp.MacAddres = cell.String()
								}
							}

							return column, j["Name"]
						}
					}
				}
			}
		}
	}
	return 0, ""
}

func (cp *Computer) AmXlsx(s bool) (err error) {

	return nil
}

func (cp *Computer) SetIpaddress() {
	// cp.Ip
}
