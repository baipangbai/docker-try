package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command %s", command)

	//声明你要这个新的mount namespace独立。
	// syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV

	//为什么要在容器中挂载/proc呢， 主要原因是因为ps、top等命令依赖于/proc目录。当隔离PID的时候，ps、top等命令还是未隔离的时候一样输出。 为了让隔离空间ps、top等命令只输出当前隔离空间的进程信息。需要单独挂载/proc目录。
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("err.Error()")
	}
	return nil
}
