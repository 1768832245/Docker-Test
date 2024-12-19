package models

import (
	"time"
)

type User struct {
	ID        int        `gorm:"primaryKey;autoIncrement"`                     // 用户唯一标识
	Nickname  string     `gorm:"size:255"`                                     // 用户名
	Username  string     `gorm:"size:255;unique;not null"`                     // 账号，确保唯一
	Password  string     `gorm:"size:255;not null"`                            // 用户密码
	CreatedAt *time.Time `gorm:"type:datetime(3);default:null;autoCreateTime"` // 用户创建时间，允许为空
	UpdatedAt *time.Time `gorm:"type:datetime(3);default:null;autoUpdateTime"` // 用户更新时间，允许为空，自动更新时间
}

type Message struct {
	ID         int        `gorm:"primaryKey;autoIncrement"`                     // 留言唯一标识
	UserID     int        `gorm:"not null"`                                     // 留言的用户ID，外键users id
	Content    string     `gorm:"type:text;not null"`                           // 留言内容
	CreatedAt  *time.Time `gorm:"type:datetime(3);default:null;autoCreateTime"` // 留言时间，允许为空
	UpdatedAt  *time.Time `gorm:"type:datetime(3);default:null;autoUpdateTime"` // 留言更新时间，允许为空，自动更新时间
	IsDeleted  bool       `gorm:"default:false"`                                // 是否删除，逻辑删除，0表示未删除，1表示已删除
	ParentID   *int       `gorm:"default:null"`                                 // 父留言ID，支持回复功能，根留言为NULL
	LikesCount int        `gorm:"default:0"`                                    // 留言点赞数（总共）

	User   User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`   // 外键
	Parent *Message `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"` // 外键
}

//每一条点赞记录

type Like struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    int       `gorm:"not null"`
	MessageID int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:datetime(3);default:null;autoCreateTime"`

	User    User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Message *Message `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE"`
}
