package web

import (
	"encoding/json"
	"fmt"
	"gomonitor/app/task"
	"gomonitor/config"
	"gomonitor/utils"
	"net/http"
)

func Starweb()  {
	config.SetConfig()
	http.HandleFunc("/exelist", ExeList)
	http.HandleFunc("/addexe", AddExeList)
	http.HandleFunc("/startexe", task.StartExe)
	if err := http.ListenAndServe(":"+config.Gconfig.WebPort, nil); err != nil {
		fmt.Println(err)
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
func ExeList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		result, err := json.Marshal(&utils.Comresult{
			Code: 200,
			Msg: "success",
			Data: GetExeList(),
		})
		if err != nil {
			w.Write([]byte("api err"))
		} else {
			w.Write(result)
		}
	}
}

func AddExeList(w http.ResponseWriter, r *http.Request) {
	//打印请求的方法
	if r.Method == "GET" {
		query :=r.URL.Query()
		if len(query["exeid"]) >0 && len(query["cmd"])>0 {
			_,res:=config.Gconfig.ExeList.Load(query["exeid"][0])

			if res{
				result, _ := json.Marshal(utils.Comresult{
					Code: 4,
					Msg:  "已存在",
				})

				w.Write(result)
			}else{
				config.Gconfig.ExeList.Store(query["exeid"][0],utils.ExeInfo{
					Exeid: query["exeid"][0],
					Cmd: query["cmd"][0],
					Status : "stop",
				})
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

