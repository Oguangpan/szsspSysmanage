/*
2018.02.28 by dead pig panndora

Process:
1. Get a list of hard drive addresses.
2. Get mac address list.
3. visit xlsx file search mac address, return the relevant data.
4. According to the data set ip address.
5. Prompt and modify the contents of the xlsx file according to the operation.
*/

package main

import (
	"github.com/Oguangpan/szsspSysManage/ynm3000"
)

// temp ip
const tmpIp = "33.66.100.255"

func main() {
	var com ynm3000.Ynm3k
	com = new(ynm3000.Computer)
	com.GetHardNumber()
	com.GetMacAddress()
}
