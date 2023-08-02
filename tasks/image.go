package tasks

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func cleanImages() {
	now := time.Now()
	err := filepath.Walk("files/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果文件的上传修改时间超过1小时，则删除该文件
		if info.ModTime().Before(now.Add(-time.Hour)) {
			log.Println("Deleting file:", path)
			os.Remove(path)
		}
		return nil
	})
	if err != nil {
		log.Printf("error when deleting thumbnails: %v", err)
	}
}
