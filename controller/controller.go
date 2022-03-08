package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Keshav-Agrawal/mongoapi/model"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

const connectionString = "mongodb+srv://keshav1:keshav1@cluster0.sjkrk.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"
const dbName = "homework"
const colName = "task"

var collection *mongo.Collection

type HomeworkSVC interface {
	GetMyAllTask(w http.ResponseWriter, r *http.Request)
	CreateTask(w http.ResponseWriter, r *http.Request)
	MarkAsDone(w http.ResponseWriter, r *http.Request)
	DeleteATask(w http.ResponseWriter, r *http.Request)
	DeleteAllTask(w http.ResponseWriter, r *http.Request)
}
type homeworkService struct {
	collection *mongo.Collection
}

func NewHomeWorkService() HomeworkSVC {
	c := InitDB()
	return &homeworkService{
		collection: c,
	}

}

//dont use init
func InitDB() *mongo.Collection {
	var err error
	clientOption := options.Client().ApplyURI(connectionString)

	client, err = mongo.Connect(context.TODO(), clientOption)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	return client.Database(dbName).Collection(colName)

}

func insertOneTask(work model.Homework, collection *mongo.Collection) {
	inserted, err := collection.InsertOne(context.Background(), work)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted 1 task in db with id: ", inserted.InsertedID)
}

func updateOneTask(workId string, collection *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex(workId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"done": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count: ", result.ModifiedCount)
}

func deleteOneTask(workId string, collection *mongo.Collection) {
	id, _ := primitive.ObjectIDFromHex(workId)
	filter := bson.M{"_id": id}
	deleteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("task got delete with delete count: ", deleteCount)
}

func deleteAllTask(collection *mongo.Collection) int64 {

	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("NUmber of task delete: ", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

func getAllTask(collection *mongo.Collection) []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var worklist []primitive.M

	for cur.Next(context.Background()) {
		var work bson.M
		err := cur.Decode(&work)
		if err != nil {
			log.Fatal(err)
		}
		worklist = append(worklist, work)

	}

	defer cur.Close(context.Background())
	return worklist
}

func (h homeworkService) GetMyAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allTask := getAllTask(h.collection)
	json.NewEncoder(w).Encode(allTask)

}

func (h homeworkService) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var work model.Homework
	_ = json.NewDecoder(r.Body).Decode(&work)
	insertOneTask(work, h.collection)
	json.NewEncoder(w).Encode(work)

}

func (h homeworkService) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneTask(params["id"], h.collection)
	json.NewEncoder(w).Encode(params["id"])
}

func (h homeworkService) DeleteATask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneTask(params["id"], h.collection)
	json.NewEncoder(w).Encode(params["id"])
}

func (h homeworkService) DeleteAllTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Allow-Control-Allow-Methods", "DELETE")

	count := deleteAllTask(h.collection)
	json.NewEncoder(w).Encode(count)
}

