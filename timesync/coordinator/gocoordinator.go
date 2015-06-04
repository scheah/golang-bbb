package main

// #include <stdio.h>
// #include <stdint.h>
// unsigned int get_cyclecount (void)
// {
//  unsigned int value;
//	// Read CCNT Register
//	asm volatile ("MRC p15, 0, %0, c9, c13, 0\t\n": "=r"(value));
//	return value;
// }
//
// void init_perfcounters (int32_t do_reset, int32_t enable_divider)
// {
//	// in general enable all counters (including cycle counter)
//	int32_t value = 1;
//
//	// peform reset:
//	if (do_reset)
//	{
//		value |= 2;     // reset all counters to zero.
//		value |= 4;     // reset cycle counter to zero.
//	}
//
//	if (enable_divider)
//		value |= 8;     // enable "by 64" divider for CCNT.
//
//	value |= 16;
//
//	// program the performance-counter control-register:
//	asm volatile ("MCR p15, 0, %0, c9, c12, 0\t\n" :: "r"(value));
//
//	// enable all counters:
//	asm volatile ("MCR p15, 0, %0, c9, c12, 1\t\n" :: "r"(0x8000000f));
//
//	// clear overflows:
//	asm volatile ("MCR p15, 0, %0, c9, c12, 3\t\n" :: "r"(0x8000000f));
// }
import "C"
import "fmt"
import "os"
import "syscall"
import "unsafe"
import "time"

var timeStampLayout = "02:Jan:2006:15:04:05.000000"
var samples = 100
var delayEntries = make([]uint64, samples)

func openPort(name string) (f *os.File, err error) {

	f, err = os.OpenFile(name, syscall.O_RDWR|syscall.O_NOCTTY, 0666)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil && f != nil {
			f.Close()
		}
	}()

	fd := f.Fd()

	// Set serial port 'name' to 115200/8/N/1 in RAW mode (i.e. no pre-process of received data
	// and pay special attention to Cc field, this tells the serial port to not return until at
	// at least syscall.VMIN bytes have been read. This is a tunable parameter they may help in Lab 3
	t := syscall.Termios{
		Iflag:  syscall.IGNPAR,
		Cflag:  syscall.CS8 | syscall.CREAD | syscall.CLOCAL | syscall.B115200,
		Cc:     [32]uint8{syscall.VMIN: 27},
		Ispeed: syscall.B115200,
		Ospeed: syscall.B115200,
	}

	// Syscall to apply these parameters
	_, _, errno := syscall.Syscall6(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(syscall.TCSETS),
		uintptr(unsafe.Pointer(&t)),
		0,
		0,
		0,
	)

	if errno != 0 {
		return nil, errno
	}

	return f, nil
}

func generateTimestamp() string {
	t := time.Now()
	return t.Format(timeStampLayout)
}

func parseTimestamp(timestamp string) time.Time {
	t, e := time.Parse(timeStampLayout, timestamp)
	if e != nil {
		fmt.Printf("Parse error occured: %v\n", e)
	}
	return t
}

func calculateDelayRTT(p0 string, p1 string, p2 string, p3 string) uint64 {
	// parse time stamp string
	t0 := parseTimestamp(p0)
	t1 := parseTimestamp(p1)
	t2 := parseTimestamp(p2)
	t3 := parseTimestamp(p3)
	delayRTT := (t3.Sub(t0) + t2.Sub(t1))
	//fmt.Printf("RTT delay: %v\n", delayRTT)
	return uint64(delayRTT.Nanoseconds())
}

func readMessage(f *os.File) string {
	buffer := make([]byte, 27)
	count, err := f.Read(buffer)
	if err != nil {
		fmt.Printf("Read error occured: %v\n", err)
	}
	fmt.Printf("Received %d bytes: %s\n", count, string(buffer))
	return string(buffer)
}

func writeMessage(f *os.File, str string) {
	buffer := []byte(str)
	count, err := f.Write(buffer)
	if err != nil {
		fmt.Printf("Write error occured: %v\n", err)
	}
	fmt.Printf("Sent %d bytes: %s\n", count, str)
}

