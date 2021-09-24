package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/evercyan/cantor/backend/cfg"
	"github.com/evercyan/letitgo/crypto"
	"github.com/evercyan/letitgo/json"
)

// Git git 配置
type Git struct {
	Repo        string `json:"repo"`
	Owner       string `json:"owner"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

// ----------------------------------------------------------------

// Get 获取文件
func (g *Git) Get(path string) (string, error) {
	resp, err := g.request("GET", g.getApi(path), "")
	if err != nil {
		return "", err
	}
	return g.getResp(resp)
}

// Update 新增或更新文件
func (g *Git) Update(path string, content string) error {
	param := map[string]interface{}{
		"message": cfg.GitMessage,
		"committer": map[string]string{
			"name":  g.Owner,
			"email": g.Email,
		},
		"content": crypto.Base64Encode(content),
	}
	// 如果存在 sha, 更新文件
	sha := g.getSha(path)
	if sha != "" {
		param["sha"] = sha
	}
	resp, err := g.request("PUT", g.getApi(path), crypto.JsonEncode(param))
	if err != nil {
		return err
	}
	_, err = g.getResp(resp)
	return err
}

// Delete 删除文件
func (g *Git) Delete(path string) error {
	sha := g.getSha(path)
	if sha == "" {
		return errors.New("获取文件 sha 失败")
	}
	param := map[string]interface{}{
		"message": cfg.GitMessage,
		"sha":     sha,
	}
	resp, err := g.request("DELETE", g.getApi(path), crypto.JsonEncode(param))
	if err != nil {
		return err
	}
	_, err = g.getResp(resp)
	return err
}

// ----------------------------------------------------------------

// getSha ...
func (g *Git) getSha(path string) string {
	resp, _ := g.Get(path)
	return json.Json(resp).Key("sha").ToString()
}

// getApi ...
func (g *Git) getApi(path string) string {
	return fmt.Sprintf(cfg.GitApiUrl, g.Owner, g.Repo, path)
}

// getResp ...
func (g *Git) getResp(resp string) (string, error) {
	message := json.Json(resp).Key("message").ToString()
	if message != "" {
		return "", errors.New(message)
	}
	return resp, nil
}

// request ...
func (g *Git) request(method string, url string, data string) (string, error) {
	body := bytes.NewReader([]byte(data))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", g.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return string(b), err
}

// ----------------------------------------------------------------

// GetFileUrl 获取文件链接
func (g *Git) GetFileUrl(path string) string {
	return fmt.Sprintf(cfg.GitFileUrl, g.Owner, g.Repo, path)
}

// GetLastVersion 获取应用最后版本号
func (g *Git) GetLastVersion() string {
	resp, err := g.request("GET", cfg.GitTagUrl, "")
	if err != nil {
		return ""
	}
	return json.Json(resp).Index(0).Key("name").ToString()
}
