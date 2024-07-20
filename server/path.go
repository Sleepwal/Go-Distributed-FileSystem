package server

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

const defaultRootFolderName = "file"

type PathKey struct {
	PathName string
	Filename string
}

type PathTransformFunc func(key string) PathKey

// FirstPathName
// @Description: 返回首路径
// @receiver p
// @return string
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// FullPath
// @Description: 返回完整路径 + 文件名
// @receiver p
// @return string
func (p PathKey) FullPath() string {
	return fmt.Sprintf("/%s/%s", p.PathName, p.Filename)
}

func CasPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i+1)*blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		PathName: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}
