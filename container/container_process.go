package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/mydocker/container"
	"os"
	"os/exec"
	"syscall"
)

//容器运行中所需的常量定义
var (
	RUNNING             string = "running"
	STOP                string = "stop"
	Exit                string = "exit"
	DefaultInfoLocation string = "/var/run/mydocker/%s/"
	ConfigName          string = "config.json"
	ContainerLogFile    string = "container.log"
	RootUrl             string = "/root"
	MntUrl              string = "/root/mnt/%s"
	WriteLayerUrl       string = "/root/writeLayer/%s"
)

//记录容器运行的状态
type ContainerInfo struct {
	Pid        string `json:"pid"`        //容器的init进程在宿主主机的PID
	Id         string `json:"id"`         //容器ID
	Name       string `json:"name"`       //容器名
	Command    string `json:"command"`    //容器内init运行的命令
	CreateTime string `json:"createTime"` //容器创建的时间
	Status     string `json:"status"`     //容器状态
	Volume     string `json:"volume"`     //目录映射
}

func NewParentProcess(tty bool, volume string, containerName string, imageName string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Error("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdout
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		//生成容器对应的目录 container.log文件
		dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
		if err := os.MkdirAll(dirURL, 0622); err != nil {
			logrus.Errorf("NewParentProcess mkdir %s error %v", dirURL, err)
			return nil, nil
		}
		stdLogFilePath := dirURL + ContainerLogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			logrus.Errorf("NewParentProcess create file %s error %v", dirURL, err)
			return nil, nil
		}
		//把生成好的文件赋值给stdout，这样就能把容器的标准输出重定向到这个文件中
		cmd.Stdout = stdLogFile
	}

	cmd.ExtraFiles = []*os.File{readPipe}

	NewWorkSpace(volume, imageName, containerName)
	cmd.Dir = fmt.Sprintf(MntUrl, containerName)

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
