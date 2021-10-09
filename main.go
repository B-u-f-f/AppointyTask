package main

import (
    "context"
    "fmt"
    "log"
    "strings"
    "net/http" 
    "encoding/json"

    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
    Name string
    Email string
    Password string
}

type post struct {
    UserId string
    Caption string
    Image_url string
    Posted_timestamp string
}

type userModel struct {
    _Id primitive.ObjectID
    Name string
    Email string
    Password string
}

type postModel struct {
    _Id primitive.ObjectID
    Caption string
    Image_url string
    Posted_timestamp string
}

type userPostModel struct {
    UserId primitive.ObjectID  
    PostIds []primitive.ObjectID 
}

var db *mongo.Database
var usersCollec *mongo.Collection
var postsCollec *mongo.Collection
var userPostsCollec *mongo.Collection

func mongoConnect() (* mongo.Database) {

    clientOptions := options.Client().ApplyURI("mongodb://127.0.0.1:27017")

    client, err := mongo.Connect(context.TODO(), clientOptions)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connected to MongoDB!")

    db := client.Database("goTask")

    return db
}

func init() {
    db = mongoConnect()
    usersCollec = db.Collection("users")
    postsCollec = db.Collection("posts")
    userPostsCollec = db.Collection("userPosts")    
}

func parseID(path string) (id string){
    var lastSlash = strings.LastIndex(path, "/")
    id = path[(lastSlash+1):]
    return id 
}

func getUserID(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "Invalid request"}`))

        return
    }

    var userID string = parseID(r.URL.Path)
    //user, ok := users[userID]
    id, err := primitive.ObjectIDFromHex(userID)
    
    if err == nil {
        var um userModel
        filter := bson.M {"_id": id}
        usersCollec.FindOne(context.TODO(), filter).Decode(&um)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(um)
    }else{
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Invalid user id"}`))
    }
}

func getPostID(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "Invalid request"}`))

        return
    }

    var postID string = parseID(r.URL.Path)
    id, err := primitive.ObjectIDFromHex(postID)

    if err == nil {
        var pm postModel
        filter := bson.M {"_id": id}
        postsCollec.FindOne(context.TODO(), filter).Decode(&pm)

        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(pm)
    }else{
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Invalid post id"}`))
    }
}

func getAllPosts(w http.ResponseWriter, r *http.Request) {
    if r.Method != "GET" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "Invalid request"}`))

        return
    }

    var userID string = parseID(r.URL.Path)

    id, err := primitive.ObjectIDFromHex(userID)

    if err == nil {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        
        var upm userPostModel
        filter := bson.M {"userid": id}
        userPostsCollec.FindOne(context.TODO(), filter).Decode(&upm)
        filter2 := bson.M {"_id" : bson.M {"$in": upm.PostIds} }
        
        ums, _ := postsCollec.Find(context.TODO(), filter2)
        json.NewEncoder(w).Encode(ums)
    }else{
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"message": "Invalid post id"}`))
    }
}

func createUser(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "Invalid request"}`))

        return
    }
    
    decoder := json.NewDecoder(r.Body)
    var u user
    e := decoder.Decode(&u);

        

    if e != nil { 
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Printf("%s\n", e)
        w.Write([]byte(`{"message": "failed user creation"}`))
    }else{

        result, err := usersCollec.InsertOne(context.TODO(), u)
        


        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        
        if err != nil {
            w.Write([]byte(`{"message": "failed user creation"}`))
        }else{
            
            var up userPostModel

            up.UserId = result.InsertedID.(primitive.ObjectID)
            up.PostIds = make([]primitive.ObjectID, 0)

            _, err := userPostsCollec.InsertOne(context.TODO(), up)

            if err != nil {
                w.Write([]byte(`{"message": "failed user post creation"}`))
            }else{
                w.Write([]byte(`{"message": "failed user post creation"}`))
            }
        }
    }
}

func createPost(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"message": "Invalid request"}`))

        return
    }
    
    //fmt.Printf(string(r.Body)
    decoder := json.NewDecoder(r.Body)
    var p post
    e := decoder.Decode(&p);

        

    if e != nil { 
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        fmt.Printf("%s\n", e)
        w.Write([]byte(`{"message": "failed post creation"}`))
    }else{
        
        result, err := postsCollec.InsertOne(context.TODO(), struct {Caption string; Image_url string; Posted_timestamp string}{p.Caption, p.Image_url, p.Posted_timestamp})


        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        
        if err != nil {
            w.Write([]byte(`{"message": "failed post creation"}`))
        }else{
            id, e := primitive.ObjectIDFromHex(p.UserId)
            
            if e != nil {
                w.Write([]byte(`{"message": "failed post creation"}`))
                return
            }

            query := bson.M {"_id": id}
            update := bson.M {"$push" : bson.M{"postids": result.InsertedID.(primitive.ObjectID)}}

            userPostsCollec.FindOneAndUpdate(context.TODO(), query, update)
        }

        w.Write([]byte(`{"message": "successful post creation"}`))
    }

}

func main(){


    http.HandleFunc("/users/", getUserID)
    http.HandleFunc("/users", createUser)
    http.HandleFunc("/posts", createPost)
    http.HandleFunc("/posts/", getPostID)
    http.HandleFunc("/posts/users/", getAllPosts)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
