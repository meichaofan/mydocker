package container

import (
	"github.com/Sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
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
	}
	cmd.ExtraFiles = []*os.File{readPipe}

	//cmd.Dir = "/root/busybox"
	rootURL := "/root"
	mntURL := "/root/mnt"
	NewWorkSpace(rootURL, mntURL, volume)
	cmd.Dir = mntURL

	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootURL string, mntURL string,volume string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)

	if volume != "" {
		volumeUrls := volumeUrlExtract(volume)
		length := len(volume)
		if length == 2 || volumeUrls[0] != "" || volumeUrls[1] != "" {
			MountVolume(mntURL,volumeUrls)
			logrus.Infof("%q", volumeUrls)
		}else{
			logrus.Infof("Volume parameter input is not correct.")
		}
	}
}

func volumeUrlExtract(volume string) []string {
	var volumeUrls []string
	volumeUrls = strings.Split(volume,":")
	return volumeUrls
}

func MountVolume(mntURL string,volumeURLs []string) {
	// make suzhu zhuji folder
	parentUrl := volumeURLs[0]
	if err := os.Mkdir(parentUrl, 0777); err != nil {
		logrus.Infof("mkdir parent dir %s error: %v", parentUrl, err)
	}

	// in container create folder
	containerUrl := volumeURLs[1]
	containerVolumeUrl := mntURL + containerUrl
	if err := os.Mkdir(containerVolumeUrl, 0777); err != nil {
		logrus.Infof("mkdir parent dir %s error: %v", containerVolumeUrl, err)
	}

	dirs := "dirs=" + parentUrl
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", containerVolumeUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Infof("Mount volume failed. %v", err)
	}
}


// tar xvf busybox.tar -C ~/busybox as readOnlyLayer
// rootURL -> /root/
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "/busybox/"
	busyboxTarURL := rootURL + "/busybox.tar"
	exits, err := pathExists(busyboxURL)
	if err != nil {
		logrus.Infof("fail to judge whether dir %s exists: %v", rootURL, err)
	}
	if exits == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			logrus.Errorf("mkdir dir %s error: %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			logrus.Errorf("untar dir %s error: %v", busyboxTarURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer/"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error: %v", writeURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	// create mnt folder as mount point
	if err := os.Mkdir(mntURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error: %v", mntURL, err)
	}

	dirs := "dirs=" + rootURL + "/writeLayer:" + rootURL + "/busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
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

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer/"
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Errorf("delete %s error: $v", writeURL, err)
	}
}

func DeleteWorkSpace(rootURL string, mntURL string,volume string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(mntURL, volumeURLs)
		} else {
			DeleteMountPoint(mntURL)
		}
	} else {
		DeleteMountPoint(mntURL)
	}
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("umount %s error: %v", mntURL, err)
	}

	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("delete %s error: %v", mntURL, err)
	}
}

func DeleteMountPointWithVolume(mntURL string, volomeURLs []string) {
	containerUrl := mntURL + volomeURLs[1]
	cmd := exec.Command("umount",containerUrl)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err:= cmd.Run();err!=nil {
		logrus.Errorf("umount volume failed .%v",err)
	}
	DeleteMountPoint(mntURL)
}