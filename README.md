<div align="center">
<img src="./appicon.png" width="200" alt="cantor" align=center />

[![goreportcard](https://goreportcard.com/badge/github.com/evercyan/cantor)](https://goreportcard.com/report/github.com/evercyan/cantor)

一个简单好用的图床应用

</div>

---

## 安装

[点击下载](https://github.com/evercyan/cantor/releases/download/v0.1.0/Cantor-v0.1.0.dmg) Mac Cantor-v0.1.0.dmg

或者 `git clone` 源码自行编译, 可支持多平台

---

## 准备

- GitHub 账号, e.g. `evercyan`
- 邮箱, e.g. `evercyan@qq.com`
- 仓库, e.g. `evercyan/repository`
- 申请 GitHub access_token [点击申请](https://github.com/settings/tokens)

---

## 使用

- 打开应用

- GitHub 配置
    ```text
    未配置时, 会自动打开编辑窗口, 也可以通过点击设置按钮触发编辑窗口
    配置存储于 `~/.cantor/config.yaml`
    ```

- 上传图片
    ```text
    点击上传按钮, 或者菜单-文件-上传图片, 可批量选择图片进行上传
    单次最多可上传 10 张图片, 单张图片最大支持 4M, 图片格式仅支持 png gif jpg jpeg
    ```

- 终端使用

    ```shell
    # 需要安装应用并完成 GitHub 图床相关配置才可在终端使用
  
    # 安装 cantor 终端程序
    go install github.com/evercyan/cantor/cmd/cantor@latest
    # 上传文件
    cantor upload ~/demo.png ~/demo1.png
    
    # 安装终端程序后, 可在 typora 里如下图所示配置图片上传, 其他应用亦是同理
    ```

![typora-1](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/93/9355e89e0a615750862320fba0db2bbd.png)

![typora-2](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/56/5687102c40a183568eb6d5e35e8612b3.png)


---

## 快照

![cantor](https://cdn.jsdelivr.net/gh/evercyan/repository/resource/1c/1c98042d7f58d999fdd080dc6bdf68aa.png)
