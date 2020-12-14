package handlers

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"microserservices/files"
	"net/http"
	"path/filepath"
)

type Files struct {
	store files.Storage
	l     *log.Logger
}

func NewFiles(l *log.Logger, s files.Storage) *Files {
	return &Files{
		store: s,
		l:     l,
	}
}


func (f *Files) UploadREST(w http.ResponseWriter, r *http.Request) {
	vars :=mux.Vars(r)
	// 请求前面做了正则表达式校验
	id:= vars["id"]
	fn:= vars["filename"]
	f.l.Println(" Handle File POST","id",id,"filename",fn)
	if err:=f.saveFile(id,fn,r.Body);err!=nil{
		http.Error(w,"保存文件失败",http.StatusBadRequest)
	}
	resp:=Resp{
		Code:    200,
		Message: "保存文件成功",
		Data:    nil,
	}
	resp.ToJSON(w)
}

func (f *Files) saveFile(id, path string, r io.ReadCloser) error {
	f.l.Println("Save file for Product","id",id,"path",path)
	fp:=filepath.Join(id,path)
	// 保存接口
	return f.store.Save(fp,r)
}


