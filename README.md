# docker-try

> 参考地址：https://github.com/xianlubird/mydocker

### code-3.1

创建并启动容器，namespace隔离，并且挂载入单独的proc

### code-3.2
在code-3.1基础上添加cgroup对容器资源进行限制。

将实现 ./x run -ti -m 100m -cpuset 1 -cpushare 512 /bin/sh 通过该方式来控制容器的内存和CPU配置

1. 设置cgroup的内存限制，将限制写入到cgroup对应目录的memory.limit_in_bytes文件中
2. 如何找到挂载了subsystem的hierarchy挂载目录呢？
   1. 通过 `/proc/self/mountinfo`可以找出与当前进程相关的mount信息
   2. Cgroups的hierarchy的虚拟文件系统是通过cgroup类型文件系统的mount挂载上去的，option 加上subsystem代表挂载的subsystem类型
   3. 比如memory，最后option 是rw，memory可以看出这一条挂载的subsystem是memory。那么在对应的/sys/fs/cgroup/memory中创建文件夹对应创建的cgroup就可以用来做内存限制。

> 疑问：cgroup是用来作为管理进程和subsystem的控制关系，subsystem是作用于cgroup节点，用来控制进程中资源的占用的。
>
> 那么hierarchy到底是干啥的？？？有什么具体的例子

```sequence
CgroupManager-->Subsystem实例: 2.创建Subsystem实例
Note Left of CgroupManager: 1.创建带资源限制的容器
CgroupManager->Subsystem实例: 3.在每个Subsystem对应的hierarchy上创建配置cgroup
Subsystem实例-->CgroupManager: 4.创建cgroup完成
CgroupManager->Subsystem实例: 5.将容器的进程移入每个Subsystem创建的cgroup中
Subsystem实例-->CgroupManager: 6.完成容器进程的资源限制
Note Left of CgroupManager: 7: 完成容器资源限制
```

#### cgroup三个概念

> hierarchy : 组织Cgroup，串联起来
>
> cgroup:管理进程和subsystem的控制关系
>
> subsystem:作用于cgroup节点，控制节点中进程的资源占用

1. cgroup hierarchy 中的节点，用于管理进程和subsystem的控制关系
2. subsystem作用于hierarchy上的cgroup节点，并控制节点中进程的资源占用
3. hierarchy将cgroup通过树状结构串起来，并通过虚拟文件系统的方式暴露给用户

