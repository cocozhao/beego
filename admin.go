package beego

import (
	"fmt"
	"github.com/astaxie/beego/toolbox"
	"net/http"
	"time"
)

var BeeAdminApp *AdminApp

//func MyFilterMonitor(method, requestPath string, t time.Duration) bool {
//	if method == "POST" {
//		return false
//	}
//	if t.Nanoseconds() < 100 {
//		return false
//	}
//	if strings.HasPrefix(requestPath, "/astaxie") {
//		return false
//	}
//	return true
//}

//beego.FilterMonitorFunc = MyFilterMonitor
var FilterMonitorFunc func(string, string, time.Duration) bool

func init() {
	BeeAdminApp = &AdminApp{
		routers: make(map[string]http.HandlerFunc),
	}
	BeeAdminApp.Route("/", AdminIndex)
	BeeAdminApp.Route("/qps", QpsIndex)
	BeeAdminApp.Route("/prof", ProfIndex)
	BeeAdminApp.Route("/healthcheck", toolbox.Healthcheck)
	BeeAdminApp.Route("/task", toolbox.TaskStatus)
	BeeAdminApp.Route("/runtask", toolbox.RunTask)
	FilterMonitorFunc = func(string, string, time.Duration) bool { return true }
}

func AdminIndex(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Welcome to Admin Dashboard"))
}

func QpsIndex(rw http.ResponseWriter, r *http.Request) {
	toolbox.StatisticsMap.GetMap(rw)
}

func ProfIndex(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	command := r.Form.Get("command")
	if command != "" {
		toolbox.ProcessInput(command, rw)
	} else {
		rw.Write([]byte("request url like '/prof?command=lookup goroutine'"))
	}
}

type AdminApp struct {
	routers map[string]http.HandlerFunc
}

func (admin *AdminApp) Route(pattern string, f http.HandlerFunc) {
	admin.routers[pattern] = f
}

func (admin *AdminApp) Run() {
	if len(toolbox.AdminTaskList) > 0 {
		toolbox.StartTask()
	}
	addr := AdminHttpAddr

	if AdminHttpPort != 0 {
		addr = fmt.Sprintf("%s:%d", AdminHttpAddr, AdminHttpPort)
	}
	for p, f := range admin.routers {
		http.Handle(p, f)
	}
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		BeeLogger.Critical("Admin ListenAndServe: ", err)
	}
}
