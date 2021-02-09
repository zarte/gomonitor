package web

import (
	"encoding/json"
	"fmt"
	"gomonitor/app/task"
	"gomonitor/config"
	"gomonitor/utils"
	"io/ioutil"
	"net/http"
	"html/template"
	"strings"
	"time"
)

func Starweb()  {
	config.SetConfig()
	//读取list
	initExeList()
	http.HandleFunc("/exelist", ExeList)
	http.HandleFunc("/addexe", AddExeList)
	http.HandleFunc("/startexe", task.StartExe)
	http.HandleFunc("/seelog", ExeLog)

	http.HandleFunc("/checklogin", CheckLogin)

	http.HandleFunc("/", WebList)

	if err := http.ListenAndServe(":"+config.Gconfig.WebPort, nil); err != nil {
		fmt.Println(err)
	}

}
func initExeList()  {
	content,err := ioutil.ReadFile(config.Gconfig.CurExePath + "/tmpexelist.txt")
	if err !=nil{
		fmt.Println(err)
	}else{
		if string(content) != "" {
			var exelist []utils.ExeInfo
			err =json.Unmarshal(content,&exelist)
			if err !=nil{
				fmt.Println(err)
			}else{
				//循环添加
				for _,item :=range exelist {
					config.Gconfig.ExeList.Store(item.Exeid,utils.ExeInfo{
						Exeid: item.Exeid,
						Cmd: item.Cmd,
						Status : "stop",
					})
				}
			}
		}
	}
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "POST" {
		user :=r.FormValue("name")
		pwd :=r.FormValue("passwd")
		if user=="admin" && pwd == "123456" {
			token := utils.GetJwt(user,&utils.UserInfo{
				Id:     "1",
				Name:   "admin",
				Avatar: "",
				Ip:     "",
			})
			result,_ := json.Marshal(&utils.Comresult{
				Code: 200,
				Msg: "success",
				Data: token,
			})
			w.Write(result)
		}else{
			result,_ := json.Marshal(&utils.Comresult{
				Code: 4,
				Msg: "登录失败",
			})
			w.Write(result)
		}
	}
}
func WebList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" { //如果请求方法为get显示login.html,并相应给前端
		t, _ := template.ParseFiles(config.Gconfig.CurExePath +"list.html")
		t.Execute(w, nil)
	}
}

func GetExeList() []utils.ExeInfo{
	var list []utils.ExeInfo
	config.Gconfig.ExeList.Range(func(k, v interface{}) bool {
		list = append(list,v.(utils.ExeInfo) )
		return true
	})
	return list
}

func ExeLog(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		Authorization :=r.Header.Get("Authorization")
		tmp := strings.Fields(Authorization)
		if len(tmp)>1 {
			//验证
			_,err :=utils.ParseToken(tmp[1])
			if err!=nil{
				result, _ := json.Marshal(&utils.Comresult{
					Code: 4,
					Msg: "success",
				})
				w.Write(result)
				return
			}
			exeid :=r.FormValue("exeid")
			if exeid !=""{
				//获取log
				result, _ := json.Marshal(utils.Comresult{
					Code: 200,
					Msg:  "success",
					Data:getlogtail(exeid),
				})
				w.Write(result)
			} else {
				w.Write([]byte("no"))
			}
		}
	}
}
func ExeList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		Authorization :=r.Header.Get("Authorization")
		tmp := strings.Fields(Authorization)
		if len(tmp)>1 {
			//验证
			_,err :=utils.ParseToken(tmp[1])
			if err!=nil{
				result, _ := json.Marshal(&utils.Comresult{
					Code: 4,
					Msg: "success",
				})
				w.Write(result)
				return
			}
			result, err := json.Marshal(&utils.Comresult{
				Code: 200,
				Msg: "success",
				Data: GetExeList(),
			})
			w.Write(result)
		}
	}
}
func getlogtail(exeid string) string {
	filename := "info_"+exeid +"_" + time.Now().Format("2006-1-2")+".log"
	bytes, err := ioutil.ReadFile(config.Gconfig.CurExePath+"/logs/"+filename)
	if err != nil {
		fmt.Println("error : %s", err)
		return ""
	}
	return string(bytes)

}
func AddExeList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "POST" {
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
		exeid :=r.FormValue("exeid")
		cmd :=r.FormValue("cmd")
		if exeid !="" && cmd!="" {
			_,res:=config.Gconfig.ExeList.Load(exeid)
			if res{
				result, _ := json.Marshal(utils.Comresult{
					Code: 4,
					Msg:  "已存在",
				})
				w.Write(result)
			}else{
				config.Gconfig.ExeList.Store(exeid,utils.ExeInfo{
					Exeid: exeid,
					Cmd: cmd,
					Status : "stop",
				})
				//记录list
				go func() {
					result, _ := json.Marshal(GetExeList())
					ioutil.WriteFile(config.Gconfig.CurExePath + "/tmpexelist.txt",result , 0666) //写入文件(字节数组)
				}()

				result, _ := json.Marshal(utils.Comresult{
					Code: 200,
					Msg:  "添加成功",
				})
				w.Write(result)
			}
		} else {
			w.Write([]byte("no"))
		}
	}
}

