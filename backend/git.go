package backend

import (
	"errors"
	"fmt"

	"github.com/evercyan/letitgo/crypto"
	ljson "github.com/evercyan/letitgo/json"
	"github.com/evercyan/letitgo/request"
)

var (
	gitApiUrl  = "https://api.github.com/repos/%s/%s/contents/%s?access_token=%s"
	gitFileUrl = "https://raw.githubusercontent.com/%s/%s/master/%s"
	gitMessage = "auto deploy"
)

type Git struct {
	Repo        string `json:"repo"`
	Owner       string `json:"owner"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
}

func (g *Git) resp(resp string) (string, error) {
	message := ljson.Json(resp).Key("message").ToString()
	if message != "" {
		return "", errors.New(message)
	}
	return resp, nil
}

func (g *Git) Get(path string) (string, error) {
	url := fmt.Sprintf(gitApiUrl, g.Owner, g.Repo, path, g.AccessToken)
	resp, err := request.Get(url)
	Log().Info("Git Get ", "resp ", resp, " err ", err)
	if err != nil {
		return "", err
	}
	return g.resp(resp)
}

func (g *Git) Update(path string, content string) error {
	param := map[string]interface{}{
		"message": gitMessage,
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
	Log().Info("Git Update ", "resp ", resp, " err ", err)
	if err != nil {
		return err
	}
	_, err = g.resp(resp)
	return err
}

func (g *Git) Delete(path string) error {
	sha := g.Sha(path)
	if sha == "" {
		return errors.New("获取文件 sha 失败")
	}
	param := map[string]interface{}{
		"message": gitMessage,
		"sha":     sha,
	}
	resp, err := request.Request("DELETE", g.Api(path), crypto.JsonEncode(param))
	Log().Info("Git Delete ", "resp ", resp, "err ", err)
	if err != nil {
		return err
	}
	_, err = g.resp(resp)
	return err
}

func (g *Git) Sha(path string) string {
	resp, _ := g.Get(path)
	return ljson.Json(resp).Key("sha").ToString()
}

func (g *Git) Api(path string) string {
	return fmt.Sprintf(gitApiUrl, g.Owner, g.Repo, path, g.AccessToken)
}

func (g *Git) Url(path string) string {
	return fmt.Sprintf(gitFileUrl, g.Owner, g.Repo, path)
}
