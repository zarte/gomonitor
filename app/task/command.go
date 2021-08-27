package task

import (
	"context"
	"encoding/json"
	"errors"
	"gomonitor/config"
	"gomonitor/utils"
	"net/http"
	"strings"
)


func StartExe(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		Authorization :=r.Header.Get("Authorization")
		tmp := strings.Fields(Authorization)
		if len(tmp)>1 {
			_, err := utils.ParseToken(tmp[1])
			if err != nil {
				result, _ := json.Marshal(&utils.Comresult{
					Code: 4,
					Msg:  "fail",
				})
				w.Write(result)
				return
			}
		}else{
			result, _ := json.Marshal(&utils.Comresult{
				Code: 4,
				Msg:  "fail",
			})
			w.Write(result)
			return
		}
		query :=r.URL.Query()
		if len(query["exeid"]) >0 && len(query["op"])>0 {
			exeId :=query["exeid"][0]
			switch query["op"][0] {
			case "start":
				res,err:=startTask(exeId)
				if res{
					result, _ := json.Marshal(utils.Comresult{
						Code: 200,
						Msg:  "success",
					})
					w.Write(result)
				} else{
					result, _ := json.Marshal(utils.Comresult{
						Code: 4,
						Msg: err.Error(),
					})
					w.Write(result)
				}
				break
			case "stop":
				res,err:=stopTask(exeId)
				if res{
					result, _ := json.Marshal(utils.Comresult{
						Code: 200,
						Msg:  "success",
					})
					w.Write(result)
				} else{
					result, _ := json.Marshal(utils.Comresult{
						Code: 4,
						Msg: err.Error(),
					})
					w.Write(result)
				}
				break
			case "del":
				_,res:=config.Gconfig.ExeList.Load(exeId)
				if res{
					config.Gconfig.ExeList.Delete(exeId)
					result, _ := json.Marshal(utils.Comresult{
						Code: 200,
						Msg:  "success",
					})
					w.Write(result)
				}else{
					result, _ := json.Marshal(utils.Comresult{
						Code: 4,
						Msg:  "任务不存在",
					})
					w.Write(result)
				}
				break
			default:
				break
			}
		}
	}
}

func ChangeExeInfo(exeid string, cmd string,status string, cancel context.CancelFunc){
	info,res := config.Gconfig.ExeList.Load(exeid)
	if res {
		config.Gconfig.ExeList.Store(exeid,utils.ExeInfo{
			Exeid:   exeid,
			Cmd:   cmd,
			Name: info.(utils.ExeInfo).Name,
			Status: status,
			Cancel: cancel,
		})
	}
}
func startTask(exeid string)(bool,error)  {
	//判断是否已经运行
	exeinfo,checkFlag :=config.Gconfig.ExeList.Load(exeid)
	if !checkFlag {
		return false,errors.New("程序不存在")
	}
	if exeinfo.(utils.ExeInfo).Status == "start" ||exeinfo.(utils.ExeInfo).Status == "run" {
		return false,errors.New("已在运行")
	}
	Command(exeinfo.(utils.ExeInfo).Cmd,exeid)
	return true,nil
}
func stopTask(exeid string) (bool,error) {
	exeInfo,err := config.Gconfig.ExeList.Load(exeid)
	if err {
		if exeInfo.(utils.ExeInfo).Status == "start" || exeInfo.(utils.ExeInfo).Status == "run"{
			exeInfo.(utils.ExeInfo).Cancel.(context.CancelFunc)()
		} else {
			return false,errors.New("程序未运行")
		}
		return true,nil
	}else{
		return false,errors.New("获取句柄失败")
	}
}