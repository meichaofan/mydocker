package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"io/ioutil"
	"mydocker/container"
	"os"
	"text/tabwriter"
)

func ListContainer() {
	//找到容器存储信息路径/var/run/mydocker
	dirURL := fmt.Sprintf(container.DefaultInfoLocation, "")
	dirURL = dirURL[:len(dirURL)-1]
	//读取该文件夹下所有文件
	files, err := ioutil.ReadDir(dirURL)
	if err != nil {
		logrus.Errorf("Read dir %s error %v", dirURL, err)
		return
	}

	var containers []*container.ContainerInfo

	//遍历 /var/run/mydocker 下所有文件
	for _, file := range files {
		tmpContainer, err := getContainerInfo(file);
		if err != nil {
			logrus.Errorf("get container info error %v", err)
			continue
		}
		containers = append(containers, tmpContainer)
	}

	// text/tabwriter
	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	fmt.Fprint(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	for _, item := range containers {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.Id,
			item.Name,
			item.Pid,
			item.Status,
			item.Command,
			item.CreateTime,
		)
	}

	//刷出输出缓冲区内容到标准输出
	if err := w.Flush(); err != nil {
		logrus.Errorf("Flush error %v", err)
		return
	}
}

func getContainerInfo(file os.FileInfo) (*container.ContainerInfo, error) {
	//获取文件名
	containerName := file.Name()
	//根据文件生成绝对路径
	configFileDir := fmt.Sprintf(container.DefaultInfoLocation, containerName)
	configFileDir = configFileDir + container.ConfigName
	//读取config.json文件内容
	content, err := ioutil.ReadFile(configFileDir)
	if err != nil {
		logrus.Errorf("read config %s file error %v", configFileDir, err)
		return nil, err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(content, &containerInfo); err != nil {
		logrus.Errorf("json unmarshal error %v", err)
		return nil, err
	}
	return &containerInfo, nil
}
