### 背景

1. Linux Namespace技术，帮助进程隔离出自己单独的空间
2. docker是怎么限制每个空间的大小，保证不会相互争抢资源呢？用到linux Cgroups技术

### 疑问点

> 1. hierarchy ? linux系统基础知识

### 简介

Linux Cgroups 提供了一组对进程及将来子进程的资源限制、控制和统计功能，这些资源包括CPU、内存、存储、网络等。通过Cgroups，可以方便地限制某个进程的资源占用，并且可以实时地监控进程的监控和统计信息



### 三大组件

#### cgroups

对进程分组的一种机制，一个cgroup包含一组进程，并可以在cgroup上增加Linux subsystem的各种参数配置，将一组进程和一组subsystem的系统参数关联起来

#### subsystem

一组资源控制的模块

#### hierarchy

把一组cgroup串成一个树状的结构。通过树状结构，就可以做到继承。



### Union File System

> UFS：为Linux、FreeBSD、NetBSD操作系统设计的，把其他文件系统联合到一个联合挂载点的文件系统服务
>
> 写时复制技术
>
> 参考链接：https://www.cnblogs.com/sparkdev/p/11237347.html

##### AUFS

> 可写分支的负载均衡
>
> AUFS 依然是docker支持的一种存储驱动类型。
>
> 下面介绍AUFS如何利用AUFS存储image和container

```go
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

//挂载memory subsystem的hierarchy的根目录位置
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

//通过go语言实现cgroup限制容器的资源
func main() {

	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	// syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	// defualtMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defualtMountFlags), "")

	fmt.Println("os.Args[0]", os.Args[0])

	if os.Args[0] == "/proc/self/exe" {
		//容器进程
		fmt.Printf("current pi %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	} else {
		//得到fork出来的进程映射在外部命名空间的pid
		fmt.Printf("%v", cmd.Process.Pid)

		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)

		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
	cmd.Process.Wait()
}

```







