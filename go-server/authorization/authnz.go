package authorization

import (
	"awesomeProject/go-server/model"
	"crypto/sha256"
	_ "crypto/sha256"
	"database/sql"
	_ "database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"math/rand"
	"net/http"
	"time"
)

const connstr = "user=postgres password=ficha2003 dbname=todo sslmode=disable"

func AuthN() {

}

func Generate_token() string {
	e := sha256.New()
	tt := fmt.Sprintf("%d%d", rand.Uint64(), rand.Uint32())
	e.Write([]byte(tt))
	token_st := e.Sum(nil)
	token := hex.EncodeToString(token_st)
	return token
}

func Create_table_user() {

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		print("Error connect to bd", err.Error())
		panic(err)
	} else {
		print("connected to bd")
	}
	db.Exec("")

	defer db.Close()
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		print("Error connect to bd", err.Error())
		panic(err)
	} else {
		print("connected to bd")
	}
	createTableUser := `
	CREATE TABLE IF NOT EXISTS users 
	( 				
						id SERIAL PRIMARY KEY,
						username VARCHAR(100),
						email VARCHAR(100) UNIQUE NOT NULL,
						password_hash VARCHAR(100) NOT NULL,
						created_at TIMESTAMP,
						updated_at TIMESTAMP,
						is_active BOOLEAN NOT NULL DEFAULT FALSE,
						Token VARCHAR(100)
						);`

	_, err = db.Exec(createTableUser)
	if err != nil {
		log.Fatal("failed create table: %v", err)
	}
	createTableTask := `       
	CREATE TABLE IF NOT EXISTS task 
	( 				
						id SERIAL PRIMARY KEY,
						title VARCHAR(100) NOT NULL,
						deskription VARCHAR(100),
						completed BOOLEAN NOT NULL DEFAULT FALSE,
						created_at TIMESTAMP,
						duedate TIMESTAMP,
	    				uuid uuid DEFAULT uuid_generate_v4() NOT NULL,
	    				user_id INTEGER NOT NULL,
						CONSTRAINT user_id
	    				FOREIGN KEY (id) 
	    				REFERENCES users(id)
						);`
	_, err = db.Exec(createTableTask)
	if err != nil {
		log.Fatal("failed create table: %v", err)
	}
	func_for_triigger := `
						CREATE OR REPLACE FUNCTION check_deadline()
						RETURNS TRIGGER AS $$
						BEGIN
						IF NEW.duedate <= NOW() + INTERVAL '1 hour' THEN
							NEW.notify := TRUE;
						END IF;
						RETURN NEW;
						END;
						$$ LANGUAGE plpgsql;
						`
	_, err = db.Exec(func_for_triigger)
	if err != nil {
		log.Fatal("failed create func_for_triigger: %v", err)
	}
	querry_triger := `
	CREATE TRIGGER set_notify_tasks_before_deadline
	BEFORE INSERT OR UPDATE ON task
	FOR EACH ROW
	EXECUTE FUNCTION check_deadline();
					`
	_, err = db.Exec(querry_triger)
	if err != nil {
		log.Fatal("failed create trigger: %v", err)
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	h := sha256.New()
	h.Write([]byte(password))
	hash_passwd_st := h.Sum(nil)
	hash_passwd := hex.EncodeToString(hash_passwd_st)
	e := sha256.New()
	tt := fmt.Sprintf("%d%d", rand.Uint64(), rand.Uint32())
	e.Write([]byte(tt))
	token_st := e.Sum(nil)
	token := hex.EncodeToString(token_st)
	create_user := model.User{Username: username, Email: email,
		Password_hash: string(hash_passwd), Created_at: time.Now(),
		Updated_at: time.Now(), Is_active: true, Token: token}
	fmt.Print(create_user)
	err = insertUser(db, create_user)
	defer db.Close()
	if err != nil {
		log.Fatal("Error insert new user", err)
	} else {
		print("Successful add new user!@")
		http.Redirect(w, r, "/home", 301)
	}
}

func insertUser(db *sql.DB, user model.User) error {
	insertSQl := ` INSERT INTO users (username, email, password_hash,
                    created_at, updated_at, is_active, token) 
 					VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(insertSQl, user.Username, user.Email,
		user.Password_hash, user.Created_at, user.Updated_at,
		user.Is_active, user.Token)
	return err
}

func AuthZ(w http.ResponseWriter, r *http.Request) {
	var password_hash string
	fmt.Print("authz")
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	email := r.FormValue("loginEmail")
	password := r.FormValue("loginPassword")
	log.Print(email)
	h := sha256.New()
	h.Write([]byte(password))
	hash_passwd_st := h.Sum(nil)
	user_hash_passwd := hex.EncodeToString(hash_passwd_st)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		print("Error connect to bd for login site", err.Error())
	}
	defer db.Close()
	query := `SELECT password_hash FROM users WHERE email = $1`
	err_find_record := db.QueryRow(query, email).Scan(&password_hash)
	if err_find_record == sql.ErrNoRows {
		log.Print("User with us email not find, email: %s ", email)
		http.Redirect(w, r, "/", 301)
	} else {
		if user_hash_passwd != password_hash {
			log.Print("No correct password")
			http.Redirect(w, r, "/", 301)
		} else {
			new_token_user := Generate_token()
			query_update_token := `
			UPDATE users
			SET token = $2
			WHERE email = $1`
			_, err = db.Exec(query_update_token, email, new_token_user)
			defer db.Close()
			log.Print("Successful login", new_token_user)
			cookie := &http.Cookie{
				Name:     "auth_token",
				Value:    new_token_user,
				Expires:  time.Now().Add(24 * time.Hour),
				HttpOnly: true,
			}

			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/home", 301)
		}
	}
}
