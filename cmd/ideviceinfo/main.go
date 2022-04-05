package main

import (
	"fmt"

	"github.com/inahga/inahgo/imobiledevice/idevice"
)

func main() {
	idevice.Debug(1)

	devices, err := idevice.List()
	if err != nil {
		panic(err)
	}
	fmt.Println(devices)

	device, err := idevice.New("")
	if err != nil {
		panic(err)
	}
	fmt.Println(device)

	handle, err := device.Handle()
	if err != nil {
		panic(err)
	}
	fmt.Println(handle)

	udid, err := device.UDID()
	if err != nil {
		panic(err)
	}
	fmt.Println(udid)

	events, err := idevice.EventSubscribe()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := idevice.EventUnsubscribe(); err != nil {
			panic(err)
		}
	}()

	for {
		event := <-events
		fmt.Println(event)
	}
}
