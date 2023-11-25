package main

import (
	"FileServer/internal/logger"
	"FileServer/internal/server"
	"fmt"
	"net"
	"strings"

	"github.com/charmbracelet/log"
)

func main() {
	server := server.NewServer()

	logger.Log.Infof("Starting server at: http://%s%s\n", getOutboundIP(), server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		logger.Log.Fatal("could not start server", "err", err)
	}
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal("could not get IP address", "err", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	stringAddr := fmt.Sprint(localAddr)
	stringAddr, _, _ = strings.Cut(stringAddr, ":")
	stringAddr = strings.ReplaceAll(stringAddr, " ", ".")

	return stringAddr
}
