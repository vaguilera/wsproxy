// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var remoteHost string = ""
var messageType int = websocket.TextMessage

type WSConn struct {
	wsConn *websocket.Conn
	a      int
}

type TCPConn struct {
	conn net.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
} // use default options

func (conn *WSConn) Upgrade(w http.ResponseWriter, r *http.Request) error {
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	log.Print("WS connection from " + clientIP)
	var err error
	conn.wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading connection:", err)
		return err
	}
	return nil
}

func (conn *WSConn) Close() {
	conn.wsConn.Close()
}

func (conn *WSConn) ReadData(ch chan []byte, cherror chan bool) {
	for {
		_, msg, err := conn.wsConn.ReadMessage()
		if err != nil {
			log.Print("Error reading from WS", err)
			cherror <- true
			return
		}
		ch <- msg
	}
}

func (conn *WSConn) SendData(data []byte) error {
	err := conn.wsConn.WriteMessage(messageType, data)
	if err != nil {
		log.Print("Error sending data to WS Socket", err)
	}
	return err
}

func (conn *TCPConn) Connect(host string) error {
	var err error
	conn.conn, err = net.Dial("tcp", host)

	if err != nil {
		log.Print("Error connecting remote host", err)
		return err
	}
	log.Print("Connected to remote host", host)
	return nil
}

func (conn *TCPConn) ReadData(ch chan []byte, cherror chan bool) {

	buffer := make([]byte, 1024)
	for {
		nbytes, err := bufio.NewReader(conn.conn).Read(buffer)
		if err != nil {
			log.Print("Error reading from TCP socket", err)
			cherror <- true

		}
		if nbytes > 0 {
			ch <- buffer[0:nbytes]
		}
	}
}

func (conn *TCPConn) SendData(data []byte) error {
	_, err := conn.conn.Write(data)
	if err != nil {
		log.Print("Error sending data to TCP Socket")
	}
	return err
}

func (conn *TCPConn) Close() {
	conn.conn.Close()
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	wsSocket := &WSConn{}
	err := wsSocket.Upgrade(w, r)
	if err != nil {
		return
	}

	tcpSocket := &TCPConn{}
	err = tcpSocket.Connect(remoteHost)
	if err != nil {
		wsSocket.Close()
		return
	}

	chWS := make(chan []byte)
	chTCP := make(chan []byte)
	chError := make(chan bool)

	go wsSocket.ReadData(chWS, chError)
	go tcpSocket.ReadData(chTCP, chError)
	for {
		select {
		case dataWS := <-chWS:
			err = tcpSocket.SendData(dataWS)
			if err != nil {
				wsSocket.Close()
				tcpSocket.Close()
				return
			}
		default:
		}

		select {
		case dataTCP := <-chTCP:
			err := wsSocket.SendData(dataTCP)
			if err != nil {
				wsSocket.Close()
				tcpSocket.Close()
				return
			}
		default:
		}

		select {
		case status := <-chError:
			if status {
				wsSocket.Close()
				tcpSocket.Close()
				return
			}
		default:
		}
	}

}

func ListenWebSocket(localhost string, localport int, remoteAddress string, remotePort int, textMode bool) {

	listenPort := strconv.Itoa(localport)
	remoteHost = remoteAddress + ":" + strconv.Itoa(remotePort)
	if !textMode {
		messageType = websocket.BinaryMessage
	}

	http.HandleFunc("/", handleConnection)
	log.Print("Starting proxy on [ws:" + listenPort + "]")
	log.Fatal(http.ListenAndServe("localhost:"+listenPort, nil))
}
