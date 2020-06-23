# coding=utf-8
from flask import Flask, jsonify, request
import json
import os
import hashlib
import shutil

HTTP_HOST = "127.0.0.1"
HTTP_PORT = 7777
CANTOR_PATH = os.path.dirname(os.path.realpath(__file__)) + '/../'
RESOURCE_PATH = CANTOR_PATH + 'resource/'


# 图片上传处理
class Cantor:
    SUPPORT_TYPES = [
        'jpg',
        'png',
        'gif',
    ]

    def __init__(self):
        return

    def upload(self, file):
        print('file', file)

        file_name = file.filename
        file_suffix = file_name.split('.')[1]
        if '.' not in file_name or file_suffix not in self.SUPPORT_TYPES:
            print("only accept file type as " + ",".join(self.SUPPORT_TYPES))
            return {"path": ""}

        file_tmp_path = os.path.join(RESOURCE_PATH, file_name)
        file.save(file_tmp_path)
        file_md5 = self.md5_file(file_tmp_path)

        file_dir = RESOURCE_PATH + file_md5[0:2] + '/'
        if not os.path.exists(file_dir):
            os.makedirs(file_dir)

        file_real_path = file_dir + file_md5 + '.' + file_suffix
        if os.path.exists(file_real_path) is False:
            shutil.move(file_tmp_path, file_real_path)
            # 上传成功后, 调用 shell, 一键发布 git
            deploy = os.system('nohup sh ' + CANTOR_PATH + '/deploy.sh &')
            print('deploy', deploy)
        else:
            os.remove(file_tmp_path)

        path = '/' + file_md5[0:2] + '/' + file_md5 + '.' + file_suffix
        print('path', path)

        return {"path": path}

    def md5_file(self, file_path):
        md5_obj = hashlib.md5()
        with open(file_path, 'rb') as file_obj:
            md5_obj.update(file_obj.read())
        file_md5 = md5_obj.hexdigest()
        return file_md5


# 入口
app = Flask(__name__)


@app.route('/')
def index():
    return 'Hello Cantor!'


@app.route('/upload', methods=['POST'])
def upload():
    file = request.files['file']
    if not file:
        return jsonify({"path": ""})
    resp = (Cantor()).upload(file)
    return jsonify(resp)


if __name__ == '__main__':
    app.run(host=HTTP_HOST, port=HTTP_PORT)
