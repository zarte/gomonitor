package config

import (
	"os"
	"path/filepath"
	"sync"
	"github.com/zarte/zloger"
)

var Gconfig = new(Sysconfig)
type Sysconfig struct {
	ChanList *sync.Map
	CurExePath string
	WebPort string
	ExeList *sync.Map
	GLoger *zloger.Loger
	SecretKey string
	Expire int
}
func SetConfig() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Gconfig.CurExePath =  "./"
	}else{
		Gconfig.CurExePath =  dir+"/"
	}
	Gconfig.CurExePath =  "./"
	Gconfig.GLoger = zloger.NewLog(Gconfig.CurExePath +"logs")


	Gconfig.WebPort = "5556"
	Gconfig.ChanList = new (sync.Map)
	Gconfig.ExeList = new (sync.Map)
	Gconfig.SecretKey = "cptbtptp"
	Gconfig.Expire = 60*60*24
}