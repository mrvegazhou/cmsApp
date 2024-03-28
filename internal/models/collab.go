package models

type CollabArticleInfo struct {
	TokenUrl string     `label:"协作地址" json:"tokenUrl"`
	Info     AppArticle `label:"文章详情" json:"info"`
}

type CollabTokenInfo struct {
	RoomName    string `label:"共享名称" json:"roomName"`
	UserName    string `label:"协作成员" json:"userName"`
	CursorColor string `label:"标识颜色" json:"cursorColor"`
	IsCollab    bool   `label:"是否共享" json:"isCollab"`
	Token       string `label:"token" json:"token"`
	IsMe        bool   `lable:"是否为本人" josn:"isMe"`
}

type CollabToken struct {
	Token string `label:"协作验证码" json:"token"`
}
