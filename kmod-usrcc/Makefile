obj-m += kmod-usrcc.o

all: gotest ctest install

gotest: test.go 
	go build test.go
	mv test gotest

install: module 
	cp kmod-usrcc.ko /lib/modules/$(shell uname -r)/kernel/drivers/clk
	echo kmod-usrcc >> /etc/modules
	depmod
	modprobe kmod-usrcc

ctest: test.c
	gcc -o ctest test.c

module: kmod-usrcc.c
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(PWD) clean
	rm -f ctest gotest *.o
