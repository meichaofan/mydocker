package subsystems

import (
	"io/ioutil"
	"path"
	"fmt"
	"strconv"
	"os"
)

type CpusetSubSystem struct {
}

func (s *CpusetSubSystem) Name() string {
	return "cpuset"
}

func (s *CpusetSubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.CpuSet != "" {
			// 这里有个问题，还需要设置 cpuset.mems
			if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), 0644); err != nil {
				return fmt.Errorf("set %s cgroup fail %v", s.Name(), err)
			}
			// 设置 cpuset.mems
			if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "cpuset.mems"), []byte(strconv.Itoa(0)), 0644); err != nil {
				return fmt.Errorf("set cpuset.mems fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

func (s *CpusetSubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

func (s *CpusetSubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		return os.Remove(subsysCgroupPath)
	} else {
		return nil
	}
}
