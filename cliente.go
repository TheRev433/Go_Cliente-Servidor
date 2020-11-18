package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

func ProcesoCliente(id uint64, i uint64, channel chan uint64) {
	for {
		fmt.Println(id, " : ", i)
		i = i + 1
		channel <- i
		time.Sleep(time.Millisecond * 500)
	}
}

func Cliente() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	var inf [2]uint64
	err = gob.NewDecoder(c).Decode(&inf)
	if err != nil {
		fmt.Println(err)
	} else {
		channel := make(chan uint64)
		go ProcesoCliente(inf[0], inf[1], channel)
		for {
			inf[1] = <- channel
			err := gob.NewEncoder(c).Encode(inf[1])
			if err != nil {
				fmt.Println(err)
				return
			} 
		}
	}
}

func main() {
	go Cliente()
	var input string
	fmt.Scanln(&input)
}
