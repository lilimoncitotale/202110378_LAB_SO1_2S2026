#include <linux/module.h>
#define INCLUDE_VERMAGIC
#include <linux/build-salt.h>
#include <linux/elfnote-lto.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;
BUILD_LTO_INFO;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif

static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0x3d49336c, "module_layout" },
	{ 0x2cd6909b, "single_release" },
	{ 0x847b9c62, "seq_lseek" },
	{ 0xdae7f963, "seq_read" },
	{ 0x849c8d9e, "remove_proc_entry" },
	{ 0x92997ed8, "_printk" },
	{ 0x40240c8, "proc_create" },
	{ 0xd0da656b, "__stack_chk_fail" },
	{ 0x15ba50a6, "jiffies" },
	{ 0x9166fada, "strncpy" },
	{ 0x3dd38aca, "init_task" },
	{ 0x27402e0a, "seq_printf" },
	{ 0x40c7247c, "si_meminfo" },
	{ 0x9fef9b3d, "init_pid_ns" },
	{ 0x1e6d26a8, "strstr" },
	{ 0x5b8239ca, "__x86_return_thunk" },
	{ 0xf4f7e987, "single_open" },
	{ 0xbdfb6dbb, "__fentry__" },
};

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "7D2092C070C7B00ECEB51B1");
