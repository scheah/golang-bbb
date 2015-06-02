#!/bin/bash
export PATH=/usr/local/go/bin:$PATH
insmod kmod-usrcc/kmod-usrcc.ko
cd /sys/devices/system/cpu/cpu0/cpufreq
cpufreq-set -g userspace
echo 1000000 > scaling_setspeed
cd ~
echo 60 > /sys/class/gpio/export
echo 48 > /sys/class/gpio/export
echo high > /sys/class/gpio/gpio60/direction
echo high > /sys/class/gpio/gpio48/direction

