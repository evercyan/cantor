# Cantor

> a self-use image repo
> 

### 说明
- 配合 mweb + github 打造个人图床服务
- mweb 请求 cantor 上传图片, cantor 处理完, 自动调用 deploy.sh 进行 github 发布
- 因为是实时 deploy, 接口响应速度略慢, 待优化

### 启动服务
```shell
# 安装依赖
pip3 install -r requirements.txt
# 启动服务
sh run.sh
# 调试服务
python3 ./src/main.py
``` 

### 配置 mweb 图床

> fork 当前仓库, 替换上图中的 "图片 URL 前缀" 中的用户名为自己的 github 用户名即可

![](./assets/mweb-config.png)