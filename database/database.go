package database

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type TimeTable struct {
	UserID int    `bson:"user_id"`
	Name   string `bson:"name"`
	Mon    []int  `bson:"mon"`
	Tues   []int  `bson:"tues"`
	Wed    []int  `bson:"wed"`
	Thurs  []int  `bson:"thurs"`
	Fri    []int  `bson:"fri"`
	Sat    []int  `bson:"sat"`
}

type USER struct {
	UserId int    `bson:"userid"`
	Name   string `bson:"name"`
}

type Attendance struct {
	UserID       string        `bson:"user_id"`
	TimestampUTC bson.DateTime `bson:"timestamp_utc"`
	DateStr      string        `bson:"date_str"`
	Date         string        `bson:"date"`
	Status       string        `bson:"status"`
}

var IsAccepted bool = false
var Names []string
var IsPresent bool = false

func (t TimeTable) give(s string) []int {
	switch s {
	case "mon":
		return t.Mon
	case "tues":
		return t.Tues
	case "wed":
		return t.Wed
	case "thurs":
		return t.Thurs
	case "fri":
		return t.Fri
	case "sat":
		return t.Sat
	default:
		return nil
	}
}

func Time(slot string) int {
	switch slot {
	case "9:00-10:00":
		return 1
	case "10:00-10:50":
		return 2
	case "11:10-12:00":
		return 3
	case "12:00-1:00":
		return 4
	case "1:45-2:45":
		return 5
	case "2:45-3:45":
		return 6
	case "3:45-4:45":
		return 7
	default:
		return -1
	}
}

func Check(Name, Day, Period string) {
	IsAccepted = false
	// Name = strings.ToLower(Name)
	Day = strings.ToLower(Day)
	Period = strings.ToLower(Period)
	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		fmt.Println("URI is empty")
		return
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("Connection failed", err)
		return
	}

	defer client.Disconnect(context.Background())
	cxn, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Ping(cxn, nil)

	if err != nil {
		fmt.Println("Connection to mongo failed", err)
		return
	}
	fmt.Println("Mongodb connected")
	database := client.Database("time_table")
	collection := database.Collection("table")

	// user := TimeTable{
	//  9,
	// 	"Charan Tej S K",
	// 	[]int{1, 3, 7},
	// 	[]int{3, 6, 8},
	// 	[]int{1, 2, 4, 6},
	// 	[]int{2, 3, 4, 5, 6},
	// 	[]int{1, 4},
	// 	[]int{},
	// }

	// result, err := collection.InsertOne(cxn, user)
	// if err != nil {
	// 	fmt.Println("Insertion failed", err)
	// 	return
	// }
	// fmt.Println("insert id", result.InsertedID)

	// cursor, err := collection.Find(cxn, bson.M{"name": "Charan"})
	// if err != nil {
	// 	fmt.Println("Operation failed", err)
	// 	return
	// }
	// defer cursor.Close(cxn)

	// for cursor.Next(cxn) {
	// 	var list bson.M
	// 	err := cursor.Decode(&list)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}
	// 	fmt.Println(list)
	// }
	var data TimeTable
	err = collection.FindOne(cxn, bson.M{"name": Name}).Decode(&data)
	if err != nil {
		fmt.Println("Error in finding:", err)
		return
	}
	Attendence(data.Name)
	if got := Time(Period); got > 0 {
		if slices.Contains(data.give(Day), got) {
			IsAccepted = true
		}
	}
}

func DataProvider() {
	Names = make([]string, 0)

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("URI is empty")
		return
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("Client connection failed", err)
		return
	}
	defer client.Disconnect(context.Background())
	cxn, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(cxn, nil)
	if err != nil {
		fmt.Println("Connection failed", err)
		return
	}
	fmt.Println("Mongodb connected")
	database := client.Database("time_table")
	collection := database.Collection("table")

	cursor, err := collection.Find(cxn, bson.M{})
	if err != nil {
		fmt.Println("finding failed", err)
		return
	}
	defer cursor.Close(cxn)
	for cursor.Next(cxn) {
		var data TimeTable
		err := cursor.Decode(&data)
		if err != nil {
			fmt.Println(err)
			continue
		}
		Names = append(Names, data.Name)
	}
}

func Attendence(Name string) {
	IsPresent = false

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		fmt.Println("URI is empty")
		return
	}
	client, err := mongo.Connect(options.Client().ApplyURI(uri))

	if err != nil {
		fmt.Println("Client connection failed:", err)
		return
	}
	defer client.Disconnect(context.Background())

	cxn, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(cxn, nil)
	if err != nil {
		fmt.Println("Connection failed:", err)
		return
	}
	fmt.Println("MongoDB connected")
	database := client.Database(os.Getenv("MONGODB_USER"))
	collection := database.Collection(os.Getenv("MONGODB_USER_COLLECTION"))
	var person USER
	err = collection.FindOne(cxn, bson.M{"name": Name}).Decode(&person)

	if err != nil {
		fmt.Println("Error in finding user:", err)
		return
	}
	now := time.Now().UTC()
	currentDate := now.Format("2006-01-02")
	data := client.Database(os.Getenv("MONGODB_DATABASE"))
	collect := data.Collection(os.Getenv("MONGODB_ATTENDANCE_COLLECTION"))

	var attend Attendance

	err = collect.FindOne(cxn, bson.M{"user_id": strconv.Itoa(person.UserId), "date_str": currentDate}).Decode(&attend)
	if err != nil {
		fmt.Println("Error inf finding:", err)
		return
	}
	if currentDate == attend.DateStr && attend.Status == "PRESENT" {
		IsPresent = true
		return
	} else {
		fmt.Println("Teacher is absent")
	}
}
