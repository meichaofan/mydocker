package subsystems

import (
	"io/ioutil"
	"path"
	"fmt"
	"os"
	"strconv"
)

// memory subsystem 的实现
type MemorySubSystem struct {
}

// 返回 cgroup 的名字
func (s *MemorySubSystem) Name() string {
	return "memory"
}

//设置 cgroupPath 对应的 cgroup 的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	/**
	GetCgroupPath 是获取当前的 subsystem 在虚拟文件系统中的路径
	 */
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			/**
			设置这个 cgroup 的内存限制，即将限制写入到cgroup对应目录的memory.limit_in_bytes文件中
			 */
			if err = ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup %s fail %v", s.Name(), err)
			}
		}
		return nil
	} else {
		return err
	}
}

// 删除 cgroupPath 对应的 cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgoupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		// 删除 cgroup 便是删除对应的 cgroupPath 的目录
		fmt.Printf("the memory cgroup path is %s", subsysCgoupPath)
		return os.Remove(subsysCgoupPath)
	} else {
		return nil
	}
}

// 将一个进程加入到 cgroupPath 对应的 cgroup 中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	if subsysCgoupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		if err = ioutil.WriteFile(path.Join(subsysCgoupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}
