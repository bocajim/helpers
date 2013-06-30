package web 

import "helpers/log"
import "net/http"
import "github.com/gorilla/mux"

var router *mux.Router

type WebHandlerFunc func(w http.ResponseWriter,r *http.Request, m *MetaData)

type MetaData struct {
	SessionCookie *http.Cookie
}

func notFoundHandler(w http.ResponseWriter,r *http.Request) {
	http.Error(w,"Sorry, page not found",http.StatusNotFound) 
}

func Initialize(listenAddress string) {

	router = mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	
	http.Handle("/", router)
	
	go func() {
		log.Printf(log.Info,"Starting HTTP server on address: [%s]",listenAddress)
		e:=http.ListenAndServe(listenAddress, nil)
		if e!=nil {
			log.Printf(log.Error,"Error starting HTTP server: %s",e.Error())
		}
	}()
}

func validationHandler(handler WebHandlerFunc, is_public bool) http.HandlerFunc {
	if is_public==true {
		return func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("sessionId")

			m := new(MetaData)
			if err == nil {
				m.SessionCookie=cookie
			} else {
				m.SessionCookie=nil
			}
	        handler(w, r, m)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionId")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		m := new(MetaData)
		m.SessionCookie=cookie
        handler(w, r, m)
	}
}

func RegisterHandlerFunc(url string, handler WebHandlerFunc, is_public bool) {
	router.HandleFunc(url, validationHandler(handler,is_public))
}

func RegisterHandlerFuncWithPrefix(prefix string, handler WebHandlerFunc, is_public bool) {
	router.PathPrefix(prefix).HandlerFunc(validationHandler(handler,is_public))
}

func GetRequestVars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

