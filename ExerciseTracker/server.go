/*
AUTOR: Mario Reiley Backend Engeneer
NOTE: It is my port in Golang 1.19

# Exercise Tracker

Build a full stack JavaScript app that is functionally similar to this: https://exercise-tracker.freecodecamp.rocks.
Working on this project will involve you writing your code using one of the following methods.
Your responses should have the following structures.

Exercise:

	{
	  username: "fcc_test",
	  description: "test",
	  duration: 60,
	  date: "Mon Jan 01 1990",
	  _id: "5fb5853f734231456ccb3b05"
	}

User:

	{
	  username: "fcc_test",
	  _id: "5fb5853f734231456ccb3b05"
	}

Log:

	{
	  username: "fcc_test",
	  count: 1,
	  _id: "5fb5853f734231456ccb3b05",
	  log: [{
	    description: "test",
	    duration: 60,
	    date: "Mon Jan 01 1990",
	  }]
	}
*/
package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type exerciseLog struct {
	Description string    `bson:"description" json:"description"`
	Duration    float64   `bson:"duration" json:"duration"`
	Date        time.Time `bson:"date" json:"date"`
}

type jerror struct {
	Msj string `json:"msg"`
}

type dataUser struct {
	Username string `bson:"username"`
}

type responseUser struct {
	ObjectID string `bson:"_id" json:"_id"`
	Username string `bson:"username" json:"username"`
}

type dataExcersise struct {
	UserId      string    `bson:"userId" json:"userId"` // UserId = ObjectID
	Description string    `bson:"description" json:"description"`
	Duration    float64   `bson:"duration" json:"duration"`
	Date        time.Time `bson:"date" json:"date"`
}

/*
  - Little trick do not like it but resolve this challengue
    ever id param is in 3 position.
*/
func getId(path string) string {
	mp := strings.Split(path, "/")
	return mp[3]
}

// get('/api/users'
// You can make a GET request to /api/users to get a list of all users.
func users(w http.ResponseWriter, req *http.Request) {

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: "Can not connect to MongoDB"})
		w.Write(res)
		return

	}
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := db.Database("myFirstDatabase").Collection("users")
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}
	var documents []bson.M
	if err = cur.All(context.TODO(), &documents); err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return
	}

	res, _ := json.Marshal(documents)

	w.Write(res)
}

// Built mongoDb query
func seek(from, to time.Time, _id string) (seek primitive.D) {

	if from.IsZero() == false && to.IsZero() == false {
		seek = bson.D{
			{Key: "$and",
				Value: bson.A{
					bson.D{{Key: "userId", Value: _id}},
					bson.D{{Key: "date", Value: bson.D{{Key: "$gte", Value: from}}}},
					bson.D{{Key: "date", Value: bson.D{{Key: "$lte", Value: to}}}},
				},
			},
		}
	} else {
		seek = bson.D{{Key: "userId", Value: _id}}
	}

	return seek
}

