// +build !windows

package utils

import (
	"os/exec"
	"syscall"
)
func KillCmdTask(c *exec.Cmd)  {
	syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
}