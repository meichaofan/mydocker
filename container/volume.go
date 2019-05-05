package container

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName, imageName)

	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)
		length := len(volume)
		if length == 2 || volumeUrls[0] != "" || volumeUrls[1] != "" {
			MountVolume(volumeUrls, containerName)
			logrus.Infof("%q", volumeUrls)
		} else {
			logrus.Infof("Volume parameter input is not correct.")
		}
	}
}

func volumeUrlExtract(volume string) []string {
	var volumeUrls []string
	volumeUrls = strings.Split(volume, ":")
	return volumeUrls
}

func MountVolume(volumeURLs []string, containerName string) error {
	// 创建宿主主机文件目录
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		logrus.Infof("mkdir parent dir %s error: %v", parentUrl, err)
	}

	// in container create folder
	containerUrl := volumeURLs[1]
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	containerVolumeUrl := mntUrl + "/" + containerUrl

	if err := os.Mkdir(containerVolumeUrl, 0777); err != nil {
		logrus.Infof("mkdir parent dir %s error: %v", containerVolumeUrl, err)
	}

	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeUrl).CombinedOutput()
	if err != nil {
		logrus.Errorf("Mount volume failed. %v", err)
		return err
	}
	return nil
}

// tar xvf busybox.tar -C ~/busybox as readOnlyLayer
// rootURL -> /root/
func CreateReadOnlyLayer(imageName string) {
	unTarFolderUrl := RootUrl + "/" + imageName + "/"
	imageUrl := RootUrl + "/" + imageName + ".tar"
	exits, err := pathExists(unTarFolderUrl)
	if err != nil {
		logrus.Infof("fail to judge whether dir %s exists: %v", unTarFolderUrl, err)
	}
	if exits == false {
		if err := os.Mkdir(unTarFolderUrl, 0622); err != nil {
			logrus.Errorf("mkdir dir %s error: %v", unTarFolderUrl, err)
		}
		if _, err := exec.Command("tar", "-xvf", imageUrl, "-C", unTarFolderUrl).CombinedOutput(); err != nil {
			logrus.Errorf("untar dir %s error: %v", unTarFolderUrl, err)
		}
	}
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error: %v", writeURL, err)
	}
}

func CreateMountPoint(containerName, imageName string) error {
	// create mnt folder as mount point
	mntUrl := fmt.Sprintf(MntUrl, containerName)
	if err := os.MkdirAll(mntUrl, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error: %v", mntUrl, err)
	}

	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl, containerName)
	tmpImageLocation := RootUrl + "/" + imageName

	dirs := "dirs=" + tmpWriteLayer + ":" + tmpImageLocation
	_, err := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntUrl).CombinedOutput()
	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	return nil
}

// pan duan wenjian lujing shi fou cun zai
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl, containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Errorf("remove writeLayer %s error: $v", writeURL, err)
	}
}

func DeleteWorkSpace(volume, containerName string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(volumeURLs, containerName)
		} else {
			DeleteMountPoint(containerName)
		}
	} else {
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	_, err := exec.Command("umount", mntURL).CombinedOutput()
	if err != nil {
		logrus.Errorf("umount %s error: %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("delete %s error: %v", mntURL, err)
	}
	return nil
}

func DeleteMountPointWithVolume(volomeURLs []string, containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerUrl := mntURL + "/" + volomeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		logrus.Errorf("umount volume failed .%v", err)
		return err
	}

	if _, err := exec.Command("umount", mntURL).CombinedOutput(); err != nil {
		logrus.Errorf("umount mountpoint %s failed: %v", mntURL, err)
		return err
	}

	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("remove mountpoint dir %s error: %v", mntURL, err)
	}
	return nil
}
