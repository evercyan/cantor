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

resp: {
  "path": "https://raw.githubusercontent.com/evercyan/cantor/master/resource/a7/a7a8bade8d9ae355d5a47f9948b64178.png"
}
```

### 配置 mweb 图床

- fork 当前仓库
- 替换 main.py 的 CANTOR_PREFIX 中的用户名和仓库名
- 如下图配置 mweb 图床

![cantor](https://raw.githubusercontent.com/evercyan/cantor/master/resource/a7/a7a8bade8d9ae355d5a47f9948b64178.png)