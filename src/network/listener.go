package network

import (
	"fmt"
	"log"
	"net"

	"github.com/AlexeyArno/golang-files-transfer/src/utility"
	webSocketWork "github.com/AlexeyArno/golang-files-transfer/src/webSocketWork"
)

var (
	usefullIps []string
)

// Listen return hello
func Listen(myAddressWPort string, port string, ipFound func(string), messages *chan string) {

	ServerAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		log.Println("Listen 1:", err)
		return
	}

	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	if err != nil {
		log.Println("Listen 2:", err)
		return
	}
	defer ServerConn.Close()
	log.Println("Listener UDP address - ", ServerConn.LocalAddr().String())

	buf := make([]byte, 1024)
	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)
		cleanIP := utility.GetCleanIPFromString(addr.String())
		fmt.Println(cleanIP)
		if cleanIP != myAddressWPort+":"+port {
			if !stringInSlice(string(buf[0:n]), usefullIps) {
				gotUsefullIP(cleanIP)
			}
		}
		if err != nil {
			log.Println("Error: ", err)
		}
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func gotUsefullIP(ip string) {

	IPFound(ip)
	for _, port := range AvailabalePortsTCP {
		webSocketWork.ConnectTo(ip + ":" + port)
	}
}