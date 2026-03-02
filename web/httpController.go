package web

import (
	"log"
	"net/http"
	"text/template"

	"github.com/andeya/pholcus/app"
	"github.com/andeya/pholcus/common/session"
	"github.com/andeya/pholcus/config"
	"github.com/andeya/pholcus/logs"
	"github.com/andeya/pholcus/runtime/status"
)

var globalSessions *session.Manager

func init() {
	var err error
	globalSessions, err = session.NewManager("memory", `{"cookieName":"pholcusSession", "enableSetCookie,omitempty": true, "secure": false, "sessionIDHashFunc": "sha1", "sessionIDHashKey": "", "cookieLifeTime": 157680000, "providerConfig": ""}`)
	if err != nil {
		log.Fatal(err)
	}
	// go globalSessions.GC()
}

func web(rw http.ResponseWriter, req *http.Request) {
	sess, _ := globalSessions.SessionStart(rw, req)
	defer sess.SessionRelease(rw)
	index, err := viewsFS.ReadFile("views/index.html")
	if err != nil {
		logs.Log.Error("read index.html: %v", err)
		http.Error(rw, "internal error", http.StatusInternalServerError)
		return
	}
	t, err := template.New("index").Parse(string(index))
	if err != nil {
		logs.Log.Error("%v", err)
	}
	data := map[string]interface{}{
		"title":   config.NAME,
		"logo":    config.ICON_PNG,
		"version": config.VERSION,
		"author":  config.AUTHOR,
		"mode": map[string]int{
			"offline": status.OFFLINE,
			"server":  status.SERVER,
			"client":  status.CLIENT,
			"unset":   status.UNSET,
			"curr":    app.LogicApp.GetAppConf("mode").(int),
		},
		"status": map[string]int{
			"stopped": status.STOPPED,
			"stop":    status.STOP,
			"run":     status.RUN,
			"pause":   status.PAUSE,
		},
		"port": app.LogicApp.GetAppConf("port").(int),
		"ip":   app.LogicApp.GetAppConf("master").(string),
	}
	t.Execute(rw, data)
}