func exchangeTimestamps(f *os.File) uint64 {
	t0 := generateTimestamp()
	writeMessage(f, t0)
	t2 := readMessage(f)
	t3 := generateTimestamp()
	writeMessage(f, t3)
	// gotta get t1...
	t1 := readMessage(f)
	calculateClockOffset(t0, t1, t2, t3)
	return calculateDelayRTT(t0, t1, t2, t3)
}

//func waitForEvent(c2 chan string, f *os.File, cycles int) {
//writeMessage(f, fmt.Sprintf("%27d", cycles))
//readMessage(f)
//c2 <- "Success: Response within acceptance window\n"
//}

func waitForEvent(c2 chan string, f *os.File) {
	readMessage(f)
	c2 <- "Success: Response within acceptance window\n"
}

func upperWindow(c1 chan string, cycles uint64) {
	C.init_perfcounters(1, 0)
	timeStart := C.ulonglong(C.get_cyclecount())
	for {
		timeElapsed := C.ulonglong(C.get_cyclecount()) - timeStart
		if timeElapsed > C.ulonglong(cycles) {
			c1 <- "MISSED WINDOW\n"
			break
		}
	}
}

func lowerWindow(c1 chan string, cycles uint64) {
	C.init_perfcounters(1, 0)
	timeStart := C.ulonglong(C.get_cyclecount())
	for {
		timeElapsed := C.ulonglong(C.get_cyclecount()) - timeStart
		if timeElapsed > C.ulonglong(cycles) {
			c1 <- "ENTER WINDOW\n"
			break
		}
	}
}

func priority_test(usb *os.File, windowSize int, waitTime int, led_blue *os.File, led_green *os.File) {
	c1 := make(chan string, 1)
	c2 := make(chan string, 1)
	go upperWindow(c1, windowSize)
	go lowerWindow(c1, waitTime)
	writeMessage(usb, fmt.Sprintf("%27d", waitTime))
	go waitForEvent(c2, usb)
	select {
	case res := <-c1:
		fmt.Println(res)
	case <-c2:
		fmt.Println("received message too early.. FAILED")
		return
	}
	select {
	case res := <-c1:
		fmt.Println(res)
	case res := <-c2:
		turnOff(led_blue)
		turnOn(led_green)
		fmt.Println(res)
	}
}

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
	usb, err := openPort("/dev/ttyUSB0")
	if err != nil {
		fmt.Printf("Failed to open the serial port!")
		return
	}
	//pin 12 and pin 15 from P9 on BBB
	//pin 12 = gpio60
	//pin 15 = gpio48
	led_blue, err := os.OpenFile("/sys/class/gpio/gpio48/value", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		return
	}
	led_green, err := os.OpenFile("/sys/class/gpio/gpio60/value", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("error occured: %v\n", err)
		return
	}
	turnOn(led_blue)
	turnOff(led_green)

	// Delay/Jitter Calculations
	totalDelay := uint64(0)
	for i := 0; i < samples; i++ {
		delayEntries[i] = exchangeTimestamps(usb)
		totalDelay += delayEntries[i]
	}
	avgdelay := totaldelay / samples
	totaljitter := uint64(0)
	jitter := uint64(0)
	for i := 0; i < samples; i++ {
		if delayEntries[i] > avgdelay {
			jitter = delayEntries[i] - avgdelay
		} else {
			jitter = avgdelay - delayEntries[i]
		}
		totaljitter += jitter
	}
	avgjitter := totaljitter / samples
	fmt.Printf("Avg. Delay = %v ns\n", avgdelay)
	fmt.Printf("Avg. Jitter = %v ns\n", avgjitter)

	windowSize := 600000000 //cycles at 1ghz is 0.6s
	waitTime := 500000000   //cycles at 1ghz is 0.5s
	priority_test(usb, windowSize, waitTime, led_blue, led_green)
	usb.Close()
	led_blue.Close()
	led_green.Close()
}