// get('/api/users/:_id/logs'
// You can make a GET request to /api/users/:_id/logs to retrieve a full exercise log of any user.
// GET user's exercise log: GET /api/users/:_id/logs?[from][&to][&limit]
// [ ] = optional
// from, to = dates (yyyy-mm-dd); limit = number
func logs(w http.ResponseWriter, req *http.Request) {
	var from time.Time
	var to time.Time
	var limit int64
	var exerciseDoc []exerciseLog
	var user bson.M
	var userName string
	var _seek primitive.D

	_id := getId(req.URL.Path)

	// time.Time{} is 0 time
	from, e := time.Parse("2006-01-02", req.URL.Query().Get("from"))
	if e != nil {
		from = time.Time{}

	}
	to, _ = time.Parse("2006-01-02", req.URL.Query().Get("to"))
	if e != nil {
		to = time.Time{}

	}

	limit, _ = strconv.ParseInt(req.URL.Query().Get("limit"), 0, 64)

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	users := db.Database("myFirstDatabase").Collection("users")
	filter, _ := primitive.ObjectIDFromHex(_id)
	found := users.FindOne(context.TODO(), bson.D{{Key: "_id", Value: filter}}).Decode(&user)
	if found != nil {
		res, _ := json.Marshal(jerror{Msj: "User not found"})
		w.Write(res)
		return

	}
	// result is Allways 1 Document
	for _, v := range user {
		userName, _ = v.(string)
	}

	excersices := db.Database("myFirstDatabase").Collection("excersises")

	opt := options.Find().SetLimit(limit)
	// make logs filter
	_seek = seek(from, to, _id)
	cur, err := excersices.Find(context.TODO(), _seek, opt)
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}

	if err = cur.All(context.TODO(), &exerciseDoc); err != nil {
		res, _ := json.Marshal(jerror{Msj: "mi error" + err.Error()})
		w.Write(res)
		return
	}

	// sentd data to client
	count := len(exerciseDoc)
	res, _ := json.Marshal(struct {
		Username string
		Count    int
		Id       string `json:"_id"`
		Log      *[]exerciseLog
	}{userName, count, _id, &exerciseDoc})

	w.Write(res)

}

// post('/api/users'
// You can POST to /api/users with form data username to create a new user.
func newUser(w http.ResponseWriter, req *http.Request) {
	user := req.PostFormValue("username")

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	collection := db.Database("myFirstDatabase").Collection("users")
	data := dataUser{Username: user}
	mongoId, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}

	mongoIdStr := mongoId.InsertedID.(primitive.ObjectID).Hex() // just this part:"63a06d1c3ac6553ce35b44fd"

	res, _ := json.Marshal(responseUser{ObjectID: mongoIdStr, Username: user})

	w.Write(res)
}

// post('/api/users/:_id/exercises'
// You can POST to /api/users/:_id/exercises with form data description, duration, and optionally date.
// If no date is supplied, the current date will be used.
func exercises(w http.ResponseWriter, req *http.Request) {
	var date time.Time
	var user bson.M
	var userName string

	userId := getId(req.URL.Path)
	description := req.PostFormValue("description")
	duration, _ := strconv.ParseFloat(req.PostFormValue("duration"), 64)

	if len(req.PostFormValue("date")) > 0 {
		date, _ = time.Parse("2006-01-02", req.PostFormValue("date"))
	} else {
		date = time.Now()
	}

	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGODB_URI))
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}
	defer func() {
		if err := db.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	users := db.Database("myFirstDatabase").Collection("users")

	filter, _ := primitive.ObjectIDFromHex(userId)
	found := users.FindOne(context.TODO(), bson.D{{Key: "_id", Value: filter}}).Decode(&user)
	if found != nil {
		res, _ := json.Marshal(jerror{Msj: "User not found"})
		w.Write(res)
		return

	}

	// result is Allways one document
	for _, v := range user {
		userName, _ = v.(string)
	}

	//  ---- codigo correcto solo descomentarizar
	collection := db.Database("myFirstDatabase").Collection("excersises")
	data := dataExcersise{UserId: userId, Description: description, Duration: duration, Date: date}
	mongoId, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		res, _ := json.Marshal(jerror{Msj: err.Error()})
		w.Write(res)
		return

	}
	_ = mongoId // ignore

	res, _ := json.Marshal(struct {
		Username    string
		Description string
		Duration    float64
		Date        string
		Id          string `json:"_id"`
	}{userName, description, duration, date.UTC().Format(time.RFC1123)[:16], userId})

	w.Write(res)
}

const MONGODB_URI = "mongodb+srv://sdkmarior:A123456789_05@cluster0.et10y.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// index
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type:", " text/css; charset=UTF-8")
		r.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./views"))))
	})

	r.Get("/api/users", users)
	r.Get("/api/users/*", logs)
	r.Post("/api/users", newUser)
	r.Post("/api/users/{id}/exercises", exercises)

	http.ListenAndServe("localhost:8080", r)

}
