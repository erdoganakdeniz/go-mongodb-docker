package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"time"
)

var client *mongo.Client
var collection *mongo.Collection

type Tweet struct {
	ID       int    `json:"_id,omitempty" bson:"_id"`
	FullText string `json:"full_text,omitempty" bson:"full_text"`
	User     struct {
		ScreenName string `json:"screen_name,omitempty" bson:"screen_name"`
	} `json:"user,omitempty" bson:"user"`
}

func GetTweetsEndpoint(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Content-Type","application/json")
	var tweets []Tweet
	ctx,_:=context.WithTimeout(context.Background(),30*time.Second)
	cursor,err:=collection.Find(ctx,bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "`+err.Error()+`"}`))
		return
	}
	if err = cursor.All(ctx, &tweets); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(w).Encode(tweets)

}
func SearchTweetsEndpoint(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type","application/json")
	queryParams:=r.URL.Query()
	var tweets []Tweet
	ctx,_:=context.WithTimeout(context.Background(),30*time.Second)
	searchStage:=bson.D{
		{"$search",bson.D{{"index","synsearch"},{"text",bson.D{{"query",queryParams.Get("q")},{"path","full_text"},{"synonyms","slang"},
			}},
		}},
	}
	cursor,err:=collection.Aggregate(ctx,mongo.Pipeline{searchStage})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message" :"`+err.Error()+`"}`))
		return
	}
	if err=cursor.All(ctx,&tweets);err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message" :"`+err.Error()+`"}`))
		return
	}
	json.NewEncoder(w).Encode(tweets)
}
func main() {
	fmt.Println("Starting an application...")
	ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second)
	defer cancel()
	client,err:=mongo.Connect(ctx,options.Client().ApplyURI("mongodb+srv://Admin:<password>@testcluster.kg1gr.mongodb.net/myFirstDatabase?retryWrites=true&w=majority"))
	defer func() {
		if err=client.Disconnect(ctx);err != nil {
			panic(err)
		}
	}()
	collection =client.Database("exampleDatabase").Collection("exampleCollection")
	router:=mux.NewRouter()
	router.HandleFunc("/tweets",GetTweetsEndpoint).Methods("GET")
	router.HandleFunc("/search",SearchTweetsEndpoint).Methods("GET")
	http.ListenAndServe(":12345",router)
}