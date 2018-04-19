# 部署同步工具

## 用途

方便开发人员部署用，在jenkins太小题大作时

## 使用说明

`./main -i 忽略文件.txt -u 远程主机用户名  -s 110.110.110.110:22 本地目录 远程目录`


## 目前已实现的功能

- [x] md5对比相同后，忽略上传
- [x] 忽略指定文件同步

## 计划实现的功能

- [] 定制tomcat同步部署，智能忽略配置，重启tomcat
- [] 同步通知，钉钉,tower, worktile类机器人
....