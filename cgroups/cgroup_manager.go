package cgroups

import (
	"mydocker/cgroups/subsystems"
	"github.com/Sirupsen/logrus"
)

type CgroupManager struct {
	// cgroup 在 hierarchy 中的路径，相当于创建的cgroup目录相对于各 root cgroup 目录的路径
	path string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{path: path}
}

// 将进程的PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Apply(c.path, pid)
	}
	return nil
}

// 设置各个 subsystem 挂载中的 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		subSysIns.Set(c.path, res)
	}
	return nil
}

// 释放各个 subsystem 挂载中的 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range subsystems.SubsystemsIns {
		if err := subSysIns.Remove(c.path); err != nil {
			logrus.Warnf("remove cgroup %s fail %v", subSysIns.Name(), err)
		}
	}
	return nil
}
