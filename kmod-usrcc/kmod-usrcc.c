/*  
 *  kmod-usrcc.c - Enables user-space access to ARM cycle counter.
 */
#include <linux/module.h>	/* Needed by all modules */
#include <linux/kernel.h>	/* Needed for KERN_INFO */
#include <linux/init.h>		/* Needed for the macros */
#define DRIVER_AUTHOR "Sean Hamilton <skhamilt@eng.ucsd.edu>"
#define DRIVER_DESC   "Enables user-space access to ARM cycle counter"
MODULE_AUTHOR(DRIVER_AUTHOR);
MODULE_DESCRIPTION(DRIVER_DESC);

/* 
 * Get rid of taint message by declaring code as GPL. 
 */
MODULE_LICENSE("GPL");


int init_module(void)
{
	printk(KERN_INFO "Loading kmod-usrcc\n");

	/* enable user-mode access to the performance counter */
	asm ("MCR p15, 0, %0, C9, C14, 0\n\t" :: "r"(1));

	/* disable counter overflow interrupts (just in case)*/
	asm ("MCR p15, 0, %0, C9, C14, 2\n\t" :: "r"(0x8000000f));


	/* 
	 * A non 0 return means init_module failed; module can't be loaded. 
	 */
	return 0;
}

void cleanup_module(void)
{
	printk(KERN_INFO "Unloading kmod-usrcc\n");
}
