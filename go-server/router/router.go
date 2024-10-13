package router

import (
	"awesomeProject/go-server/authorization"
	"awesomeProject/go-server/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

var uid string = "{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}"

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", middleware.No_Autorization).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./go-server/static/"))))
	router.HandleFunc("/authz", authorization.AuthZ).Methods("POST")
	router.HandleFunc("/registration", authorization.RegisterUser).Methods("GET", "POST")
	router.HandleFunc("/home", middleware.Home).Methods("GET")
	router.HandleFunc("/add_task", middleware.Add_Task).Methods("POST", "GET")
	router.HandleFunc("/add_task_in table", middleware.Add_TASK_IN_Table).Methods("POST", "PUT")
	router.HandleFunc("/change_record/"+uid, middleware.Change_Task).Methods("POST", "GET")
	router.HandleFunc("/save_changed_record/"+uid, middleware.Save_Changed_Record).Methods("POST")
	router.HandleFunc("/home/filter_tasks", middleware.Filte_Tasks_Home).Methods("GET", "POST")
	return router
}
