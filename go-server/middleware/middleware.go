package middleware

import (
	"awesomeProject/go-server/model"
	"database/sql"
	"fmt"
	_ "fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"html/template"
	"log"
	"net/http"
	"slices"
	"time"
)

const connstr = "user=postgres password=ficha2003 dbname=todo sslmode=disable"

func No_Autorization(w http.ResponseWriter, r *http.Request) {
	file, errr := template.ParseFiles("go-server/templates/home_no_authorization.html")
	if errr != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	file.Execute(w, nil)
}

func Check_Authz_User(token_user string) (int, string) {
	var id_user int
	var username_user string
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		print("Error connect to bd", err.Error())
		panic(err)
	} else {
		print("connected to bd")
	}
	query := `SELECT id, username FROM users WHERE token = $1`
	err = db.QueryRow(query, token_user).Scan(&id_user, &username_user)
	//for rows.Next() {
	//	rows.Scan(&id_user, &username_user)
	//}
	//log.Printff(username_user)
	defer db.Close()
	if err != nil {
		print("Error executing query", err.Error())
		return -1, ""
	} else {
		return id_user, username_user
	}

}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie_user, err := r.Cookie("auth_token")
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Printf("Error give token for user", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id_user_after_check, username_user := Check_Authz_User(cookie_user.Value)
	//log.Printf(id_user_after_check, username_user)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	db1, err := sql.Open("postgres", connstr)
	query := `SELECT id, title, deskription, completed,
       		  created_at, duedate, uuid, user_id 
       		  FROM task WHERE user_id = $1`
	records_task, err := db1.Query(query, id_user_after_check)
	if err != nil {
		log.Fatal(err)
	}
	var tasks []model.Task
	for records_task.Next() {
		var task model.Task
		err := records_task.Scan(&task.Id, &task.Title, &task.Description, &task.Completed,
			&task.Created_at, &task.DueDate, &task.Uuid, &task.UserId)
		if err != nil {
			log.Fatal(err)
			http.Error(w, "Errur scannig database result", http.StatusInternalServerError)
		}
		//if err == records_task.Scan(&task.Id, &task.Title,
		//		&task.Description, &task.Completed,
		//		&task.Created_at, &task.DueDate, &task.UserId) {
		//		log.Printfln(err)
		//		http.Error(w, "Errur scannig database result", http.StatusInternalServerError)
		//		return
		//	}
		tasks = append(tasks, task)
	}
	data := map[string]interface{}{
		"tasks":    tasks,
		"username": username_user,
	}
	file, errr := template.ParseFiles("go-server/templates/home.html")
	if errr != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	file.Execute(w, data)
	defer db.Close()
}

func Add_Task(w http.ResponseWriter, r *http.Request) {

	cookie_user, _ := r.Cookie("auth_token")
	id_user_after_check, _ := Check_Authz_User(cookie_user.Value)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}

	//log.Printf(r.Method)
	file, errr := template.ParseFiles("go-server/templates/add_task.html")
	if errr != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
	file.Execute(w, nil)
}

