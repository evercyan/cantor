package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/evercyan/brick/xencoding"
	"github.com/evercyan/brick/xjson"
	"github.com/evercyan/cantor/config"
)

// Git git 配置
type Git struct {
	Repo        string `json:"repo" yaml:"repo"`
	Owner       string `json:"owner" yaml:"owner"`
	Email       string `json:"email" yaml:"email"`
	AccessToken string `json:"access_token" yaml:"accessToken"`
}

// ----------------------------------------------------------------

// Get 获取文件
func (t *Git) Get(path string) (string, error) {
	resp, err := t.request("GET", t.getApi(path), "")
	if err != nil {
		return "", err
	}
	return t.getResp(resp)
}

// Update 新增或更新文件
func (t *Git) Update(path string, content string) error {
	param := map[string]interface{}{
		"message": config.GitMessage,
		"committer": map[string]string{
			"name":  t.Owner,
			"email": t.Email,
		},
		"content": xencoding.Base64Encode(content),
	}
	// 如果存在 sha, 更新文件
	sha := t.getSha(path)
	if sha != "" {
		param["sha"] = sha
	}
	resp, err := t.request("PUT", t.getApi(path), xencoding.JSONEncode(param))
	if err != nil {
		return err
	}
	_, err = t.getResp(resp)
	return err
}

// Delete 删除文件
func (t *Git) Delete(path string) error {
	sha := t.getSha(path)
	if sha == "" {
		return errors.New("获取文件 sha 失败")
	}
	param := map[string]interface{}{
		"message": config.GitMessage,
		"sha":     sha,
	}
	resp, err := t.request("DELETE", t.getApi(path), xencoding.JSONEncode(param))
	if err != nil {
		return err
	}
	_, err = t.getResp(resp)
	return err
}

// ----------------------------------------------------------------

// getSha ...
func (t *Git) getSha(path string) string {
	resp, _ := t.Get(path)
	return xjson.New(resp).Key("sha").ToString()
}

// getApi ...
func (t *Git) getApi(path string) string {
	return fmt.Sprintf(config.GitApiURL, t.Owner, t.Repo, path)
}

// getResp ...
func (t *Git) getResp(resp string) (string, error) {
	message := xjson.New(resp).Key("message").ToString()
	if message != "" {
		return "", errors.New(message)
	}
	return resp, nil
}

// request ...
func (t *Git) request(method string, url string, data string) (string, error) {
	body := bytes.NewReader([]byte(data))
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("token %s", t.AccessToken))
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
func (t *Git) GetFileUrl(path string) string {
	return fmt.Sprintf(config.GitFileURL, t.Owner, t.Repo, path)
}

// GetLastVersion 获取应用最后版本号
func (t *Git) GetLastVersion() string {
	resp, err := t.request("GET", config.GitTagURL, "")
	if err != nil {
		return ""
	}
	return xjson.New(resp).Index(0).Key("name").ToString()
}

// GetContent ...
func (t *Git) GetContent(path string) string {
	resp, err := t.Get(path)
	if err != nil {
		return ""
	}
	return xencoding.Base64Decode(xjson.New(resp).Key("content").ToString())
}
