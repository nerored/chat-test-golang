/*
	文件扫描
*/
package cli

import (
	"os"
	"path/filepath"
)

func LoadFileList(dir string, check func(info os.FileInfo) bool) (fileList []string) {
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) (result error) {
		if check != nil && !check(info) {
			return
		}

		fileList = append(fileList, path)
		return
	})

	return
}
