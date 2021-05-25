package subsystems

import "github.com/sirupsen/logrus"

type CgroupManager struct {
	//cgroup 在hierarchy中的路径，相当于创建的cgroup目录相对于各root cgroup目录的路径
	Path     string
	Resource *ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager  {
	return &CgroupManager{
		Path: path,
	}
}

//Apply 将进程PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error  {
	for _, subSysIns := range(SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

//Set 设置各个subsystem挂载中的cgroup的资源限制
func (c *CgroupManager)Set(res *ResourceConfig) error  {
	for _, subSysIns := range (SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// 释放各个subsystem挂载中的cgroup
func (c *CgroupManager) Destory() error  {
	for _, subSysIns := range SubsystemsIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}


