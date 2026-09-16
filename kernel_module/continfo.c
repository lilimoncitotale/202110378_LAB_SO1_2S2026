//includes obligatorios
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/init.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/mm.h>
#include <linux/jiffies.h>
#include <linux/pid_namespace.h>
#include <linux/nsproxy.h>
#include <linux/string.h>
#include <linux/slab.h>


//protototipos de funciones
static bool es_contenedor(struct task_struct *task);


//funcion que escribe el contenido del archivo /proc/continfo
static int continfo_show(struct seq_file *m, void *v){
	//escribir variables antes del si_meminfo
	struct sysinfo info;
	struct task_struct *task;
	unsigned long total_mb, free_mb, used_mb,vsz_kb, rss_kb, total_ram_kb;
	unsigned long mem_percent, cpu_percent;
	unsigned long long  total_time, system_time;
	char container_id[13] = "N/A";
	
	


	
	

	////////////////////////////////
	si_meminfo(&info);

	


	//convertir a MB
	total_mb = (info.totalram * PAGE_SIZE) / (1024 * 1024);
	free_mb = (info.freeram * PAGE_SIZE) / (1024 * 1024);
	used_mb = total_mb - free_mb;
	total_ram_kb = (info.totalram * PAGE_SIZE) / 1024;

	seq_printf(m, "Memoria total: %lu MB\n", total_mb);
	seq_printf(m, "Memoria libre: %lu MB\n", free_mb);
	seq_printf(m, "Memoria usada: %lu MB\n", used_mb);
	seq_printf(m, "Hola desde /proc\n");
	seq_printf(m, "Carnet: 202110378\n");

	seq_printf(m,"------------------LISTADO DE PROCESOS ACTIVOS---------------------		");
	//PROCESOS//////////////////////////////////
	//proceso de contenedores
	seq_printf(m, "\n------------------PROCESOS DE CONTENEDORES----------------------\n");
	seq_printf(m, "\n=== PROCESOS DE CONTENEDORES ===\n");
	seq_printf(m, "%-8s %-20s %-12s %-12s %-8s %-8s\n",
           "PID", "NOMBRE", "VSZ(KB)", "RSS(KB)", "MEM%", "CPU%");
	seq_printf(m, "%-8s %-20s %-12s %-12s %-8s %-8s\n",
           "-----", "------------------", "---------", "---------", "------", "-----");
	for_each_process(task){
		if(es_contenedor(task)){
			//inicializar variables a 0 para evitar errores en procesos de contenedores
			vsz_kb = 0;
            rss_kb = 0;
            mem_percent = 0;
            cpu_percent = 0;
			strcpy(container_id, "N/A");

			//obtenemod el id usaremos comm
			strncpy(container_id, task->comm, 12);
			container_id[12] = '\0'; // Asegurarse de que la cadena esté terminada en nulo
			

			//calular metricas

			
			
			if(task->mm){
				
				vsz_kb = task->mm->total_vm * (PAGE_SIZE / 1024);
				rss_kb = get_mm_rss(task->mm) * (PAGE_SIZE / 1024);
				if(total_ram_kb >0){
					mem_percent = (rss_kb * 1000) / total_ram_kb;
				}
			}

			total_time = task->utime + task->stime;
			system_time = jiffies;
			if(system_time > 0){
				cpu_percent = (total_time * 100) / system_time;
			}

			seq_printf(m, "%-8d %-20s %-12lu %-12lu %-8lu %-8lu %-12s\n",task->pid, task->comm, vsz_kb, rss_kb, mem_percent, cpu_percent, container_id);
           
		}
	}
	//procesos del sistema
	seq_printf(m, "\n------------------PROCESOS DEL SISTEMA----------------------\n");
	seq_printf(m, "%-8s %-20s %-12s %-12s %-8s %-8s\n",
           "PID", "NOMBRE", "VSZ(KB)", "RSS(KB)", "MEM%", "CPU%");
	seq_printf(m, "%-8s %-20s %-12s %-12s %-8s %-8s\n",
           "-----", "------------------", "---------", "---------", "------", "-----");
	for_each_process(task){

		if (!es_contenedor(task)){//verificar que no sea null
			//inicializar variables a 0 para evitar errores en procesos del sistema
			vsz_kb = 0;
			rss_kb = 0;
			mem_percent = 0;
			cpu_percent = 0;

			if(task->mm){
				vsz_kb = task->mm->total_vm * (PAGE_SIZE / 1024);
				rss_kb = get_mm_rss(task->mm) * (PAGE_SIZE / 1024);
				if(total_ram_kb >0){
					mem_percent = (rss_kb * 1000) / total_ram_kb;
				}
			
			}
			total_time = task->utime + task->stime;
			system_time = jiffies;
			if(system_time > 0){
				cpu_percent = (total_time * 100) / system_time;
			}

			seq_printf(m, "%-8d %-20s %-12lu %-12lu %-8lu %-8lu\n",task->pid, task->comm, vsz_kb, rss_kb, mem_percent, cpu_percent);
			
		}
	}

	seq_printf(m, "fin de listado");

	return 0;
}

static bool es_contenedor(struct task_struct *task){
	//buscamos en comm primero
	if(strstr(task->comm, "docker") != NULL || strstr(task->comm, "containerd") != NULL || strstr(task->comm, "runc") != NULL){
		return true;

	}

	//tambien verificamos el PID namespace
	if(task->nsproxy && task->nsproxy->pid_ns_for_children != &init_pid_ns){
		return true;
	}
	return false;
}

//funcion que se ejecuta al abrir el archivo /proc/continfo
static int continfo_open(struct inode *inode, struct file *file){
	return single_open(file, continfo_show, NULL);
}

//estructura que define las operaciones del archivo /proc/continfo
static const struct proc_ops continfo_ops = {
	.proc_open = continfo_open,
	.proc_read = seq_read,
	.proc_lseek = seq_lseek,
	.proc_release = single_release,
};

//FUNCION DE INICIO
//se ejecuta cuando se carga el modulo (insmod) OBLIGATORIO
static int __init continfo_init(void){
    proc_create("continfo_pr2_so1_202110378", 0444, NULL, &continfo_ops);
	printk(KERN_INFO "Modulo continfo cargado\n");
    return 0;
}

//FUNCION DE SALIDA
//se ejecuta cuando se elimina el modulo (rmmod) OBLIGATORIO
static void __exit continfo_exit(void){
	remove_proc_entry("continfo_pr2_so1_202110378", NULL);
    printk(KERN_INFO "Adios mundo\n");
}

//registro de funciones OBLIGATORIO
module_init(continfo_init);
module_exit(continfo_exit);

//informacion del modulo OBLIGATORIO
MODULE_LICENSE("GPL");
MODULE_AUTHOR("limon");
MODULE_DESCRIPTION("modulo de telemetria para contenedores");
MODULE_VERSION("1.0");