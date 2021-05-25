package subsystems

type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

//Subsystem 这里将cgroup抽象成了path，原因是cgroup在hierarchy的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	//返回subsystem的名字，比如cpu memory
	Name() string
	//设置某个cgroup在这个subsystem中的资源限制
	Set(path string, res *ResourceConfig) error
	//将进程添加到某个cgroup中
	Apply(path string, pid int) error
	//移除某个cgroup
	Remove(path string) error
}

var (
	SubsystemsIns = []Subsystem{
		&MemorySubSystem{},
		&CpusetSubSystem{},
		&CpuSubSystem{},
	}
)
