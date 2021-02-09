// +build !windows

package task

import (
	"bufio"
	"context"
	"fmt"
	"gomonitor/utils"
	"gomonitor/config"
	"io"
	"os/exec"
	"syscall"
)

func Command(cmd string, exeid string) error {

	c := exec.Command("bash", "-c", cmd)  // mac or linux
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ChangeExeInfo(exeid,cmd,"run",cancel)

	// 监听关闭命令
	go func(c *exec.Cmd) {
		_:<-ctx.Done()
		utils.KillCmdTask(c)
	}(c)

	go func() {
		reader := bufio.NewReader(stdout)
		for {
			readString,_, err := reader.ReadLine()
			if err != nil || err == io.EOF {
				ChangeExeInfo(exeid,cmd,"stop",nil)
				fmt.Println("守护进程结束")
				return
			}
			byte2String :=utils.ConvertByte2String(readString, "UTF8")
			fmt.Print(byte2String)
			config.Gconfig.GLoger.InfoLog(byte2String,exeid+"_")
		}
	}()
	go func() {
		reader := bufio.NewReader(stderr)
		for {
			readString,_, err := reader.ReadLine()
			if err != nil || err == io.EOF {
				return
			}
			byte2String := utils.ConvertByte2String(readString, "UTF8")
			fmt.Print(byte2String)
			config.Gconfig.GLoger.InfoLog(byte2String,exeid+"_")
		}
	}()
	err = c.Start()

	return err
}
