package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func cleanImages() {
	now := time.Now()
	files, err := filepath.Glob("files/")
	if err != nil {
		log.Printf("遍历图片文件时发生错误 %v", err)
	}
	for _, file := range files {
		// 获取文件的时间
		createdAt, err := os.Stat(file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// 如果文件的上传修改时间超过1小时，则删除该文件
		if createdAt.ModTime().Before(now.Add(-time.Hour)) {
			fmt.Println("Deleting file:", file)
			os.Remove(file)
		}
	}
}
