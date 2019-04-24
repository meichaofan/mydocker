package main

import (
	"mydocker/container"
	"github.com/Sirupsen/logrus"
	"os"
)

//启动init进程
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Run(); err != nil {
		logrus.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}
