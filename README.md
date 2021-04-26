# docker-try

### code-3.1
创建并启动容器，namespace隔离，并且挂载入单独的proc

### code-3.2
在code-3.1基础上添加cgroup对容器资源进行限制。

将实现 ./x run -ti -m 100m -cpuset 1 -cpushare 512 /bin/sh 通过该方式来控制容器的内存和CPU配置

#### cgroup三个概念

1. cgroup hierarchy 中的节点，用于管理进程和subsystem的控制关系
2. subsystem作用于hierarchy上的cgroup节点，并控制节点中进程的资源占用
3. hierarchy将cgroup通过树状结构串起来，并通过虚拟文件系统的方式暴露给用户