func Add_TASK_IN_Table(w http.ResponseWriter, r *http.Request) {
	cookie_user, err := r.Cookie("auth_token")
	if err != nil {
		log.Printf("Error give token for user", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id_user_after_check, _ := Check_Authz_User(cookie_user.Value)
	//log.Printf(id_user_after_check, username_user)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	r.FormValue("add_task")
	err = r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	title := r.FormValue("taskName")
	description := r.FormValue("taskDescription")
	due_date := r.FormValue("dueDateTime")
	querry := `INSERT INTO task (title, deskription, 
               Completed, created_at, duedate, user_id )
			   VALUES ($1, $2, $3, $4, $5, $6)`
	db, err := sql.Open("postgres", connstr)
	result, err := db.Exec(querry, title, description, false, time.Now(), due_date, id_user_after_check)
	if err != nil {
		log.Fatal("Error add new task", err)
	}
	rowsAffected, _ := result.RowsAffected()
	fmt.Println(rowsAffected)
	defer db.Close()
	http.Redirect(w, r, "/home", http.StatusFound)
}

func Change_Task(w http.ResponseWriter, r *http.Request) {
	cookie_user, err := r.Cookie("auth_token")
	if err != nil {
		log.Printf("Error give token for user", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id_user_after_check, username_user := Check_Authz_User(cookie_user.Value)
	//log.Printf(id_user_after_check, username_user)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Printf("Error connect bd in change task", err)
	}
	uid_record := mux.Vars(r)["id"]
	querry := ` SELECT title, deskription, 
               Completed, created_at, duedate, uuid, user_id FROM task WHERE uuid=$1`

	if err != nil {
		log.Fatal("Receiving error  record task by its uuid", err)
		http.Redirect(w, r, "/home", http.StatusFound)
	}
	var task model.Task
	err = db.QueryRow(querry, uid_record).Scan(&task.Title, &task.Description, &task.Completed,
		&task.Created_at, &task.DueDate, &task.Uuid, &task.UserId)
	if err != nil {
		log.Fatal(err)
	}
	file, errr := template.ParseFiles("go-server/templates/change_task.html")
	data := map[string]interface{}{
		"username": username_user,
		"task":     task,
	}
	file.Execute(w, data)
	if errr != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

func Save_Changed_Record(w http.ResponseWriter, r *http.Request) {
	cookie_user, err := r.Cookie("auth_token")
	if err != nil {
		log.Printf("Error give token for user", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id_user_after_check, _ := Check_Authz_User(cookie_user.Value)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	uid_task := mux.Vars(r)["id"]
	new_title := r.FormValue("title")
	new_description := r.FormValue("description")
	new_due_date := r.FormValue("due_date")
	var new_completed bool
	if r.FormValue("completed") == "" {
		new_completed = false
	} else {
		new_completed = true
	}
	querry := `UPDATE task 
			   SET title=$1, deskription=$2, Completed=$3,  duedate=$4 
			   WHERE uuid=$5 AND user_id=$6`
	db, err := sql.Open("postgres", connstr)
	_, err = db.Exec(querry, new_title, new_description, new_completed,
		new_due_date, uid_task, id_user_after_check)
	if err != nil {
		log.Printf("Error save changed record", err)
	}
	defer db.Close()
	http.Redirect(w, r, "/home", http.StatusFound)
}

func check_filter_tasks(filter1 string, filter2 string) []string {
	f_1 := []string{"True", "False"}
	f_2 := []string{"Early deadline", "Late deadline", "Early creation", "Late creation"}
	er1 := slices.Contains(f_1, filter1)
	er2 := slices.Contains(f_2, filter2)
	if er2 == true {
		index := slices.Index(f_2, filter2)
		log.Print("index: ", index)
		if index == 0 {
			filter2 = "duedate DESC"
		} else if index == 1 {
			filter2 = "duedate ASC"
		} else if index == 2 {
			filter2 = "created_at DESC"
		} else if index == 3 {
			filter2 = "created_at ASC"
		}
	}
	if (er1 == false) && (er2 == false) {
		res := []string{"Null"}
		return res
	} else if er2 == false {
		res := []string{filter1}
		return res
	} else if er1 == false {
		res := []string{filter2}
		return res
	} else {

		res := []string{filter1, filter2}
		return res
	}
	res := []string{"Null"}
	return res
}

func Filte_Tasks_Home(w http.ResponseWriter, r *http.Request) {
	cookie_user, err := r.Cookie("auth_token")
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		log.Printf("Error give token for user", err)
		http.Redirect(w, r, "/", http.StatusFound)
	}
	id_user_after_check, username_user := Check_Authz_User(cookie_user.Value)
	if id_user_after_check == -1 {
		http.Redirect(w, r, "/", http.StatusFound)
	}
	r.ParseForm()
	flter_1 := r.FormValue("filter1")
	flter_2 := r.FormValue("filter2")

	checked_filter := check_filter_tasks(flter_1, flter_2)
	log.Println(checked_filter[0], checked_filter[1], "HHHHHHHHHHHHHHH")
	if checked_filter[0] == "Null" {
		log.Printf("Error execute filter for the tasks")
		http.Redirect(w, r, "/home", http.StatusFound)
	} else {
		db1, err := sql.Open("postgres", connstr)
		var records_task *sql.Rows
		if (len(checked_filter) == 1) && (len(checked_filter[0]) > 7) {
			query := `SELECT id, title, deskription, completed,
		    		  created_at, duedate, uuid, user_id
					  WHERE user_id = $1
		    		  FROM task ORDER BY ` + checked_filter[0]

			records_task, err = db1.Query(query, id_user_after_check)
		} else if (len(checked_filter) == 1) && (len(checked_filter[0]) < 7) {
			query := `SELECT id, title, deskription, completed,
		    		  created_at, duedate, uuid, user_id
		    		  FROM task
		    		  WHERE user_id = $1 AND completed = $2`
			records_task, err = db1.Query(query, id_user_after_check, checked_filter[0])
		} else {
			query := `SELECT id, title, deskription, completed,
		    		  created_at, duedate, uuid, user_id
		    		  FROM task 
		    		  WHERE user_id = $1 AND completed = $2
		    		  ORDER BY ` + checked_filter[1]
			records_task, err = db1.Query(query, id_user_after_check, checked_filter[0])
		}

		if err != nil {
			log.Fatal(err)
		}
		var tasks []model.Task
		for records_task.Next() {
			var task model.Task
			err := records_task.Scan(&task.Id, &task.Title, &task.Description, &task.Completed,
				&task.Created_at, &task.DueDate, &task.Uuid, &task.UserId)
			if err != nil {
				log.Fatal(err)
				http.Error(w, "Errur scannig database result", http.StatusInternalServerError)
			}
			tasks = append(tasks, task)
		}

		data := map[string]interface{}{
			"tasks":    tasks,
			"username": username_user,
		}
		file, errr := template.ParseFiles("go-server/templates/home.html")
		if errr != nil {
			http.Error(w, "Error rendering template", http.StatusInternalServerError)
		}
		file.Execute(w, data)
		defer db.Close()

	}

}
