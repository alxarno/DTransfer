package main

import (
	"log"
	"net/http"

	"github.com/AlexeyArno/golang-files-transfer/src/gui"
	"github.com/AlexeyArno/golang-files-transfer/src/network"
	webSocketWork "github.com/AlexeyArno/golang-files-transfer/src/webSocketWork"
	"github.com/zserge/webview"
	"golang.org/x/net/websocket"
)

var close = make(chan struct{})

func main() {
	data, err := gui.Asset("data/index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	})
	http.Handle("/websocket_data", websocket.Handler(webSocketWork.DataHandler))

	port, err := network.ReservTCPPort()
	if err != nil {
		panic(err)
	}
	log.Println("Listen TCP port: ", port)
	webSocketWork.RegisterTCPPort(port)

	go func(port string) {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			panic(err)
		}
	}(port)

	// close <- struct{}{}
	network.StartServer()
	w := webview.New(webview.Settings{
		Title:                  "DTransfer",
		URL:                    "http://127.0.0.1:" + port,
		Width:                  300,
		Height:                 430,
		ExternalInvokeCallback: gui.HandleRPC})
	gui.RegisterGUI(&w)
	defer w.Exit()
	w.Run()
	<-close

}
