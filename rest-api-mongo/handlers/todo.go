package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	protos "github.com/frost060/go-microservice-basic/basic-messaging-service/protos/notifications"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/logging"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/models"
	"github.com/gorilla/mux"
)

// TodoHandler model
type TodoHandler struct {
	todoRepo         models.TodoRepository
	log              *logging.LogWrapper
	messagingService protos.NotificationClient
}

// NewTodoHandler creates a new TodoHandler for handling route '/todo'
func NewTodoHandler(todoRepo models.TodoRepository, l *logging.LogWrapper, mss protos.NotificationClient) *TodoHandler {
	return &TodoHandler{todoRepo, l, mss}
}

// GetAllTodos returns all the todos
func (t *TodoHandler) GetAllTodos(rw http.ResponseWriter, r *http.Request) {
	rs, err := t.todoRepo.GetAll()
	if err != nil {
		http.Error(rw, "Failed to load database items", http.StatusInternalServerError)
		return
	}

	bs, err := json.Marshal(rs)
	if err != nil {
		http.Error(rw, "Failed to marshal data", http.StatusInternalServerError)
		return
	}

	_, _ = rw.Write(bs)
}

// GetTodo returns a todo given a todo ID
func (t *TodoHandler) GetTodo(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	todoID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID format", http.StatusBadRequest)
		return
	}

	todo, err := t.todoRepo.Get(todoID)
	if err != nil {
		http.Error(rw, "Failed to read database", http.StatusInternalServerError)
		return
	}

	err = todo.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Failed to marshal data", http.StatusInternalServerError)
	}
}

// AddTodo creates a new todo and save it to db
func (t *TodoHandler) AddTodo(rw http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	err := todo.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Error occurred while unmarshalling data", http.StatusBadRequest)
		return
	}

	err = t.todoRepo.Save(&todo)
	if err != nil {
		http.Error(rw, "Error occurred while saving todo to DB", http.StatusInternalServerError)
	}

	rw.WriteHeader(http.StatusOK)
}

// DeleteTodo deletes a todo from db
func (t *TodoHandler) DeleteTodo(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	todoID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(rw, "Invalid ID format", http.StatusBadRequest)
		return
	}

	err = t.todoRepo.Delete(todoID)
	if err != nil {
		http.Error(
			rw,
			fmt.Sprintf("Error occured while deleting todo with ID: %d", todoID),
			http.StatusInternalServerError,
		)
		return
	}
}
