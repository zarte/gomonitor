// +build windows

package task

import (
	"bufio"
	"context"
	"fmt"
	"gomonitor/config"
	"gomonitor/utils"
	"io"
	"os/exec"
)


func Command(cmd string, exeid string) error {
	c := exec.Command("cmd", "/C", cmd) 	// windows
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
			byte2String :=utils.ConvertByte2String(readString, "GB18030")
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
			byte2String := utils.ConvertByte2String(readString, "GB18030")
			fmt.Print(byte2String)
			config.Gconfig.GLoger.InfoLog(byte2String,exeid+"_")
		}
	}()
	err = c.Start()
	return err
}


