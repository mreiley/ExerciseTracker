# Exercise Tracker Microservice

### www.freecodecamp.com proposes this project in JavaScript so I did it in Golang.

Description: Here I used MongoDB Atlas services like repository and did some CRUD operations, additionaly  I try CHI, a lightweight, idiomatic and 
composable router for building Go HTTP services. 


For practique diferentes Golang Packages. I used:

```
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
```

### Task
* You can POST to /api/users with form data username to create a new user.
* The returned response from POST /api/users with form data username will be an object with username and _id properties.
*  You can make a GET request to /api/users to get a list of all users.
* The GET request to /api/users returns an array.
* Each element in the array returned from GET /api/users is an object literal containing a user's username and _id.
* You can POST to /api/users/:_id/exercises with form data description, duration, and optionally date. If no date is supplied, the current date will be used.
* The response returned from POST /api/users/:_id/exercises will be the user object with the exercise fields added.
* You can make a GET request to /api/users/:_id/logs to retrieve a full exercise log of any user.
* A request to a user's log GET /api/users/:_id/logs returns a user object with a count property representing the number of exercises that belong to that user.
* A GET request to /api/users/:_id/logs will return the user object with a log array of all the exercises added.
* Each item in the log array that is returned from GET /api/users/:_id/logs is an object that should have a description, duration, and date properties.
* The description property of any object in the log array that is returned from GET /api/users/:_id/logs should be a string.
* The duration property of any object in the log array that is returned from GET /api/users/:_id/logs should be a number.
* The date property of any object in the log array that is returned from GET /api/users/:_id/logs should be a string. Use the dateString format of the Date API.
* You can add from, to and limit parameters to a GET /api/users/:_id/logs request to retrieve part of the log of any user. from and to are dates in yyyy-mm-dd format. limit is an integer of how many logs to send back.


## REST API Response Format

Your responses should have the following structures.

```
GET http://localhost:8080/
send index.html

users
GET http://localhost:8080/api/users
[
  {
    "_id": "63a06d1c3ac6553ce35b44fd",
    "username": "Mario Reiley"
  },
  {
    "_id": "63aa3ffca2155dcd6e4b7c65",
    "username": "Zohar de Reiley"
  },
  {
    "_id": "63aa404ba2155dcd6e4b7c66",
    "username": "Alan Daniel Reiley"
  }
]
Exercise:

POST http://localhost:8080/api/users/{id}/exercises"
{
  username: "fcc_test",
  description: "test",
  duration: 60,
  date: "Mon Jan 01 1990",
  _id: "5fb5853f734231456ccb3b05"
}

POST http://localhost:8080/api/users
User:
{
  username: "fcc_test",
  _id: "5fb5853f734231456ccb3b05"
}

GET http://localhost:8080/api/users/*
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
```
