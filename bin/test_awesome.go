package main

import (
	"com"
	"com/tcp"
	"strconv"
	"time"
)

const maxNumberOfClients int = 10

var (
	send_ch    chan tcp.IDable
	receive_ch chan interface{}
	status_ch  chan tcp.ClientStatus
	active     bool
	localID    int
	clients    []*Client
)

type Client struct {
	Id     int
	Active bool
}

func main() {
	clients = make([]*Client, 0, maxNumberOfClients)
	send_ch = make(chan tcp.IDable, 1)
	receive_ch = make(chan interface{})
	status_ch = make(chan tcp.ClientStatus, 1)
	wait := make(chan bool)
	localID, _ = com.Init(send_ch, receive_ch, status_ch, maxNumberOfClients)
	// if err != nil {
	// 	println(err.Error())
	// }
	//localID = stat.ID
	go status_listener()
	go receive()
	<-wait
}

func send(client_ *Client) {
	message := com.Header{1, localID, client_.Id, com.ElevUpdate{4, 0}}
	for {
		//	println("send:" + strconv.FormatBool(client_.Active))
		if client_.Active {
			println("data sendt")
			send_ch <- message
		}
		time.Sleep(1 * time.Second)
	}
}

func receive() {
	for {
		message := <-receive_ch
		m := message.(com.Header)
		switch data2 := m.Data.(type) {
		case com.ElevUpdate:
			println("data mottatt:" + data2.String())
		default:
			println("default")
		}
	}
}
func status_listener() {
	var exists bool
	for {
		exists = false
		status := <-status_ch
		for n, _ := range clients {
			if status.ID == clients[n].Id {
				//			println(strconv.FormatBool(clients[n].Active))
				clients[n].Active = status.Active
				//			println(strconv.FormatBool(clients[n].Active))
				exists = true
				break
			}
		}
		if status.ID == localID {
			println("er master:" + strconv.FormatBool(status.IsMaster))

		} else if !exists {
			client_ := Client{status.ID, status.Active}
			clients = append(clients, &client_)
			go send(&client_)
		}
	}
}
