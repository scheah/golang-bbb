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
import "sort"
import "log"

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

type int64arr []int64
func (a int64arr) Len() int { return len(a) }
func (a int64arr) Swap(i, j int){ a[i], a[j] = a[j], a[i] }
func (a int64arr) Less(i, j int) bool { return a[i] < a[j] }

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

<<<<<<< Updated upstream
func calculateDelayRTT(p0 string, p1 string, p2 string, p3 string) uint64 {
=======
func calculateDelayRTT(p0 string, p1 string, p2 string, p3 string) int64 {
>>>>>>> Stashed changes
	// parse time stamp string
	t0 := parseTimestamp(p0)
	t1 := parseTimestamp(p1)
	t2 := parseTimestamp(p2)
	t3 := parseTimestamp(p3)
	delayRTT := (t3.Sub(t0) + t2.Sub(t1))
	//fmt.Printf("RTT delay: %v\n", delayRTT)
<<<<<<< Updated upstream
	return uint64(delayRTT.Nanoseconds())
=======
	return int64(delayRTT.Nanoseconds())
>>>>>>> Stashed changes
}

func exchangeTimestamps(f *os.File) int64 {
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
	_, err := f.Read(buffer)
	if err != nil {
		fmt.Printf("Read error occured: %v\n", err)
	}
	//fmt.Printf("Received %d bytes: %s\n", count, string(buffer))
	return string(buffer)
}

func writeMessage(f *os.File, str string) {
	buffer := []byte(str)
	_, err := f.Write(buffer)
	if err != nil {
		fmt.Printf("Write error occured: %v\n", err)
	}
	//fmt.Printf("Sent %d bytes: %s\n", count, str)
}

func waitTimer(cycles int64, f *os.File) int {
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

func priority_test(f *os.File, delay int64) {
	rtt := delay
	str := readMessage(f)
	cycles, err := strconv.ParseInt(strings.TrimSpace(str), 0, 64)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("readMessage string: %s\n", str)
	fmt.Printf("rtt: %d\n", rtt)
	fmt.Printf("cycles: %d\n", cycles)
	waitTimer(cycles-rtt-rtt, f)
	//writeMessage(f, fmt.Sprintf("%27s", "COMPLETE"))
}

func calculateDelayJitter(f *os.File) (aDelay int64, aJitter int64) {
	samples := 100
	delayEntries := make([]int64, samples)
	// Delay/Jitter Calculations
	//totalDelay := uint64(0)
	for i := 0; i < samples; i++ {
		delayEntries[i] = exchangeTimestamps(f)
		log.Printf("%v\n", i)
		//totalDelay += delayEntries[i]
	}
	//avgdelay := totalDelay / uint64(samples)
	sort.Sort(int64arr(delayEntries))
	avgdelay := (delayEntries[49] + delayEntries[59])/2
	totaljitter := int64(0)
	jitter := int64(0)
	extremes := 0
	jitterEntries := make([]int64, samples)
	for i := 0; i < samples; i++ {
		if delayEntries[i] > 300000000 {
			extremes += 1
			continue
		}
		if delayEntries[i] > avgdelay {
			jitter = delayEntries[i] - avgdelay
		} else {
			jitter = avgdelay - delayEntries[i]
		}
		jitterEntries[i] = jitter;
		//fmt.Printf("Cur Jitter = %d \n", jitter)
		totaljitter += jitter
	}
	sort.Sort(int64arr(jitterEntries))
	samples -= extremes
	avgjitter := totaljitter / int64(samples)
	fmt.Printf("----------------------Delay Entries----------------------------------\n");
	for i := 0; i < samples; i++ {
		fmt.Printf("%v\n", delayEntries[i]);
	}
	fmt.Printf("----------------------Jitter Entries---------------------------------\n");
	for i := 0; i < samples; i++ {
		fmt.Printf("%v\n", jitterEntries[i]);
	}
	fmt.Printf("----------------------Final Statistics--------------------------------\n");
	fmt.Printf("median  delay  = %v \n", avgdelay);
	fmt.Printf("average jitter = %v \n", avgjitter);
	fmt.Printf("Num outliers = %v \n", extremes)
	return int64(avgdelay), int64(avgjitter)
}

func main() {
	f, err := openPort("/dev/ttyUSB0")
	if err != nil {
		fmt.Printf("Failed to open the serial port!")
	}
	//avgDelay, avgJitter := calculateDelayJitter(f)
	avgDelay, _ := calculateDelayJitter(f)
	//fmt.Printf("Avg. Delay = %v ns\n", avgDelay)
	//fmt.Printf("Avg. Jitter = %v ns\n", avgJitter)
	priority_test(f, avgDelay)
	f.Close()
}
