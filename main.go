// ampel (German for traffic light) is a simple application managing
// shared services in a collaborative environment. A user can reserve a resource
// for a given amount of time and free it after use
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Global list of all services currently known
var Services ServiceList

// Adds commandfline parameters as services in the default status (free)
func AddEnvironments(list []string) {
	for _, envName := range list {
		Services = Services.Append(NewService(envName))
	}
}

func main() {
	AddEnvironments(os.Args[1:])

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/service", Index)
	router.GET("/service/:service", ServiceStatus)
	router.POST("/service/:service/stop", LockService)
	router.POST("/service/:service/go", FreeService)
	router.GET("/list/available", ListAvailable)
	router.GET("/list/reserved", ListReserved)
	router.ServeFiles("/static/*filepath", http.Dir("static/"))
	if _, err := os.Stat("server.crt"); err == nil {
		if _, err := os.Stat("server.key"); os.IsNotExist(err) {
			fmt.Println("Found server.crt but missing server.key")
			os.Exit(10)
		}
		log.Fatal(http.ListenAndServeTLS(":8443", "server.crt", "server.key", router))

	} else {
		log.Fatal(http.ListenAndServe(":8080", router))
	}
}

// Shows a list of all services together with a form to actually reserve or free them
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch outputType := r.URL.Query().Get("output"); outputType {
	case "json":
		outputJSON(w, Services)
	default:
		renderList(w, "Status", "list_services", list_services, Services)
	}
}

// Returns the status of a single service. The HTTP status codes also reflect the service status
// 200/OK is returned, if the resource is free, 423/Locked, if it is locked.
// If the requested resource is not known, 404/Not found is returned.
func ServiceStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var title string
	var tpl string
	var svc *Service

	svc = Services.Find(ps.ByName("service"))

	if svc.Name == ps.ByName("service") {
		svc.UpdateUntil()
		if svc.Status {
			w.WriteHeader(http.StatusLocked)
			title = " is currently locked"
			tpl = svc_locked
		} else {
			title = " is free"
			tpl = svc_free
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		title = " not found"
		tpl = svc_not_found
		fmt.Println(time.Now().String() + " Failed to get the status for non-existing service " + svc.Name)
	}
	renderPage(w, ps.ByName("service"), svc, tpl, title)
}

// Sets the status of a single service to "Locked". The HTTP status codes also reflect the service status
// 403/Denied is returned, if the resource is already locked, 200/OK, if it is was locked successfully.
// If the requested resource is not known, 404/Not found is returned.
func LockService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var title string
	var tpl string
	var svc *Service
	var err error

	svc = Services.Find(ps.ByName("service"))
	if svc.Name == ps.ByName("service") {
		svc.UpdateUntil()
		if svc.Status {
			w.WriteHeader(http.StatusLocked)
			title = " is already locked"
			tpl = svc_locked
			fmt.Println(time.Now().String() + " " + r.FormValue("who") + "tried to lock the already locked service " + svc.Name)
			w.WriteHeader(http.StatusForbidden)
		} else {
			svc.mutex.Lock()
			svc.Status = true
			svc.Who = r.FormValue("who")
			if r.FormValue("until") == "" {
				svc.Until = time.Now().Add(time.Hour)
			} else {
				svc.Until, err = time.Parse("2006-01-02 15:04:05", r.FormValue("until"))
				if err != nil {
					fmt.Println(err)
					svc.Until = time.Now().Add(time.Hour)
				}
			}
			svc.mutex.Unlock()
			title = " is now locked"
			tpl = svc_locked
			fmt.Println(time.Now().String() + " " + r.FormValue("who") + " locks " + svc.Name + " until " + svc.Until.String())
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		title = " not found"
		tpl = svc_not_found
		fmt.Println(time.Now().String() + " " + r.FormValue("who") + " tried to lock non-existing service " + svc.Name)
	}
	renderPage(w, ps.ByName("service"), svc, tpl, title)
}

// Sets the status of a single service to "Free". The HTTP status codes also reflect the service status
// 409/Conflict is returned, if the resource is already free, 200/OK, if it is was freed successfully.
// If the requested resource is not known, 404/Not found is returned.
func FreeService(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var title string
	var tpl string
	var svc *Service

	svc = Services.Find(ps.ByName("service"))
	if svc.Name == ps.ByName("service") {
		svc.UpdateUntil()
		if !svc.Status {
			title = " is already free"
			tpl = svc_free
			fmt.Println(time.Now().String() + " " + r.FormValue("who") + " failed to unlock already free service " + svc.Name)
			w.WriteHeader(http.StatusConflict)
		} else {
			svc.mutex.Lock()
			svc.Status = false
			svc.Who = ""
			svc.mutex.Unlock()
			title = " is now free"
			tpl = svc_free
			fmt.Println(time.Now().String() + " Service " + svc.Name + " unlocked by " + r.FormValue("who"))
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		title = " not found"
		tpl = svc_not_found
		fmt.Println(time.Now().String() + " " + r.FormValue("who") + " failed to unlock non-existing service " + svc.Name)
	}
	renderPage(w, ps.ByName("service"), svc, tpl, title)
}

func ListAvailable(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := Services.StatusList(false)
	switch outputType := r.URL.Query().Get("output"); outputType {
	case "json":
		outputJSON(w, list)
	default:
		renderList(w, "Available", "list_available", list_available, list)
	}
}

func ListReserved(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	list := Services.StatusList(true)
	switch outputType := r.URL.Query().Get("output"); outputType {
	case "json":
		outputJSON(w, list)
	default:
		renderList(w, "Reserved", "list_reserved", list_reserved, list)
	}
}
