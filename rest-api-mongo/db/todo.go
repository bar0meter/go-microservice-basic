package db

import (
	"context"
	log "github.com/frost060/go-microservice-basic/rest-api-mongo/logging"

	"github.com/frost060/go-microservice-basic/rest-api-mongo/models"
	"github.com/frost060/go-microservice-basic/rest-api-mongo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TodoRepo model having todo collection and logger instances
type TodoRepo struct {
	Instance *mongo.Collection
}

// NewTodoRepo creates a new todo repo
func NewTodoRepo(db *mongo.Database) *TodoRepo {
	instance := db.Collection(TODO)
	return &TodoRepo{instance}
}

// GetAll returns all the todo from db
func (t *TodoRepo) GetAll() (*[]models.Todo, error) {
	var todos []models.Todo
	cursor, err := t.Instance.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var todo models.Todo
		err = cursor.Decode(&todo)
		if err != nil {
			continue
		}

		todos = append(todos, todo)
	}

	return &todos, nil
}

// FindByID returns a todo from db given todo id
func (t *TodoRepo) Get(id int64) (*models.Todo, error) {
	var todo models.Todo

	err := t.Instance.FindOne(context.TODO(),
		bson.M{"_id": id}).Decode(&todo)
	if err != nil {
		log.Error("Error occurred while querying DB", "error", err)
		return nil, err
	}

	return &todo, nil
}

// Save saves a todo to the db
func (t *TodoRepo) Save(todo *models.Todo) error {
	todo.ID = utils.GenerateID()
	insertResult, err := t.Instance.InsertOne(context.TODO(), &todo)
	if err != nil {
		log.Error("Error while saving document to todo db", "error", err)
		return err
	}

	log.Info("Inserted a single todo document: ", insertResult.InsertedID)
	return nil
}

// Delete deletes a todo from db given a todo id
func (t *TodoRepo) Delete(id int64) error {
	_, err := t.Instance.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
