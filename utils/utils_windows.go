// +build windows

package utils

import (
	"fmt"
	"os/exec"
	"strconv"
)

func KillCmdTask(c *exec.Cmd)  {
	kill := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(c.Process.Pid))
	err := kill.Run()
	if err != nil {
		fmt.Println(err)
	}
}