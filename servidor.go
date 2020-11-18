package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

func Proceso(id uint64, i uint64, channel chan uint64, quit chan bool) {
	for {
		select {
		case <-quit:
			channel <- id
			channel <- i
			return
		default:
			fmt.Println(id, " : ", i)
			//channels <- strconv.FormatUint(id, 10) + ": " + strconv.FormatUint(i, 10)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func manageProcess(channels []chan uint64, chansQuit []chan bool, add chan uint64, send chan net.Conn, receive chan bool) {
	for {
		select {
		case <- receive:
			inf := [2]uint64{<-add, <- add}
			channels = append(channels, make(chan uint64, 2))
			chansQuit = append(chansQuit, make(chan bool))
			go Proceso(inf[0], inf[1], channels[len(channels)-1], chansQuit[len(chansQuit)-1])
		case c := <- send:
			go handleClient(c, channels[0], chansQuit[0], add, receive)
			channels = channels[1:]
			chansQuit = chansQuit[1:]
		}
	}
}


func servidor(channels []chan uint64, chansQuit []chan bool) {
	add := make(chan uint64, 2)
	send := make(chan net.Conn)
	receive := make(chan bool)
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	go manageProcess(channels, chansQuit, add, send, receive)
	for {
		c, err := s.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}	
		send <- c	
	}
}

func handleClient(c net.Conn, channel chan uint64, quit chan bool, add chan uint64, receive chan bool) {

	quit <- true
	inf := [2]uint64{<-channel, <- channel}
	err := gob.NewEncoder(c).Encode(inf)
	if err != nil {
		fmt.Println(err)
	} 
	for {
		err := gob.NewDecoder(c).Decode(&inf[1])
		if err != nil {
			add <- inf[0]
			add <- inf[1]
			receive <- true
			return
		} 
	}
}

func main() {
	var idCount uint64 = 0
	channels := make([]chan uint64, 5)
	chansQuit := make([]chan bool, 5)
	for i:= 0; i < 5; i++{
		channels[i] = make(chan uint64, 2)
		chansQuit[i] = make(chan bool)
		go Proceso(idCount, 0, channels[i], chansQuit[i])
		idCount++
	}
	go servidor(channels, chansQuit)
	var input string
	fmt.Scanln(&input)

}
