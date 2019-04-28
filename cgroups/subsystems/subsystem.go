package subsystems

// 用于传递资源限制配置的结构体，包括内存限制，CPU时间片权重，CPU核心数
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// Subsystem接口，每个Subsystem可以实现下面4个接口
// 这里将cgroup抽象成了path，原因是cgroup在hierarchy的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	// 返回Subsystem的名字，比如cpu memory
	Name() string
	// 设置某个cgroup在这个Subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	// 将进程添加到某个cgroup中
	Apply(path string, pid int) error
	// 移除某个cgroup
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&MemorySubSystem{},
		&CpuSubSystem{},
		&CpusetSubSystem{},
	}
)
