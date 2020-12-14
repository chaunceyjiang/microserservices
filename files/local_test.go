package files

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setupLocal(t *testing.T) (*Local, string, func()) {
	dir, err := ioutil.TempDir("", "tests")
	if err != nil {
		t.Fatal(err)
	}
	l := NewLocal(dir)

	return l, dir, func() {
		// 删除 临时文件
		os.RemoveAll(dir)
	}
}

func TestNewLocal(t *testing.T) {
	savePath := "/1/test.png"
	fileContents := "Hello world"
	l, dir, cleanup := setupLocal(t)
	defer cleanup()
	err := l.Save(savePath, bytes.NewBufferString(fileContents))
	assert.NoError(t, err)

	// 检查文件是否创建成功
	f, err := os.Open(filepath.Join(dir, savePath))
	assert.NoError(t, err)

	// 检查文件内容
	d, err := ioutil.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, string(d), fileContents)
}
