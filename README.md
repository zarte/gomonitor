# gomonitor
golang 后台程序执行管理程序
默认首页http://127.0.0.1:5556/
默认用户spirit,adcptbtptp
## 目的
为了不登录服务器，通过浏览器管理服务器上需要常驻后台执行的程序，svn推送更新代码后重启程序以及执行日志查看。
## 运行
monitor.exe  
可选参数  
-staticflag 开启静态服务器功能


## ui界面
### 登录界面
![login](4.png)
### 任务添加
![login](3.png)
### 任务列表
![login](1.png)
### 任务日志
![login](2.png)

## 常用地址
列表  
http://192.168.9.60:5556/exelist 
执行  
http://192.168.9.60:5556/startexe?exeid=1&op=start  
停止  
http://192.168.9.60:5556/startexe?exeid=1&op=stop  

## 更新日志
* 1.0 20200201 基本功能
* 1.1 20200209 ui界面，日志功能，程序列表持久化

210630 1.2
1. 程序列表默认升序排序
2. 新增标识名字段
3. 定时删除30天前的log
4. 修改默认账号
5. 新增静态文件服务器功能  
  
210827 1.3  
1. 实时控制台监控与日志查看分离    
2. 删除命令功能  
3. 命令命存储bug修复  
