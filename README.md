# Cantor

> self-use 图床服务

---

### 流程一下

- http 请求 cantor 上传图片
- cantor 按规则存储图片
- cantor 调用 deploy.sh 自动发布
- cantor 返回远程图片地址

### 备注一下

```shell
# 安装依赖
pip3 install -r requirements.txt

# 启动服务
sh run.sh

# 调试服务
python3 ./main.py
``` 

```shell
# 可直接调用 api 访问
curl POST 'http://127.0.0.1:7777/upload' --form 'file=@/tmp/abc.png'

{
  "path": "https://raw.githubusercontent.com/evercyan/cantor/master/resource/b0/b0a94e0bf957bbc6bfcb8504953b6ae7.png"
}
```

### 配置 mweb

- fork 当前仓库
- 替换 main.py 的 CANTOR_PREFIX 中的用户名和仓库名
- 增加 mweb 图床配置

![cantor](https://raw.githubusercontent.com/evercyan/cantor/master/resource/b0/b0a94e0bf957bbc6bfcb8504953b6ae7.png)