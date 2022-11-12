package version

var Version = "v1"

// バージョン情報を返す。
func GetVersion() string {
	return Version
}

type DetailedInfo struct {
	Author      string `json:"author"`
	Version     string `json:"version"`
	License     string `json:"license"`
	Description string `json:"description"`
}

func GetDetailedInfo() DetailedInfo {
	return DetailedInfo{
		Author:      "Taito Haga",
		Version:     Version,
		License:     "MIT",
		Description: "Edit, Publish, and Search your dictionary online",
	}
}
