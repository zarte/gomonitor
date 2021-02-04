package config

import (
	"os"
	"path/filepath"
	"sync"
)

var Gconfig = new(Sysconfig)
type Sysconfig struct {
	ChanList *sync.Map
	CurExePath string
	WebPort string
	ExeList *sync.Map
}
func SetConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Gconfig.CurExePath =  "./"
	}else{
		Gconfig.CurExePath =  dir+"/"
	}
	Gconfig.CurExePath =  "./"
	Gconfig.WebPort = "5556"
	Gconfig.ChanList = new (sync.Map)
	Gconfig.ExeList = new (sync.Map)
}