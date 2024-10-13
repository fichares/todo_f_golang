/*
package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type User struct {
	Name                  string
	Age                   uint16
	Money                 int16
	Avg_grades, Happiness float64
	Hobbies               []string
}

func (u *User) getAllinfo() string {
	return fmt.Sprintf("User Name is: %s. He is %d and he has Money equal: %d",
		u.Name, u.Age, u.Money)
}

func (u *User) setNewName(newName string) {
	u.Name = newName
}
func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "index.html", nil)
}

//func home_page(w http.ResponseWriter, r *http.Request) {
//	bob := User{"Bob", 25, -50, 4.2, 12.0, []string{"Football", "Basketball"}}
//	//bob.setNewName("Alex")
//	//fmt.Fprint(w, bob.getAllinfo())
//	tmpl, _ := template.ParseFiles("templates/home_page.html")
//	tmpl.Execute(w, bob)
//}ACA

func handleRequest() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}
func main() {
	handleRequest()
}

*/
