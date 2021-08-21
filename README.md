<div align="center">
<img src="./appicon.png" width="200" alt="cantor" align=center />

[![goreportcard](https://goreportcard.com/badge/github.com/evercyan/cantor)](https://goreportcard.com/report/github.com/evercyan/cantor)

基于 [wails](https://github.com/wailsapp/wails) + [vue](https://cn.vuejs.org/) + [github-api](https://docs.github.com/cn/rest/reference/repos#contents) 实现的个人图床 mac app

[点我下载 Cantor-v0.0.6.dmg](https://github.com/evercyan/cantor/releases/download/v0.0.6/Cantor-v0.0.6.dmg)
</div>

---

## QA

- Cantor 名字由来
  - 有一位数学家名叫 Georg Cantor
  - 他的成就之一就是集合论
  - Cantor, 意指 "康托尔集", 实指 "图床"


- 使用应用前的准备工作
  - 自己的 GitHub 账号, e.g. `evercyan`
  - 自己的邮箱(commit 使用), e.g. `evercyan@qq.com`
  - 新建一个 GitHub 仓库, e.g. `evercyan/repository`
  - 申请 GitHub access_token [点击申请](https://github.com/settings/tokens)

  
- 如何使用应用
  - 打开应用
  - 设置 GihHub 配置
  ```
  未配置时, 会自动打开配置窗口
  后面可以通过点击设置按钮触发
  配置存储于 `~/.cantor/config.json`
  ```
  - 上传图片
  ```text
  点击左侧上传按钮, 选择图片文件进行添加
  应用会通过 GitHub Api 将图片文件上传至配置中的仓库
  同时会写入本地数据库 ~/.cantor/cantor.db
  并同步 cantor.db 到 仓库/resource/cantor.db
  ```

- blabla...
  - app 实测支持 Mac 10.14+, 其他平台需自行下载源码编译
  - 应用日志存储于 ~/.cantor/app.log
  - 上传图片大小最大为 2M
  - 上传图片格式只支持 png gif jpg jpeg

---

## Snapshot

![cantor-1](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/76/763cda4bd4b0e2fd359799311383cf65.png)

![cantor-2](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/43/431c46df6fc2171e08d6fbfa174d6562.png)

![cantor-3](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/d7/d77e72f7234e37c3c5906023a541fb71.png)
