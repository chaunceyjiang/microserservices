package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 本地存储实现
type Local struct {
	basePath string // 基础路径
}

// fullPath 将给的路径追加到基础路径前面
func (l *Local) fullPath(path string) string {
	return filepath.Join(l.basePath,path)
}
func (l *Local) Save(path string, file io.Reader) error {
	// 获取全路径
	fp:=l.fullPath(path)

	// 创建目录全路径
	d:=filepath.Dir(fp)
	err:=os.MkdirAll(d,os.ModePerm)
	if err!=nil{
		return fmt.Errorf("文件目录创建失败: %w",err)
	}
	// 判断文件是否存在
	_,err=os.Stat(fp)
	if err == nil{
		// 文件存在，删除
		if err := os.Remove(fp);err!=nil{
			return fmt.Errorf("删除失败： %w",err)
		}
	}else if !os.IsNotExist(err) {
		// 不是 文件不存在错误
		return fmt.Errorf("其他错误: %w",err)
	}

	// 创建文件
	f,err:=os.Create(fp)
	if err != nil {
		return fmt.Errorf("文件创建失败 %w",err)
	}
	defer f.Close()

	// 拷贝数据
	_,err =io.Copy(f,file)
	if err!=nil{
		return fmt.Errorf("文件写入失败: %w",err)
	}
	return nil
}

func NewLocal(base string) *Local {
	return &Local{basePath: base}
}