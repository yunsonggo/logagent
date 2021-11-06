# logagent
### 使用Go语言完成的一个收集管理日志的生态系统(我是 golang 的铁粉)

通过web管理页面设置需要收集日志的项目,设置该项目的日志目录.
当然需要完成必须的配置项.

### 功能及流程

  1. beego    web管理
  2. mysql    存储配置
  3. etcd     保存配置
  4. tail     获取etcd的配置,根据配置收集日志
  5. zookeeper 为kafka提供依赖
  6. kafka    生产tail收集的日志消息到队列
  7. es       消费kafka中的日志消息
  8. kibana   方便搜索及管理日志消息

### 启动 以windows环境为例:

  #### zookeeper
     
      .\bin zKServer.cmd(zKServer.sh)
     
  #### kafka
  
      .\bin\windows kafka-server-start.bat .\config\server.properties
      
  #### etcd
  
      .\bin etcd.exe
      
  #### elasticsearch
  
      .\bin elasticsearch-service.bat
      
  #### kibana
  
      .\bin kibana.bat
      
  #### beego
  
      go run main.go
