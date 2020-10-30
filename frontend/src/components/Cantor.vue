<template>
  <el-container class="page-cantor" v-loading="loading" :style="areaStyle">
    <el-aside class="side-area" style="width: 71px">
      <el-menu @select="onSelect">
        <el-menu-item index="upload" title="上传">
          <i class="el-icon-upload"></i>
        </el-menu-item>
        <el-menu-item index="setting" title="配置">
          <i class="el-icon-s-tools"></i>
        </el-menu-item>
        <el-menu-item index="about" title="关于">
          <i class="el-icon-info"></i>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-main class="main-area">
      <el-table
        :key="tableKey"
        :data="list.slice((page.num - 1) * page.size, page.num * page.size)"
        :fit="true"
        stripe
        style="width: 100%"
      >
        <el-table-column label="图片" width="100">
          <template slot-scope="scope">
            <el-image :src="scope.row.file_url" :preview-src-list="fileUrlList">
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
          :render-header="
            (h, data) => renderHeader(h, data, '双击文件名称进行编辑')
          "
        >
          <template slot-scope="scope">
            <span v-if="scope.row.is_edit">
              <el-input
                maxlength="50"
                placeholder="请输入名称"
                show-word-limit
                ref="file_name"
                v-model="scope.row.file_name"
                style="width: 100%"
                @blur="onUpdateData(scope.$index, scope.row)"
              >
              </el-input>
            </span>
            <span @dblclick="onEditData(scope.$index, scope.row)" v-else>
              {{ scope.row.file_name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="file_size"
          label="大小"
          width="100"
        ></el-table-column>
        <el-table-column
          prop="create_at"
          label="时间"
          width="200"
        ></el-table-column>
        <el-table-column label="操作" width="100">
          <template slot-scope="scope">
            <el-button
              type="primary"
              icon="el-icon-link"
              size="mini"
              circle
              @click="onCopy(scope.row.file_url)"
              title="复制链接"
            ></el-button>
            <el-button
              type="danger"
              icon="el-icon-delete"
              size="mini"
              circle
              @click="onDelete(scope.row.file_path)"
              title="删除"
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

      <el-drawer
        :with-header="false"
        :visible.sync="drawerConfig"
        size="39%"
        class="config-area"
      >
        <el-form ref="form" :model="config" :rules="rules" label-width="60px">
          <el-form-item label="仓库" prop="repo">
            <el-input v-model="config.repo" placeholder="cantor"></el-input>
          </el-form-item>
          <el-form-item label="账号" prop="owner">
            <el-input v-model="config.owner" placeholder="evercyan"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input
              v-model="config.email"
              placeholder="evercyan@qq.com"
            ></el-input>
          </el-form-item>
          <el-form-item label="私钥" prop="access_token">
            <el-input v-model="config.access_token"></el-input>
            <el-link
              @click="onLink('https://github.com/settings/tokens')"
              type="info"
              :underline="false"
              icon="el-icon-info"
              >申请 github 私钥</el-link
            >
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="onSetting" :loading="loading"
              >设置</el-button
            >
            <el-button type="danger" @click="drawerConfig = false"
              >取消</el-button
            >
          </el-form-item>
        </el-form>
      </el-drawer>

      <el-drawer
        :with-header="false"
        :visible.sync="drawerAbout"
        size="39%"
        class="about-area"
      >
        <img src="@/assets/images/logo.png" class="logo" />
        <div>
          <el-link
            @click="onLink('https://github.com/evercyan/cantor')"
            type="primary"
          >
            Cantor {{ version.current }}
          </el-link>
          <br />
          <el-link
            @click="onLink(version.link)"
            type="danger"
            v-if="version.current < version.last"
          >
            升级 {{ version.last }}
          </el-link>
        </div>
        <el-divider></el-divider>
        <el-steps
          :active="3"
          finish-status="process"
          direction="vertical"
          class="step-area"
        >
          <el-step title="配置" description="github 配置"></el-step>
          <el-step title="上传" description="点击上传文件到 repo"></el-step>
          <el-step title="使用" description="复制或打开链接"></el-step>
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
      tableKey: 0,
      editFileName: "",
      version: {
        current: "",
        last: "",
        link: "",
      },
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
      var fileUrlList = [];
      for (var i = 0; i < list.length; i++) {
        if (fileUrlList.indexOf(list[i].file_url) === -1) {
          fileUrlList.push(list[i].file_url);
        }
      }
      this.fileUrlList = fileUrlList;
    },
  },
  mounted() {
    // 重置 area 高度
    window.addEventListener("resize", this.refreshAreaHeight);
    this.refreshAreaHeight();
    this.configInit();
    this.listInit();
  },
  methods: {
    configInit: function () {
      window.backend.App.GetConfig().then((resp) => {
        console.log("GetConfig", resp);
        this.config = resp.data.config;
        if (this.config.repo == "") {
          this.drawerConfig = true;
        }
        this.version = resp.data.version;
        this.version.link =
          "https://github.com/evercyan/cantor/releases/download/" +
          this.version.last +
          "/Cantor-" +
          this.version.last +
          ".dmg";
      });
    },
    listInit: function () {
      this.loading = true;
      window.backend.App.GetUploadList().then((resp) => {
        console.log("GetUploadList", resp);
        this.loading = false;
        if (resp.code != 0) {
          this.$message.error(resp.message);
          return;
        }
        this.list = resp.data;
        this.pageInit();
      });
    },
    pageInit: function () {
      this.page.total = this.list.length;
      this.page.num = 1;
    },
    refreshAreaHeight: function () {
      this.areaStyle.height = window.innerHeight + "px";
    },
    onSelect: function (action) {
      if (action == "setting") {
        this.drawerConfig = true;
        return;
      }
      if (action == "about") {
        this.drawerAbout = true;
        return;
      }
      if (action == "upload") {
        this.loading = true;
        window.backend.App.UploadFile().then((resp) => {
          console.log("UploadFile", resp);
          this.loading = false;
          if (resp.code != 0) {
            this.$message.error(resp.message);
            return;
          }
          this.$message.success(resp.data);
          this.listInit();
        });
      }
    },
    onSetting: function () {
      this.$refs["form"].validate((valid) => {
        if (!valid) {
          return;
        }
        this.loading = true;
        window.backend.App.SetConfig(JSON.stringify(this.config)).then(
          (resp) => {
            console.log("SetConfig", resp);
            this.loading = false;
            if (resp.code != 0) {
              this.$message.error(resp.message);
              return;
            }
            this.$message.success(resp.data);
            this.drawerConfig = false;
            this.listInit();
          }
        );
      });
    },
    onDelete: function (filePath) {
      this.$confirm("是否确认删除?", "", { type: "warning" })
        .then(() => {
          this.loading = true;
          window.backend.App.DeleteFile(filePath).then((resp) => {
            console.log("DeleteFile", resp);
            this.loading = false;
            if (resp.code != 0) {
              this.$message.error(resp.message);
              return;
            }
            this.$message.success(resp.data);
            this.listInit();
          });
        })
        .catch(() => {});
    },
    onLink: function (fileUrl) {
      window.wails.Browser.OpenURL(fileUrl);
    },
    onCopy: function (fileUrl) {
      window.backend.App.CopyFileUrl(fileUrl).then((resp) => {
        console.log("CopyFileUrl", resp);
        if (resp.code != 0) {
          this.$message.error(resp.message);
          return;
        }
        this.$message.success(resp.data);
      });
    },
    onEditData(index, row) {
      console.log("onEditData", index, row);
      row.is_edit = true;
      this.editFileName = row.file_name;
      this.refreshTableKey();

      setTimeout(() => {
        this.$refs["file_name"].focus();
      }, 20);
    },
    onUpdateData(index, row) {
      console.log("onUpdateData", index, row);
      if (row.file_name === "") {
        this.$message.error("文件名称不能为空");
        row.file_name = this.editFileName;
        return;
      }
      row.is_edit = false;
      this.refreshTableKey();
      if (this.editFileName == row.file_name) {
        return;
      }
      this.loading = true;
      window.backend.App.UpdateFileName(row.file_path, row.file_name).then(
        (resp) => {
          console.log("DeleteFile", resp);
          this.loading = false;
          if (resp.code != 0) {
            this.$message.error(resp.message);
            return;
          }
          this.$message.success(resp.data);
        }
      );
    },
    refreshTableKey() {
      this.tableKey = new Date().getTime();
    },
    renderHeader(h, { column }, notice) {
      return (
        <span>
          {notice ? (
            <el-tooltip effect="dark" content={notice} placement="top-start">
              <span>
                {column.label}&nbsp;
                <i
                  class="el-icon-info"
                  style="color:#409eff;margin-left:5px;"
                ></i>
              </span>
            </el-tooltip>
          ) : (
            <span> {column.label} </span>
          )}
        </span>
      );
    },
  },
};
</script>
