package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Participant struct {
	Name  string             `json:"name,omitempty" bson:"name,omitempty"`
	Email string             `json:"email,omitempty" bson:"email,omitempty"`
	RSVP  string             `json:"rsvp,omitempty" bson:"rsvp,omitempty"`
}
type Meeting struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Participant  primitive.ObjectID `json:"participant,omitempty" bson:"participant,omitempty"`
	startTime  time.Time        `json:"startTime,omitempty" bson:"startTime,omitempty"`
    endTime    time.Time        `json:"endTime,omitempty" bson:"endTime,omitempty"`
    createdAt   time.Time      `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
func ScheduleMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var meeting Meeting
	_ = json.NewDecoder(request.Body).Decode(&meeting)
	collection := client.Database("Apointy").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := collection.InsertOne(ctx, meeting)
	json.NewEncoder(response).Encode(result)
}
func GetMeeting(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	id, _ := primitive.ObjectIDFromHex(request.URL.Query()["id"])
	var meeting Meeting
	collection := client.Database("Apointy").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.FindOne(ctx, Meeting{ID: id}).Decode(&meeting)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meeting)
 }
func GetMeetingsWithinTF(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	s, _ := request.URL.Query()["start"]
	e, _ := request.URL.Query()["end"]
	var meets Meeting
	collection := client.Database("Apointy").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.Find(ctx, Meeting{startTime:s,endTime:e}).Decode(&meets)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meets)
 }
func GetMeetings(response http.ResponseWriter, request *http.Request) { 
	response.Header().Set("content-type", "application/json")
	participant, _ := request.URL.Query()["participant"]
	var meets Meeting
	collection := client.Database("Apointy").Collection("meetings")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := collection.Find(ctx, Meeting{Participant.Email: participant}).Decode(&meets)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(meets)
}

func main() {
	fmt.Println("Starting the application...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)
	http.HandleFunc("/meetings", ScheduleMeeting)
	http.HandleFunc("/meeting?{id}", GetMeeting)
	http.HandleFunc("/meetings?s={st}&e={et}", GetMeetingsWithinTF)
	http.HandleFunc("/meetings?participant={email}participant={email}", GetMeetings)
	log.Fatal(http.ListenAndServe(":10000", nil))
}