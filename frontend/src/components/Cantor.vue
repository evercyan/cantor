<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-aside class="side-area" style="width: 71px;">
      <el-menu @select="onSelect">
        <el-tooltip effect="dark" content="上传文件" placement="right">
          <el-menu-item index="upload">
            <i class="el-icon-upload"></i>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip effect="dark" content="设置" placement="right">
          <el-menu-item index="setting">
            <i class="el-icon-s-tools"></i>
          </el-menu-item>
        </el-tooltip>
        <el-tooltip effect="dark" content="关于" placement="right">
          <el-menu-item index="about">
            <i class="el-icon-info"></i>
          </el-menu-item>
        </el-tooltip>
      </el-menu>
    </el-aside>

    <el-main class="main-area">
      <el-table
        :data="list.slice((page.num-1)*page.size, page.num*page.size)"
        :fit="true"
        stripe
        style="width: 100%"
      >
        <el-table-column label="图片" width="100">
          <template slot-scope="scope">
            <el-image :src="scope.row.file_url" :preview-src-list="fileUrlList"></el-image>
          </template>
        </el-table-column>
        <el-table-column label="名称">
          <template slot-scope="scope">
            <el-link type="primary" @click="onLink(scope.row.file_url)">{{scope.row.file_name}}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="file_size" label="大小" width="100"></el-table-column>
        <el-table-column prop="create_at" label="时间" width="200"></el-table-column>
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-button
              type="danger"
              icon="el-icon-delete"
              size="mini"
              @click="onDelete(scope.row.file_path)"
            ></el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="page.total"
        :page-size="page.size"
        :current-page.sync="page.num"
        v-if="list.length > 0"
      ></el-pagination>

      <el-drawer :with-header="false" :visible.sync="drawerConfig" size="50%">
        <el-form ref="form" :model="config" :rules="rules" label-width="60px">
          <el-form-item label="仓库" prop="repo">
            <el-input v-model="config.repo" placeholder="cantor"></el-input>
          </el-form-item>
          <el-form-item label="账号" prop="owner">
            <el-input v-model="config.owner" placeholder="evercyan"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="config.email" placeholder="evercyan@qq.com"></el-input>
          </el-form-item>
          <el-form-item label="私钥" prop="access_token">
            <el-input v-model="config.access_token"></el-input>
            <span class="help-text">
              <i class="el-icon-info"></i> 在此设置: https://github.com/settings/tokens
            </span>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSetting" :loading="loading">设置</el-button>
            <el-button type="danger" @click="drawerConfig = false;">取消</el-button>
          </el-form-item>
        </el-form>
      </el-drawer>

      <el-drawer :with-header="false" :visible.sync="drawerAbout" size="50%" class="drawer-about">
        <img src="@/assets/images/logo.png" class="logo" />
        <div>利用 github repo 打造个人图床应用</div>
        <el-divider></el-divider>
        <el-steps :active="3" direction="vertical" class="step-area">
          <el-step title="配置" description="设置 github 相关配置"></el-step>
          <el-step title="上传" description="点击左侧上传按钮上传文件到 repo"></el-step>
          <el-step title="链接" description="点击打开列表中文件, 复制 github 链接"></el-step>
        </el-steps>
      </el-drawer>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: "page-cantor",
  data() {
    return {
      areaStyle: {
        height: "",
      },
      drawerConfig: false,
      drawerAbout: false,
      loading: false,
      config: {},
      list: [],
      fileUrlList: [],
      page: {
        num: 1,
        size: 10,
        total: 0,
      },
      rules: {
        repo: [{ required: true, message: "请输入仓库", trigger: "blur" }],
        owner: [{ required: true, message: "请输入账号", trigger: "blur" }],
        email: [
          {
            required: true,
            type: "email",
            message: "请输入正确的邮箱",
            trigger: "blur",
          },
        ],
        access_token: [
          { required: true, message: "请输入私钥", trigger: "blur" },
        ],
      },
    };
  },
  watch: {
    list(list) {
      var _this = this;
      var fileUrlList = [];
      for (var i = 0; i < list.length; i++) {
        if (fileUrlList.indexOf(list[i].file_url) === -1) {
          fileUrlList.push(list[i].file_url);
        }
      }
      _this.fileUrlList = fileUrlList;
    },
  },
  mounted() {
    var _this = this;
    // 重置 area 高度
    window.addEventListener("resize", _this.refreshAreaHeight);
    _this.refreshAreaHeight();
    _this.configInit();
    _this.listInit();
  },
  methods: {
    configInit: function () {
      var _this = this;
      _this.wails("GetConfig", "", function (result) {
        _this.config = result;
        if (_this.config.repo == "") {
          _this.drawerConfig = true;
        }
      });
    },
    listInit: function () {
      var _this = this;
      _this.loading = true;
      _this.wails(
        "GetUploadList",
        "",
        function (result) {
          _this.loading = false;
          _this.list = result;
          _this.pageInit();
        },
        function (error) {
          _this.loading = false;
          _this.$message.error(error);
        }
      );
    },
    pageInit: function () {
      this.page.total = this.list.length;
      this.page.num = 1;
    },
    refreshAreaHeight: function () {
      this.areaStyle.height = window.innerHeight + "px";
    },
    onSelect: function (action) {
      var _this = this;
      console.log("onSelect", "action", action);
      if (action == "setting") {
        _this.drawerConfig = true;
        return;
      }

      if (action == "about") {
        _this.drawerAbout = true;
        return;
      }

      if (action == "upload") {
        _this.loading = true;
        _this.wails(
          "UploadFile",
          "",
          function (result) {
            console.log("UploadFile result", result);
            _this.loading = false;
            _this.$message.success("上传成功");
            _this.listInit();
          },
          function (error) {
            _this.loading = false;
            _this.$message.error(error);
          }
        );
        return;
      }
    },
    onSetting: function () {
      var _this = this;
      _this.$refs["form"].validate((valid) => {
        if (valid) {
          _this.loading = true;
          _this.wails(
            "SetConfig",
            JSON.stringify(_this.config),
            function (result) {
              _this.loading = false;
              _this.$message.success(result);
              _this.drawerConfig = false;
              _this.listInit();
            },
            function (error) {
              _this.loading = false;
              _this.$message.error(error);
            }
          );
        }
      });
    },
    onDelete: function (filePath) {
      var _this = this;
      _this
        .$confirm("是否确认删除?", "", { type: "warning" })
        .then(() => {
          _this.loading = true;
          _this.wails(
            "DeleteFile",
            filePath,
            function (result) {
              console.log("UploadFile result", result);
              _this.loading = false;
              _this.$message.success("删除成功");
              _this.listInit();
            },
            function (error) {
              _this.loading = false;
              _this.$message.error(error);
            }
          );
        })
        .catch(() => {});
    },
    onLink: function (fileUrl) {
      window.wails.Browser.OpenURL(fileUrl);
    },
  },
};
</script>