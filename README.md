<div align="center">
<img src="./appicon.png" width="200" alt="cantor" align=center />


[![goreportcard](https://goreportcard.com/badge/github.com/evercyan/cantor)](https://goreportcard.com/report/github.com/evercyan/cantor)
[![996.icu](https://img.shields.io/badge/link-996.icu-red.svg)](https://996.icu)

基于 [wails](https://github.com/wailsapp/wails) + vue + [github-api](https://docs.github.com/cn/rest/reference/repos#contents) 实现的个人图床 mac app [点我下载](https://github.com/evercyan/cantor/releases/download/v0.0.3/cantor-v0.0.3.tar.gz)
</div>

---

## 必读 QA

```txt
Q: what's cantor?
A: it means "康托尔集", 意指图床
```

```txt
Q: 使用应用前准备工作
A: 在 github 新建仓库和申请 accss_token
```

[点我申请](https://github.com/settings/tokens)

```txt
Q: 使用步骤
- 打开应用, 填写相关配置 (配置存储于 ~/.cantor/config.json 中)
- 点击上传图片文件, 应用会通过 github api 将图片文件上传至配置的仓库中, 并在应用列表中显示, 可直接打开链接复用
- 上传文件列表的数据存储在 仓库名/resource/database.json 中
```

```txt
Q: repo 目录说明
|____resource           图床资源目录
| |____database.json    上传文件记录
|____backend            后端代码目录  
|____frontend           前端代码目录
|____assets             个人资源目录
```

```txt
Q: 其他 blabla..
- 暂只支持 Mac 10.14+
- 应用日志位置在 /tmp/cantor.log
- 上传图片最大 2M, 只支持 png gif jpg jpeg, 有需要可以下载代码自己调整 (./backend/app.go)
```

---

## DIY

### 准备工作

```sh
# 安装 wails
go get -u github.com/wailsapp/wails/cmd/wails
wails -help

# 下载 cantor
git clone https://github.com/evercyan/cantor
# 本 repo 下 resource, assets 均是个人的资源目录, clone 后可直接删除
```

### 浏览器调试

```sh
# 安装前端组件
cd ./cantor/frontend/
npm install
# 启动前端服务
npm run serve

# 启动后端服务
cd ./cantor/
sh run.sh debug

# 打开 http://127.0.0.1:8080/
# 浏览器调试模式下, 不支持上传文件(需要调用系统文件选择功能), 其余功能正常
```

### 编译可执行文件和 app

```sh
# 生成可执行文件 ./build/cantor
sh run.sh test

# 生成 mac app ./build/cantor.app
sh run.sh build
```

---

## Snapshot

![list](https://raw.githubusercontent.com/evercyan/cantor/master/resource/85/8583ac8715210074a080f90111cb55c1.png)

![config](https://raw.githubusercontent.com/evercyan/cantor/master/resource/39/3951a5451f83f22e4a4867dd8bde4b93.png)

![about](https://raw.githubusercontent.com/evercyan/cantor/master/resource/65/65add3fdae4cd2fddd0d711d3863cbc9.png)
