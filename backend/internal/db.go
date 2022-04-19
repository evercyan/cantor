package internal

import (
	"github.com/evercyan/brick/xfile"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// File ...
type File struct {
	ID       int    `json:"id" gorm:"column:id;AUTO_INCREMENT;not null"`
	Name     string `json:"file_name" gorm:"column:file_name;not null"`
	Md5      string `json:"file_md5" gorm:"column:file_md5;not null"`
	Size     string `json:"file_size" gorm:"column:file_size;not null"`
	Path     string `json:"file_path" gorm:"column:file_path;not null"`
	CreateAt string `json:"create_at" gorm:"column:create_at;not null"`
}

// TableName ...
func (f *File) TableName() string {
	return "file"
}

// NewDB ...
func NewDB(dbFilePath string) *gorm.DB {
	isExist := xfile.IsExist(dbFilePath)
	db, err := gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: gormLogger.Discard,
	})
	if err != nil {
		panic("创建数据库失败")
	}
	if !isExist {
		// db 文件在 open 前不存在时, 需要创建表
		db.AutoMigrate(&File{})
	}
	return db
}
