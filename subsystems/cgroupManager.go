package subsystems

type CgroupManager struct {
	//cgroup 在hierarchy中的路径，相当于创建的cgroup目录相对于各root cgroup目录的路径
	Path     string
	Resource *ResourceConfig
}
