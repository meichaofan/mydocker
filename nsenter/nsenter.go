package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

//这里的__attribute__((constructor))的是，一旦这个包被引用，那么这个函数就会被自动执行
//类似于构造函数，会在程序一启动的时候运行
__attribute__((constructor)) void enter_namespace(void) {
	char *mydocker_pid;
	//从环境变量获取取药进去的PID
	mydocker_pid = getenv("mydocker_pid");
	if (mydocker_pid){
		//fprintf(stdout,"got mydocker_pid=%s\n",mydocker_pid);
	}else{
		//fprintf(stdout,"missing mydocker_pid env skip nsenter");
		return;
	}
	//从环境变量中获取需要执行的命令
	char *mydocker_cmd;
	mydocker_cmd = getenv("mydocker_cmd");
	if (mydocker_cmd){
		//fprintf(stdout,"got mydocker_cmd=%s\n",mydocker_cmd);
	}else{
		//fprintf(stdout,"missing mydocker_cmd env skip nsenter");
		return;
	}
	int i;
	char nspath[1024];
	//需要进入5中namespace
	char *namespace[] = {"ipc","uts","mnt","net","pid"};
	for(i=0;i<5;i++){
		sprintf(nspath,"/proc/%s/ns/%s",mydocker_pid,namespace[i]);
		int fd = open(nspath,O_RDONLY);
		//这里调用setns系统调用进入对应的namespace
		if(setns(fd,0)==-1){
			//fprintf(stderr, "setns on %s namespace failed: %s\n", namespaces[i], strerror(errno));
		}else{
			//fprintf(stdout, "setns on %s namespace successed\n",namespace[i]);
		}
		close(fd);
	}
	//在进入的Namespace中执行指定的命令
	int res = system(mydocker_cmd);
	//退出
	exit(0);
	return;
}
*/
import "C"
