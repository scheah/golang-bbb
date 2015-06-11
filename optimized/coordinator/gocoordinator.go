package main

//#cgo LDFLAGS: -L. -lftd2xx -Wl,-rpath /usr/local/lib
//#include <stdint.h>
//#include <stdio.h>
//#include <stdlib.h>
//#include <sys/time.h>
//#include "ftd2xx.h"
//
//FT_STATUS	ftStatus;
//FT_HANDLE	ftHandle0;
//
//char * readMsg() {
//	DWORD RxBytes = 27;
//	DWORD BytesReceived;
//	char * msg = malloc(RxBytes);
//	ftStatus = FT_Read(ftHandle0,msg,RxBytes,&BytesReceived);	
//	if (ftStatus == FT_OK) {
//		if (BytesReceived == RxBytes) {
//			printf("reading: %s\n", msg);
//			printf("FT_Read OK\n", msg);
//			return msg;
//		}
//		else {
//			printf("FT_Read Timeout\n");
//		}
//	}
//	else {
//		printf("FT_Read Failed\n");
//	}
// 	
//}
//
//void writeMsg(char * msg) {
//	DWORD BytesWritten;
//	//char TxBuffer[64]; // Contains data to write to device
//	//strncpy(TxBuffer, msg, sizeof(TxBuffer));
// 	//TxBuffer[sizeof(TxBuffer) - 1] = '\0';
//	printf("writing: %s\n", msg);
//	ftStatus = FT_Write(ftHandle0, msg, 27, &BytesWritten);
//	if (ftStatus == FT_OK) {
//		printf("FT_Write OK\n");
//	}
//	else {
//		printf("FT_Write Failed\n");
//	}
//}
//
//
//int setup()
//{
//	int iport;
//	static FT_PROGRAM_DATA Data;
//	static FT_DEVICE ftDevice;
//	DWORD libraryVersion = 0;
//	int retCode = 0;
//
//	ftStatus = FT_GetLibraryVersion(&libraryVersion);
//	if (ftStatus == FT_OK)
//	{
//		printf("Library version = 0x%x\n", (unsigned int)libraryVersion);
//	}
//	else
//	{
//		printf("Error reading library version.\n");
//		return 1;
//	}
//
//	iport = 0;
//
//	printf("Opening port %d\n", iport);
//	
//	ftStatus = FT_Open(iport, &ftHandle0);
//	if(ftStatus != FT_OK) {
//		/* 
//			This can fail if the ftdi_sio driver is loaded
//		 	use lsmod to check this and rmmod ftdi_sio to remove
//			also rmmod usbserial
//		 */
//		printf("FT_Open(%d) failed\n", iport);
//		return 1;
//	}
//
//	printf("FT_Open succeeded.  Handle is %p\n", ftHandle0);
//
//  ftStatus = FT_ResetDevice(ftHandle0);
//  if (ftStatus == FT_OK)
//  	printf("FT_ResetDevice OK\n");
//  else
//      printf("FT_ResetDevice FAILED\n");
//	ftStatus = FT_GetDeviceInfo(ftHandle0,
//	                            &ftDevice,
//	                            NULL,
//	                            NULL,
//	                            NULL,
//	                            NULL); 
//	if (ftStatus != FT_OK) 
//	{ 
//		printf("FT_GetDeviceType FAILED!\n");
//		retCode = 1;
//		goto exit;
//	}  
//
//	printf("FT_GetDeviceInfo succeeded.  Device is type %d.\n", 
//	       (int)ftDevice);
//
//	/* MUST set Signature1 and 2 before calling FT_EE_Read */
//	Data.Signature1 = 0x00000000;
//	Data.Signature2 = 0xffffffff;
//	Data.Manufacturer = (char *)malloc(256); /* E.g "FTDI" */
//	Data.ManufacturerId = (char *)malloc(256); /* E.g. "FT" */
//	Data.Description = (char *)malloc(256); /* E.g. "USB HS Serial Converter" */
//	Data.SerialNumber = (char *)malloc(256); /* E.g. "FT000001" if fixed, or NULL */
//	if (Data.Manufacturer == NULL ||
//	    Data.ManufacturerId == NULL ||
//	    Data.Description == NULL ||
//	    Data.SerialNumber == NULL)
//	{
//		printf("Failed to allocate memory.\n");
//		retCode = 1;
//		goto exit;
//	}
//
//	ftStatus = FT_EE_Read(ftHandle0, &Data);
//	if(ftStatus != FT_OK) {
//		printf("FT_EE_Read failed\n");
//		retCode = 1;
//		goto exit;
//	}
//
//	printf("FT_EE_Read succeeded.\n\n");
//	
//	printf("Signature1 = %d\n", (int)Data.Signature1);			
//	printf("Signature2 = %d\n", (int)Data.Signature2);			
//	printf("Version = %d\n", (int)Data.Version);				
//								
//	printf("VendorId = 0x%04X\n", Data.VendorId);				
//	printf("ProductId = 0x%04X\n", Data.ProductId);
//	printf("Manufacturer = %s\n", Data.Manufacturer);			
//	printf("ManufacturerId = %s\n", Data.ManufacturerId);		
//	printf("Description = %s\n", Data.Description);			
//	printf("SerialNumber = %s\n", Data.SerialNumber);			
//	printf("MaxPower = %d\n", Data.MaxPower);				
//	printf("PnP = %d\n", Data.PnP) ;					
//	printf("SelfPowered = %d\n", Data.SelfPowered);			
//	printf("RemoteWakeup = %d\n", Data.RemoteWakeup);			
//	if (ftDevice == FT_DEVICE_232R)
//	{
//		/* Rev 6 (FT232R) extensions */
//		printf("232R:\n");
//		printf("-----\n");
//		printf("\tUseExtOsc = 0x%X\n", Data.UseExtOsc);			// Use External Oscillator
//		printf("\tHighDriveIOs = 0x%X\n", Data.HighDriveIOs);			// High Drive I/Os
//		printf("\tEndpointSize = 0x%X\n", Data.EndpointSize);			// Endpoint size
//
//		printf("\tPullDownEnableR = 0x%X\n", Data.PullDownEnableR);		// non-zero if pull down enabled
//		printf("\tSerNumEnableR = 0x%X\n", Data.SerNumEnableR);		// non-zero if serial number to be used
//
//		printf("\tInvertTXD = 0x%X\n", Data.InvertTXD);			// non-zero if invert TXD
//		printf("\tInvertRXD = 0x%X\n", Data.InvertRXD);			// non-zero if invert RXD
//		printf("\tInvertRTS = 0x%X\n", Data.InvertRTS);			// non-zero if invert RTS
//		printf("\tInvertCTS = 0x%X\n", Data.InvertCTS);			// non-zero if invert CTS
//		printf("\tInvertDTR = 0x%X\n", Data.InvertDTR);			// non-zero if invert DTR
//		printf("\tInvertDSR = 0x%X\n", Data.InvertDSR);			// non-zero if invert DSR
//		printf("\tInvertDCD = 0x%X\n", Data.InvertDCD);			// non-zero if invert DCD
//		printf("\tInvertRI = 0x%X\n", Data.InvertRI);				// non-zero if invert RI
//
//		printf("\tCbus0 = 0x%X\n", Data.Cbus0);				// Cbus Mux control
//		printf("\tCbus1 = 0x%X\n", Data.Cbus1);				// Cbus Mux control
//		printf("\tCbus2 = 0x%X\n", Data.Cbus2);				// Cbus Mux control
//		printf("\tCbus3 = 0x%X\n", Data.Cbus3);				// Cbus Mux control
//		printf("\tCbus4 = 0x%X\n", Data.Cbus4);				// Cbus Mux control
//
//		printf("\tRIsD2XX = 0x%X\n", Data.RIsD2XX); // non-zero if using D2XX
//	}
//	ftStatus = FT_SetBaudRate(ftHandle0, 115200); // Set baud rate to 115200
//	if (ftStatus == FT_OK) {
//		printf("FT_SetBaudRate OK\n");
//	}
//	else {
//		printf("FT_SetBaudRate Failed\n");
//	}
//	ftStatus = FT_SetFlowControl(ftHandle0, FT_FLOW_RTS_CTS, 0x11, 0x13);
//	if (ftStatus == FT_OK) {
//		printf("FT_SetFlowControl OK\n");
//	}
//	else {
//		printf("FT_SetFlowControl Failed\n");
//	}
//	UCHAR LatencyTimer = 1;
//	ftStatus = FT_SetLatencyTimer(ftHandle0, LatencyTimer );
//	if (ftStatus == FT_OK) {
//		printf("Set LatencyTimer: %u\n", LatencyTimer );
//	}
//	else {
//		printf("FT_SetLatencyTimer failed\n");
//		retCode = 1;
//		goto exit;
//	}
//	DWORD TransferSize = 64;
//	ftStatus = FT_SetUSBParameters(ftHandle0, TransferSize, TransferSize);
//	if (ftStatus == FT_OK) {
//		printf("In/Out transfer size set to 64 bytes\n");
//	}
//	else {
//		printf("FT_SetUSBParameters failed\n");
//		retCode = 1;
//		goto exit;
//	}
//
//
//exit:
//	free(Data.Manufacturer);
//	free(Data.ManufacturerId);
//	free(Data.Description);
//	free(Data.SerialNumber);
//	printf("Returning %d\n", retCode);
//	return retCode;
//}
//
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
import "sort"
import "log"

