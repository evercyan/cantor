package git

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/evercyan/cantor/backend/configs"
	"github.com/evercyan/letitgo/crypto"
	j "github.com/evercyan/letitgo/json"
	"github.com/evercyan/letitgo/request"
)

// Git ...
type Git struct {
	Repo        string `json:"repo"`
	Owner       string `json:"owner"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// resp ...
func (g *Git) resp(resp string) (string, error) {
	message := j.Json(resp).Key("message").ToString()
	if message != "" {
		return "", errors.New(message)
	}
	return resp, nil
}

// Get 获取文件
func (g *Git) Get(path string) (string, error) {
	url := fmt.Sprintf(configs.GitApiUrl, g.Owner, g.Repo, path, g.AccessToken)
	resp, err := request.Get(url)
	if err != nil {
		return "", err
	}
	return g.resp(resp)
}

// Update 新增或更新文件
func (g *Git) Update(path string, content string) error {
	param := map[string]interface{}{
		"message": configs.GitMessage,
		"committer": map[string]string{
			"name":  g.Owner,
			"email": g.Email,
		},
		"content": crypto.Base64Encode(content),
	}
	// 如果存在 sha, 更新文件
	sha := g.Sha(path)
	if sha != "" {
		param["sha"] = sha
	}
	resp, err := request.Request("PUT", g.Api(path), crypto.JsonEncode(param))
	if err != nil {
		return err
	}
	_, err = g.resp(resp)
	return err
}

// Delete 删除文件
func (g *Git) Delete(path string) error {
	sha := g.Sha(path)
	if sha == "" {
		return errors.New("获取文件 sha 失败")
	}
	param := map[string]interface{}{
		"message": configs.GitMessage,
		"sha":     sha,
	}
	resp, err := request.Request("DELETE", g.Api(path), crypto.JsonEncode(param))
	if err != nil {
		return err
	}
	_, err = g.resp(resp)
	return err
}

// Sha ...
func (g *Git) Sha(path string) string {
	resp, _ := g.Get(path)
	return j.Json(resp).Key("sha").ToString()
}

// Api ...
func (g *Git) Api(path string) string {
	return fmt.Sprintf(configs.GitApiUrl, g.Owner, g.Repo, path, g.AccessToken)
}

// Url ...
func (g *Git) Url(path string) string {
	return fmt.Sprintf(configs.GitFileUrl, g.Owner, g.Repo, path)
}

// LastVersion ...
func (g *Git) LastVersion() string {
	resp, err := request.Get(configs.GitTagUrl)
	if err != nil {
		return ""
	}
	return j.Json(resp).Index(0).Key("name").ToString()
}

// UploadFileList ...
func (g *Git) UploadFileList() []map[string]string {
	resp, _ := g.Get(configs.GitDBFile)
	list := []map[string]string{}
	if resp != "" {
		content := crypto.Base64Decode(j.Json(resp).Key("content").ToString())
		json.Unmarshal([]byte(content), &list)
	}
	return list
}
