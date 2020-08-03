# Cantor

> self-use 图床服务

---

### 说明一下

- 配合 mweb + github 打造个人图床服务
- mweb 请求 cantor 上传图片, cantor 处理完, 自动调用 deploy.sh 进行 github 发布

### 启动服务

```shell
# 安装依赖
pip3 install -r requirements.txt
# 启动服务
sh run.sh
# 调试服务
python3 ./src/main.py
``` 

```shell
# 可直接调用 api 访问
curl POST 'http://127.0.0.1:7777/upload' --form 'file=@/tmp/abc.png'
```

### 配置 mweb 图床

- fork 当前仓库
- 替换下图中的 "图片 URL 前缀" 中的用户名为自己的 github 用户名即可
- https://raw.githubusercontent.com/evercyan/cantor/master/resource

![cantor](https://raw.githubusercontent.com/evercyan/cantor/master/resource/0d/0de45d50af211316ea73ab7350202866.png)