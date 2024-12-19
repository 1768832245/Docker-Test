package flag

import (
	"RR7/models"
	"fmt"
)

func DatabaseAutoMigrate() {
	var err error

	//自动建表**
	err = models.DB.Set("gorm:table_option", "Engine=InnoDB").
		AutoMigrate(
			&models.User{},
			&models.Message{},
			&models.Like{},
		)

	if err != nil {
		fmt.Println("自动建表失败")
	} else {
		fmt.Println("自动建表成功")
	}
}
