# Cantor

> 基于 go wails + vue2 + github api 实现的个人图床 app

[点我下载试用](https://github.com/evercyan/cantor/releases/download/v0.0.2/cantor.tar.gz)

---

#### QA

```
Q: cator?
A: it means "康托尔集", 意指图床

Q: 业务流程?
A: 在 github 仓库图床仓库, 申请 accss_token 后, 在应用中配置好相关信息, 通过调用 github api 上传文件

Q: 系统运行日志
A: /tmp/cantor.log

Q: 支持哪些系统 
A: 仅 Mac 10.14+ 亲测
```

---

#### Run

```sh
# 安装 wails 
go get -u github.com/wailsapp/wails/cmd/wails
wails -help

# 下载 cantor
git clone https://github.com/evercyan/cantor

# 安装前端组件
cd ./cantor/frontend/
npm install

# 启动后端服务
cd ./cantor/
sh run.sh debug

# 启动前端服务
cd ./cantor/frontend
npm run serve

# 打开 http://127.0.0.1:8080/
```

```sh
# 生成可执行文件 ./build/cantor
sh run.sh test 

# 生成 mac app ./build/cantor.app
sh run.sh build
```

---

#### Snapshot

![list](https://raw.githubusercontent.com/evercyan/cantor/master/resource/85/8583ac8715210074a080f90111cb55c1.png)

![config](https://raw.githubusercontent.com/evercyan/cantor/master/resource/39/3951a5451f83f22e4a4867dd8bde4b93.png)

![about](https://raw.githubusercontent.com/evercyan/cantor/master/resource/65/65add3fdae4cd2fddd0d711d3863cbc9.png)

