package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"strconv"
	"math/rand"
)


type Todo struct {
	Key			string		`json:"key"`
	Name		string		`json:"name"`
	Complete	bool		`json:"complete"`
}

var todos []Todo;

func getTodoList(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)

}

func getTodo(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _,item := range todos {
		if item.Key == params["key"]{
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Todo{})

}

func createTodo(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	var todo Todo
	_ = json.NewDecoder(r.Body).Decode(&todo)
	todo.Key = strconv.Itoa(rand.Intn(1000000))
	todos = append(todos, todo)

	json.NewEncoder(w).Encode(todos)

}

func updateTodo(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.Key == params["key"]{

			var newTodos []Todo
			var todo Todo
			_ = json.NewDecoder(r.Body).Decode(&todo)
			todo.Key = params["key"]
			newTodos = append(newTodos, todo)
			newTodos = append(newTodos, todos[index+1:]...)
			todos = append(todos[:index], newTodos...)
			break
		}
	}

	json.NewEncoder(w).Encode(todos)

}

func deleteTodo(w http.ResponseWriter, r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range todos {
		if item.Key == params["key"]{
			todos = append(todos[:index], todos[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(todos)

}


func main() {

	r := mux.NewRouter()

	todos = append(todos, Todo{Key: "1", Name: "Eat", Complete: false})
	todos = append(todos, Todo{Key: "2", Name: "Sleep", Complete: false})

	r.HandleFunc("/api/todos", getTodoList).Methods("GET")
	r.HandleFunc("/api/todos/{key}", getTodo).Methods("GET")
	r.HandleFunc("/api/todos", createTodo).Methods("POST")
	r.HandleFunc("/api/todos/{key}", updateTodo).Methods("PUT")
	r.HandleFunc("/api/todos/{key}", deleteTodo).Methods("DELETE")

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})

	// start server listen
	// with error handling
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(originsOk, headersOk, methodsOk)(r)))

}
