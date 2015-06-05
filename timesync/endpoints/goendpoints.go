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
import "strconv"
import "strings"

var gClockOffset time.Duration
var gDelayRTT time.Duration
var timeStampLayout = "02:Jan:2006:15:04:05.000000"
var samples = 100

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

func exchangeTimestamps(f *os.File) uint64 {
	// client = coordinator, server = endpoints. We are server
	// t0 is the client's timestamp of the request packet transmission,
	// t1 is the server's timestamp of the request packet reception,
	t0 := readMessage(f)
	t1 := generateTimestamp()

	// t2 is the server's timestamp of the response packet transmission and
	// t3 is the client's timestamp of the response packet reception.
	t2 := generateTimestamp()
	writeMessage(f, t2)
	t3 := readMessage(f)
	// gotta send t1
	writeMessage(f, t1)
	//calculateClockOffset(t0, t1, t2, t3)
	return calculateDelayRTT(t0, t1, t2, t3)
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

func waitTimer(cycles uint64, f *os.File) int {
	// init counters:
	C.init_perfcounters(1, 0)
	//fmt.Printf("cyles to wait: %d\n", cycles)
	timeStart := C.ulonglong(C.get_cyclecount())
	timeElapsed := C.ulonglong(0)
	for {
		timeElapsed = C.ulonglong(C.get_cyclecount()) - timeStart
		if timeElapsed > C.ulonglong(cycles) {
			writeMessage(f, fmt.Sprintf("%27s", "COMPLETE"))
			break
		}
	}
	return int(timeElapsed)
}

func priority_test(f *os.File, delay uint64) {
	rtt := delay
	str := readMessage(f)
	cycles, err := strconv.ParseUint(strings.TrimSpace(str), 0, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("readMessage string: %s\n", str)
	fmt.Printf("rtt: %d\n", rtt)
	fmt.Printf("cycles: %d\n", cycles)
	waitTimer(cycles-rtt-rtt, f)
	//writeMessage(f, fmt.Sprintf("%27s", "COMPLETE"))
}

func calculateDelayJitter(f *os.File) (aDelay uint64, aJitter uint64) {
	samples := 100
	delayEntries := make([]uint64, samples)
	// Delay/Jitter Calculations
	totalDelay := uint64(0)
	for i := 0; i < samples; i++ {
		delayEntries[i] = exchangeTimestamps(f)
		totalDelay += delayEntries[i]
		//fmt.Printf("RTT delay: %v\n", delayEntries[i])
		//fmt.Printf("Running total: %v\n", totalDelay);
	}
	avgdelay := totalDelay / uint64(samples)
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
	avgjitter := totaljitter / uint64(samples)
	return avgdelay, avgjitter
}

func main() {
	f, err := openPort("/dev/ttyUSB0")
	if err != nil {
		fmt.Printf("Failed to open the serial port!")
	}
	avgDelay, avgJitter := calculateDelayJitter(f)
	fmt.Printf("Avg. Delay = %v ns\n", avgDelay)
	fmt.Printf("Avg. Jitter = %v ns\n", avgJitter)
	priority_test(f, avgDelay)
	f.Close()
}
