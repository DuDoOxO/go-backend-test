package model

import "time"

type User struct {
	Name      string    `json:"name"`
	Line      LineInfo  `json:"lineInfo" bson:"line_info"`
	CreatedAt time.Time `json:"createAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

type LineInfo struct {
	LineUserId  string `json:"userId" bson:"line_user_id"`
	DisplayName string `json:"displayName" bson:"display_name"`
	Language    string `json:"language" bson:"language"`
	PicUrl      string `json:"pictureUrl" bson:"pic_url"`
	StatusMsg   string `json:"statusMessage" bson:"status_msg"`
}

type LineMessage struct {
	UserId    string    `json:"userId" bson:"user_id"`
	Message   string    `json:"message" bson:"message"`
	CreatedAt time.Time `json:"createAt" bson:"created_at"`
}
