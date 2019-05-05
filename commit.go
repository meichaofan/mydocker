package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"mydocker/container"
	"os/exec"
)

func commitContainer(containerName, imageName string) {
	mntURL := fmt.Sprintf(container.MntUrl, containerName)
	mntURL += "/"
	imageTar := container.RootUrl + "/" + imageName + ".tar"
	fmt.Printf("%s", imageTar)
	if _, err := exec.Command("tar", "-cvf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		logrus.Errorf("tar folder %s error: %v \n", mntURL, err)
	}
}