var timeStampLayout = "02:Jan:2006:15:04:05.000000"

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

func calculateDelayRTT(p0 string, p1 string, p2 string, p3 string) int64 {
	// parse time stamp string
	t0 := parseTimestamp(p0)
	t1 := parseTimestamp(p1)
	t2 := parseTimestamp(p2)
	t3 := parseTimestamp(p3)
	delayRTT := (t3.Sub(t0) + t2.Sub(t1))
	//fmt.Printf("RTT delay: %v\n", delayRTT)
	return int64(delayRTT.Nanoseconds())
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

func exchangeTimestamps() int64 {
	t0 := generateTimestamp()
	C.writeMsg(C.CString(t0))
	t2 := C.GoString(C.readMsg())
	t3 := generateTimestamp()
	C.writeMsg(C.CString(t3))
	// gotta get t1...
	t1 := C.GoString(C.readMsg())
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

func upperWindow(c1 chan string, cycles int64) {
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

func lowerWindow(c1 chan string, cycles int64) {
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

func priority_test(usb *os.File, windowSize int64, waitTime int64, led_blue *os.File, led_green *os.File) {
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

func calculateDelayJitter() (aDelay int64, aJitter int64) {
	samples := 100
	delayEntries := make([]int64, samples)
	// Delay/Jitter Calculations
	//totalDelay := uint64(0)
	for i := 0; i < samples; i++ {
		delayEntries[i] = exchangeTimestamps()
		log.Printf("%v\n", i)
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
	/*usb, err := openPort("/dev/ttyUSB0")
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
	//avgDelay, avgJitter := calculateDelayJitter(usb)
	//calculateDelayJitter(usb)
	//fmt.Printf("Avg. Delay = %v ns\n", avgDelay)
	//fmt.Printf("Avg. Jitter = %v ns\n", avgJitter)
	//windowSize := int64(600000000) //cycles at 1ghz is 0.6s
	//waitTime := int64(500000000)   //cycles at 1ghz is 0.5s
	//time.Sleep(1*time.Second)
	//priority_test(usb, windowSize, waitTime, led_blue, led_green)
	//usb.Close()
	led_blue.Close()
	led_green.Close()
		*/
	C.setup()
	calculateDelayJitter()
}
