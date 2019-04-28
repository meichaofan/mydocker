package subsystems

import (
	"os"
	"bufio"
	"strings"
	"path"
	"fmt"
)

// 通过subsystem找到cgroup挂载点
// memory -> /sys/fs/cgroup/memory
func FindCgroupMountPoint(subsystem string) string {
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

/**
返回 例如 - /sys/fs/cgroup/docker
 */
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountPoint(subsystem)

	fmt.Printf("The %s cgroupRoot is %s", subsystem, cgroupRoot)
	fmt.Println()

	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0775); err != nil {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}

		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
