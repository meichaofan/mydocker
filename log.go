package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"os"
)

func logContainer(containerName string) {
	//找到对应文件夹的位置
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	logFileLocation := dirURL + container.ContainerLogFile
	//打开日志文件
	file, err := os.Open(logFileLocation)
	defer file.Close()

	if err != nil {
		logrus.Errorf("container log file %s open error %v\n", logFileLocation, err)
		return
	}

	//将文件内的内容都读取出来
	content, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("container log  file %s read error %v\n", logFileLocation, err)
		return
	}

	//使用fmt.fprint函数将读取出来的内容输出到标准输出
	fmt.Fprint(os.Stdout, string(content))
}
