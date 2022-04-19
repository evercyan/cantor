package config

// Git 配置
// https://docs.github.com/cn/rest/reference/repos#contents
const (
	GitApiURL   = "https://api.github.com/repos/%s/%s/contents/%s"
	GitTagURL   = "https://api.github.com/repos/evercyan/cantor/tags"
	GitFileURL  = "https://cdn.jsdelivr.net/gh/%s/%s/%s"
	GitDBFile   = "resource/cantor.db"
	GitFilePath = "resource/%s/%s%s"
	GitMarkFile = "mark"
	GitMessage  = "upload by cantor"
	GitRepoURL  = "https://github.com/evercyan/cantor"
	GitAppURL   = GitRepoURL + "/releases/tag/%s"
)
