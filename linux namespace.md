### linux namespace介绍

> docker是一个使用了linux namespace和cgroups的虚拟化工具
>
> 1. 什么是Linux namespace，它在docker内怎么被使用？

### 疑问知识点

> 1. pstree -pl 进程树
> 2. readlink /proc/2041/ns/uts `2041为进程号`
> 3. 什么是IPC？
> 4. /proc文件内容
> 5. mount挂载操作

### 概念

本质：隔离，隔离不同级别

| Namespace类型     | 系统调用参数  |        |
| ----------------- | ------------- | ------ |
| Mount Namespace   | CLONE_NEWNS   | 2.4.19 |
| UTS Namespace     | CLONE_NEWUTS  |        |
| IPC Namespace     | CLONE_NEWIPC  |        |
| PID Namespace     | CLONE_NEWPID  |        |
| Network Namespace | CLONE_NEWNET  |        |
| User Namespace    | CLONE_NEWUSER |        |

#### UTS Namespace

UTS Namespace:隔离nodename和domainname 两个系统标识，在UTS namespace中每个Namespace允许有自己的hostname

#### IPC Namespace

隔离System V IPC 和 POSIX message queues，每一个IPC Namespace都有自己的System V IPC 和 POSIX message queue

#### PID Namespace

> ps -ef 在容器内，前台运行的那个进程PID是1，但是在容器外，使用ps -ef会发现同样的进程拥有不同的PID

用来隔离进程ID的，同样一个进程在不同的PID namespace里可以拥有不同的PID

#### Mount Namespace

隔离各个进程看到的挂载点视图。在不同的Namespace的进程中，看到的文件系统层次是不一样的。在Mount Namespace中调用mount()和umount()仅仅只会影响当前Namespace内的文件系统，而对全局文件系统没有影响

#### User Namespace

User Namespace主要是隔离用户的用户组ID。一个进程的User ID和Group ID 在User Namespace内外可以是不同的

> 宿主机上以一个非root用户运行创建一个User Namespace，然后在User Namespace里面却映射成root用户

#### Network Namespace

- 隔离网络设备、Ip地址端口等网络栈的Namespace。可以让每个容器拥有自己独立的（虚拟的）网络设备，而且容器内的应用可以绑定到自己的端口，每个Namespace内的端口都不会互相冲突。
- 在宿主机搭建**网桥**后，就能方便容器之间通信，且不同容器上的应用可以使用相同的端口。

```go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	fmt.Println("main starting...")

	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	// cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(1), Gid: uint32(1)}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(-1)
}

```

### mount知识点

> 参考链接：https://www.sohu.com/a/260181668_467784
>
> mount namespace



