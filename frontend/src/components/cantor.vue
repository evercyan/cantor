<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-aside class="side-area" style="width: 78px;--wails-draggable: drag;">
      <el-menu @select="onSelect">
        <el-menu-item index="refresh" title="Cantor">
          <img src="@/assets/images/logo.png" alt="Cantor" class="logo"/>
        </el-menu-item>
        <el-menu-item index="upload" title="上传1" style="margin-top: 20px;">
          <i class="el-icon-upload"></i>
        </el-menu-item>
        <el-menu-item index="setting" title="配置">
          <i class="el-icon-s-tools"></i>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main class="main-area" style="--wails-draggable:no-drag">
      <el-template>
        <div class="main-drag-area" style="--wails-draggable:drag"></div>
      </el-template>
      <el-table
          :key="tableKey"
          :data="datas.filterList.slice((page.num - 1) * page.size, page.num * page.size)"
          :fit="true"
          stripe
          style="width: 100%"
          class="main-table-area"
          height="1"
          empty-text="快去上传图片吧~"
      >
        <el-table-column label="图片" width="100" align="center">
          <template slot-scope="scope">
            <el-image :src="scope.row['file_url']" :preview-src-list="datas.urlList">
              <div slot="placeholder" class="image-slot">
                <i class="el-icon-loading"></i>
              </div>
              <div slot="error" class="image-slot">
                <i class="el-icon-picture-outline"></i>
              </div>
            </el-image>
          </template>
        </el-table-column>

        <el-table-column
            label="名称"
            prop="file_name"
            align="center"
            :render-header="(h, data) => renderHeader(h, data, '双击文件名称可进行编辑')"
        >
          <template slot-scope="scope">
                        <span v-if="scope.row.is_name_edit">
                          <el-input
                              maxlength="50"
                              placeholder="请输入文件名称"
                              show-word-limit
                              ref="file_name"
                              v-model="scope.row.file_name"
                              style="width: 100%"
                              @blur="goUpdateName(scope.$index, scope.row)"
                          />
                        </span>
            <span @dblclick="onEditName(scope.$index, scope.row)" style="cursor: pointer;" v-else>
                            {{ scope.row.file_name }}
                        </span>
          </template>
        </el-table-column>

        <el-table-column
            prop="file_size"
            label="大小"
            width="100"
            align="center"
        ></el-table-column>

        <el-table-column
            prop="create_at"
            label="时间"
            width="200"
            align="center"
        ></el-table-column>

        <el-table-column label="操作" width="200" align="center">
          <!-- eslint-disable-next-line -->
          <template slot="header" slot-scope="scope">
            <el-input
                maxlength="8"
                placeholder="请输入关键字"
                v-model="form.keyword"
                style="width: 80%"
                @keyup.native="onSearch"
            />
          </template>
          <template slot-scope="scope">
            <el-button
                type="primary"
                icon="el-icon-link"
                size="mini"
                circle
                @click="goCopyUrl(scope.row['file_url'])"
                title="复制链接"
            ></el-button>
            <el-button
                type="danger"
                icon="el-icon-delete"
                size="mini"
                circle
                @click="goDeleteFile(scope.$index, scope.row)"
                title="删除图片"
            ></el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
          background
          layout="sizes, prev, pager, next, total"
          :total="page.total"
          :page-size="page.size"
          :page-sizes="page.sizes"
          :pager-count="page.count"
          :current-page.sync="page.num"
          @size-change="onChangePageSize"
          class="main-page-area"
          v-if="datas.filterList.length > 0"
      ></el-pagination>

      <el-drawer
          :with-header="false"
          :visible.sync="drawer.config"
          size="39%"
          class="config-area"
      >
        <el-form ref="form" :model="datas.config" :rules="rules" label-width="60px">
          <el-form-item label="仓库" prop="repo">
            <el-input v-model="datas.config.repo" placeholder="repository"></el-input>
          </el-form-item>
          <el-form-item label="账号" prop="owner">
            <el-input v-model="datas.config.owner" placeholder="evercyan"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input
                v-model="datas.config.email"
                placeholder="evercyan@qq.com"
            ></el-input>
          </el-form-item>
          <el-form-item label="私钥" prop="access_token">
            <el-input v-model="datas.config.access_token"></el-input>
            <el-link
                @click="onOpenUrl('https://github.com/settings/tokens')"
                type="info"
                :underline="false"
                icon="el-icon-info"
            >
              申请 GitHub 私钥
            </el-link>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="goUpdateConfig" :loading="loading">
              设置
            </el-button>
            <el-button type="danger" @click="drawer.config = false">
              取消
            </el-button>
          </el-form-item>
        </el-form>
      </el-drawer>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: 'page-cantor',
  data() {
    return {
      // wails app
      app: window.go.backend.App,
      tableKey: 0,
      loading: false,
      // 表单
      form: {
        name: '',
        keyword: '',
      },
      // 窗口信息
      areaStyle: {
        height: '',
      },
      // 弹出框
      drawer: {
        config: false,
      },
      // 数据
      datas: {
        config: {},
        list: [],
        filterList: [],
        urlList: [],
      },
      // 分页
      page: {
        num: 1,
        size: 20,
        sizes: [20, 50, 100],
        total: 0,
        count: 5,
      },
      // 校验
      rules: {
        repo: [
          {
            required: true,
            message: '请输入仓库',
            trigger: 'blur',
          },
        ],
        owner: [
          {
            required: true,
            message: '请输入账号',
            trigger: 'blur',
          },
        ],
        email: [
          {
            required: true,
            type: 'email',
            message: '请输入正确的邮箱',
            trigger: 'blur',
          },
        ],
        access_token: [
          {
            required: true,
            message: '请输入私钥',
            trigger: 'blur',
          },
        ],
      },
    };
  },
  mounted() {
    window.addEventListener('resize', this.onRefreshSize);
    this.onRefreshSize();
    this.goGetConfig();

    // 事件监听
    let _this = this;
    window.runtime.EventsOn("event.upload.begin", function () {
      _this.loading = true;
      setTimeout(() => {
        _this.loading = false;
      }, 30000);
    });
    window.runtime.EventsOn("event.upload.success", function (msg) {
      _this.loading = false;
      _this.$message.success(msg);
      _this.goGetList();
    });
    window.runtime.EventsOn("event.upload.fail", function (msg) {
      _this.loading = false;
      _this.$message.error(msg);
    });
  },
  watch: {
    'datas.filterList'(filterList) {
      // 图片链接
      this.datas.urlList = [];
      for (let i = 0; i < filterList.length; i++) {
        let item = filterList[i];
        if (this.datas.urlList.indexOf(item['file_url']) === -1) {
          this.datas.urlList.push(item['file_url']);
        }
      }
      // 页码
      this.page.total = this.datas.filterList.length;
      this.page.num = 1;
    },
  },
  methods: {
    // 获取配置
    goGetConfig: function () {
      this.app.GetConfig().then((resp) => {
        this.log('Go', 'GetConfig', resp);
        this.datas.config = resp.data;
        if (this.datas.config.repo === '') {
          this.drawer.config = true;
        } else {
          this.goGetList();
        }
      });
    },
    // 获取图片列表
    goGetList: function (slow) {
      this.loading = true;
      let duration = 20;
      if (slow) {
        this.datas.list = [];
        this.datas.filterList = [];
        duration = 1000;
      }
      setTimeout(() => {
        this.app.GetList().then((resp) => {
          this.log('Go', 'GetList', resp);
          this.loading = false;
          if (resp.code !== 0) {
            this.$message.error(resp.msg);
            return;
          }
          this.datas.list = resp.data;
          this.datas.filterList = this.datas.list
        });
      }, duration);
    },
    // 上传文件
    goUploadFile() {
      this.loading = true;
      this.app.BatchUploadFile().then((resp) => {
        this.log('Go', 'BatchUploadFile', resp);
        this.loading = false;
        if (resp.code !== 0) {
          this.$message.error(resp.msg);
          return;
        }
        this.$message.success(resp.data);
        this.goGetList();
      });
    },
    // 更新配置
    goUpdateConfig: function () {
      this.$refs['form'].validate((valid) => {
        if (!valid) {
          return;
        }
        this.loading = true;
        this.app.SetConfig(JSON.stringify(this.datas.config)).then(
            (resp) => {
              this.log('Go', 'SetConfig', resp);
              this.loading = false;
              if (resp.code !== 0) {
                this.$message.error(resp.msg);
                return;
              }
              this.$message.success(resp.data);
              this.drawer.config = false;
              this.goGetList();
            }
        );
      });
    },
    // 删除文件
    goDeleteFile: function (index, row) {
      this.$confirm('是否确认删除?', '', {type: 'warning'})
          .then(() => {
            this.loading = true;
            this.app.DeleteFile(row['file_path']).then((resp) => {
              this.log('Go', 'DeleteFile', resp);
              this.loading = false;
              if (resp.code !== 0) {
                this.$message.error(resp.msg);
                return;
              }
              this.$message.success(resp.data);
              this.datas.filterList.splice(index, 1)
            });
          })
          .catch(() => {
          });
    },
    // 复制链接
    goCopyUrl: function (fileUrl) {
      this.app.CopyFileUrl(fileUrl).then((resp) => {
        this.log('Go', 'CopyFileUrl', resp);
        if (resp.code !== 0) {
          this.$message.error(resp.msg);
          return;
        }
        this.$message.success(resp.data);
      });
    },
    // 更新文件名称
    goUpdateName(index, row) {
      this.log('Vue', 'goUpdateName', [index, row]);
      if (row['file_name'] === '') {
        this.$message.error('文件名称不能为空');
        row['file_name'] = this.form.name;
        return;
      }
      row.is_name_edit = false;
      this.refreshTableKey();
      if (this.form.name === row['file_name']) {
        return;
      }
      this.loading = true;
      this.app.UpdateFileName(row['file_path'], row['file_name']).then(
          (resp) => {
            this.log('Go', 'UpdateFileName', resp);
            this.loading = false;
            if (resp.code !== 0) {
              row['file_name'] = this.form.name;
              this.$message.error(resp.msg);
              return;
            }
            this.$message.success(resp.data);
          }
      );
    },

    // ----------------------------------------------------------------

    // 切换页码
    onChangePageSize(val) {
      this.page.size = val;
    },
    // 页面大小调整
    onRefreshSize: function () {
      this.areaStyle.height = window.innerHeight + 'px';
    },
    // 左侧按钮点选
    onSelect: function (action) {
      this.log('Vue', 'onSelect', action);
      switch (action) {
        case 'setting':
          this.drawer.config = true;
          break
        case 'refresh':
          window.runtime.WindowReload();
          break
        case 'upload':
          this.goUploadFile();
          break
      }
    },
    // 打开链接
    onOpenUrl: function (fileUrl) {
      this.log('Vue', 'onOpenUrl', fileUrl);
      window.runtime.BrowserOpenURL(fileUrl);
    },
    // 编辑名称
    onEditName(index, row) {
      this.log('Vue', 'onEditName', [index, row]);
      row.is_name_edit = true;
      this.form.name = row['file_name'];
      this.refreshTableKey();
      setTimeout(() => {
        this.$refs['file_name'].focus();
      }, 20);
    },
    // 渲染列头
    renderHeader(h, {column}, notice) {
      return (
          <span>
                  {notice ? (
                      <el-tooltip effect="dark" content={notice} placement="top-start">
                      <span>
                        {column.label}&nbsp;<i class="el-icon-info" style="color:#409eff; margin-left:5px;"></i>
                      </span>
                      </el-tooltip>
                  ) : (
                      <span> {column.label} </span>
                  )}
                </span>
      );
    },
    // 刷新表格
    refreshTableKey() {
      this.tableKey = new Date().getTime();
    },
    // 搜索
    onSearch() {
      let keyword = this.form.keyword;
      this.log('Vue', 'onSearch', keyword);
      if (keyword === '') {
        this.datas.filterList = this.datas.list;
        return;
      }
      this.datas.filterList = [];
      for (let info of this.datas.list) {
        if (info['file_name'].indexOf(keyword) !== -1) {
          this.datas.filterList.push(info);
        }
      }
    },

    // ----------------------------------------------------------------

    // 输出日志
    log(type, notice, msg) {
      console.log(
          `%c Cantor %c ${type} %c ${notice} `,
          'background: #880000; color: #fff; border-radius: 3px 0px 0px 3px; padding: 1px; font-size: 0.7rem',
          'background: #000088; color: #fff; padding: 1px; font-size: 0.7rem',
          'background: #008800; color: #fff; border-radius: 0px 3px 3px 0px; padding: 1px; font-size: 0.7rem',
      );
      if (msg) {
        console.log(msg);
      }
    },
  },
};
</script>
