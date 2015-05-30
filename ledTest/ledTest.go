// ledTest
package main

import (
	"fmt"
	"os"
	"time"
)

func turnOn(f *os.File) {
	on := []byte("1")
	_, err := f.Write(on)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
	}
}

func turnOff(f *os.File) {
	off := []byte("0")
	_, err := f.Write(off)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
	}
}

func main() {

	f, err := os.OpenFile("/sys/class/gpio/gpio48/value", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
	}
	turnOff(f)
	time.Sleep(time.Duration(5) * time.Second)
	turnOn(f)
	f.Close()
}
