# gomonitor
golang 后台程序执行管理系统
访问http://127.0.0.1:5556/
默认用户admin,123456

列表  
http://192.168.9.60:5556/exelist  
添加已改为post
http://192.168.9.60:5556/addexe?exeid=1&cmd=ping baidu.com  
执行  
http://192.168.9.60:5556/startexe?exeid=1&op=start  
停止  
http://192.168.9.60:5556/startexe?exeid=1&op=stop  