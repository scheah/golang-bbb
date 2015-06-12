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

import (
	"C"
	"fmt"
	"time"
)

func main() {
	// init counters:
	C.init_perfcounters(1, 0)

	// measure the counting overhead:
	overhead := C.get_cyclecount()
	overhead = C.get_cyclecount() - overhead

	sleeptime := time.Duration(1) * time.Second
	t := C.get_cyclecount()
	time.Sleep(sleeptime) //delay
	//fmt.Println("test data to benchmark Println function call")

	t = C.get_cyclecount() - t

	fmt.Printf("Sleeping took exactly %d cycles and overhead of get_cyclecount() was %d cycles\n", t-overhead, overhead)
}
