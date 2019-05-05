/**
运行容器内init进程
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"math/rand"
	"mydocker/cgroups"
	"mydocker/cgroups/subsystems"
	"mydocker/container"
	"os"
	"strconv"
	"strings"
	"time"
)

//启动init进程
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, volume string, containerName string, imageName string) {
	parent, writePipe := container.NewParentProcess(tty, volume, containerName, imageName)
	if parent == nil {
		logrus.Error("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	//记录容器信息
	err := recordContainerInfo(parent.Process.Pid, comArray, containerName, volume)
	if err != nil {
		logrus.Errorf("Record container info error %v", err)
		return
	}

	// 资源限制
	// use mydocker-cgroup as cgroup name
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)

	if tty {
		parent.Wait()
		// container exit
		container.DeleteWorkSpace(volume, containerName)
		deleteContainerInfo(containerName)
	}
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func recordContainerInfo(containerPID int, commandArray []string, containerName string, volume string) error {
	//生成10位数的容器id
	id := randStringBytes(10)
	//以当前时间作为容器创建时间
	createTime := time.Now().Format("2016-01-02 15:04:05")
	command := strings.Join(commandArray, " ")
	//如果不指定容器名，那么就以容器id作为容器名
	if containerName == "" {
		containerName = id
	}
	//生成容器信息的结构体实例
	containerInfo := &container.ContainerInfo{
		Id:         id,
		Pid:        strconv.Itoa(containerPID),
		Command:    command,
		CreateTime: createTime,
		Status:     container.RUNNING,
		Name:       containerName,
		Volume:     volume,
	}

	//将容器对象json序列化成字符串
	jsonBytes, err := json.Marshal(containerInfo)
	if err != nil {
		logrus.Error("Record container info error: %v", err)
		return err
	}
	jsonStr := string(jsonBytes)

	//拼凑一下存储容器信息路径
	dirUrl := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	//如果该路径不存在，则创建之
	if err := os.MkdirAll(dirUrl, 0622); err != nil {
		logrus.Error("mkdir error %s error %v", dirUrl, err)
		return err
	}

	fileName := dirUrl + "/" + container.ConfigName

	//创建最终文件 config,json
	file, err := os.Create(fileName)
	defer file.Close()

	if err != nil {
		logrus.Error("Create file %s error %v", fileName, err)
		return err
	}
	//将json化之后的数据，写入到json文件中
	if _, err := file.WriteString(jsonStr); err != nil {
		logrus.Errorf("file write string error %v", err)
		return err
	}

	return err
}

//当以tty方式创建容器退出后，删除相应的记录文件
func deleteContainerInfo(containerName string) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		logrus.Errorf("Remove dir %s error %v", dirURL, err)
	}
}

//生成container名字
func randStringBytes(n int) string {
	letterBytes := "1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
