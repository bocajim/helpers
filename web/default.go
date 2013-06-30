package web

import (
	"net/http"
	"strings"
)

var rootPath string

func DefaultHandlerInit(root string, folders string) {
	rootPath=root
	f:=strings.Split(folders,",")
	for _,v:=range f {
		if strings.HasPrefix(v,"/") {
			RegisterHandlerFuncWithPrefix(v,defaultHandler,true)
		} else {
			RegisterHandlerFuncWithPrefix("/"+v,defaultHandler,true)
		}
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request, m *MetaData) {
	if r.Method == "GET" {
	    http.ServeFile(w, r, rootPath+r.URL.Path)
	} else {
		http.Error(w,"Invalid operation, only 'GET' supported.",http.StatusBadRequest)
	}
}
