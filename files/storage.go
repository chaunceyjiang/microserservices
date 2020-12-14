package files

import "io"

// Storage 存储接口
// 方便定义本地存储，或者云存储
type Storage interface {
	Save(path string,file io.Reader) error
}
