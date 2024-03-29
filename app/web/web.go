package web

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"gomonitor/app/task"
	"gomonitor/config"
	"gomonitor/utils"
	"io/ioutil"
	"net/http"
	"html/template"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	//日志删除期限
	diffTime = 3600 * 24 * 30
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
	var staticServiceFlag bool
	flag.BoolVar(&staticServiceFlag,"staticflag",false, "Use -staticflag")
	flag.Parse()
	if staticServiceFlag {
		http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir(config.Gconfig.CurExePath))))
	}

	clearlog()
	fmt.Println("listen "+config.Gconfig.WebPort)
	if err := http.ListenAndServe(":"+config.Gconfig.WebPort, nil); err != nil {
		fmt.Println(err)
	}

}
func clearlog()  {
	//定时任务
	ticker := time.NewTicker(time.Hour * 24)
	go func() {
		for range ticker.C {
			fmt.Println("clearlog:" + time.Now().Format("2006-01-02 15:04:05"))
			nowTime := time.Now().Unix()
			err := filepath.Walk(config.Gconfig.CurExePath+"logs/", func(path string, f os.FileInfo, err error) error {
				if f == nil {
					return err
				}
				if f.IsDir() {
					return nil
				}
				fileTime := f.ModTime().Unix()
				if (nowTime - fileTime) > diffTime { //判断文件是否超过7天
					fmt.Printf("Delete file %v !\r\n", path)
					err2:=os.RemoveAll(path)
					if err2 != nil {
						fmt.Println(err2)
					}
				}
				return nil
			})
			if err != nil {
				fmt.Printf("filepath.Walk() returned %v\r\n", err)
			}
		}
	}()
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
						Name: item.Name,
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
		if user=="spirit" && pwd == "adcptbtptp" {
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
	sortlist(list)
	return list
}
func sortlist(list []utils.ExeInfo) {
	for i := 0; i < len(list)-1; i++ {
		for j := i+1; j < len(list); j++ {
			if  list[i].Exeid>list[j].Exeid{
				list[i],list[j] = list[j],list[i]
			}
		}
	}
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
			getall :=r.FormValue("all")
			tline :=r.FormValue("line")
			if exeid !=""{
				if getall == "1"{
					//获取全部log
					result, _ := json.Marshal(utils.Comresult{
						Code: 200,
						Msg:  "success",
						Data:getlogtail(exeid,0),
					})
					w.Write(result)
				}else {
					line := 50
					if tline !=""{
						t,e := strconv.Atoi(tline)
						if e==nil {
							line = t
						}
					}
					//仅获取最新部分
					result, _ := json.Marshal(utils.Comresult{
						Code: 200,
						Msg:  "success",
						Data:getlogtail(exeid,line),
					})
					w.Write(result)
				}
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
					Msg: "no login",
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
func getlogtail(exeid string, line int) string {
	filename := "info_"+exeid +"_" + time.Now().Format("2006-1-2")+".log"
	if line>0 {
		content,err := tail(config.Gconfig.CurExePath+"/logs/"+filename,line)
		if err!=nil{
			fmt.Println("error : %s", err)
			return ""
		}
		return content
	}else{
		bytes, err := ioutil.ReadFile(config.Gconfig.CurExePath+"/logs/"+filename)
		if err != nil {
			fmt.Println("error : %s", err)
			return ""
		}
		return string(bytes)
	}
}
func tail(filepath string, n int) (string,error)  {
	_,err := os.Stat(filepath)
	if err == nil {
		//获取文件总行数
		file,err := os.Open(filepath)
		if err != nil{
			return "",err
		}
		defer file.Close()
		fd:=bufio.NewReader(file)
		count :=0
		for {
			_,err := fd.ReadString('\n')
			if err!= nil{
				break
			}
			count++
		}
		startline :=  0
		if count>n {
			startline = count-n
		}
		//开始读取
		file.Seek(0,0)
		//fd.Reset(file)
		count =0
		content :=""
		for {
			line,err := fd.ReadString('\n')
			if err!= nil{
				break
			}
			if count>=startline {
				content +=line
			}
			count++
		}
		return  content,nil
	}else {
		return "",err
	}
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
		name :=r.FormValue("name")
		if exeid !="" && cmd!=""  && name!="" {
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
					Name: name,
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

