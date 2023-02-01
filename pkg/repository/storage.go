package repository

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	sessionMenu  = "menu"
	databaseName = "p2p"
)

type Day struct {
	ID      primitive.ObjectID `bson:"_id"`
	Date    string             `bson:"date"`
	Circles []Circle           `bson:"circles"`
	DayPlus float64            `bson:"day_plus"`
}

func NewMongoClient(connString string) (*mongo.Client, error) {

	//Create Client
	client, err := mongo.NewClient(options.Client().ApplyURI(connString))

	if err != nil {
		return nil, err
	}

	//Test connection
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	//Ping
	err = client.Ping(context.TODO(), nil)

	return client, err
}

func NewDay() *Day {
	t := time.Now()
	date := t.Format("02-01-2006")

	return &Day{ID: primitive.NewObjectID(), Date: date}
}

func (d *Day) AddCircle(circle *Circle) {
	circle.ID = primitive.NewObjectID()
	d.Circles = append(d.Circles, *circle)
}

type Circle struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Enter        float64            `bson:"enter"`
	Exit         float64            `bson:"exit"`
	TetherSUM    float64            `bson:"tether_sum"`
	Balance      float64            `bson:"balance"`
	CurrencyType string             `bson:"currency_type"`
	Profit       float64            `bson:"profit"`
}

type UserSessionStorage struct {
	storage *mongo.Collection
}

type UserSession struct {
	ID       primitive.ObjectID `bson:"_id"`
	ChatID   int                `bson:"chat_id"`
	State    string             `bson:"state"`
	DialogID int                `bson:"dialog_id"`
	Circle   *CircleForm        `bson:"circle"`
}

type CircleForm struct {
	ID           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Enter        string             `bson:"enter"`
	TetherSum    string             `bson:"tether_sum"`
	Exit         string             `bson:"exit"`
	Balance      string             `bson:"balace"`
	CurrencyType string             `bson:"currency_type"`
	Profit       string             `bson:"profit"`
}

func (cf *CircleForm) CalculateCircle() (*Circle, error) {

	var expectedError error

	enter, err := strconv.ParseFloat(cf.Enter, 64)
	if err != nil {
		expectedError = err
	}
	exit, err := strconv.ParseFloat(cf.Exit, 64)
	if err != nil {
		expectedError = err
	}
	balance, err := strconv.ParseFloat(cf.Balance, 64)
	if err != nil {
		expectedError = err
	}
	tetherSum, err := strconv.ParseFloat(cf.TetherSum, 64)
	if err != nil {
		expectedError = err
	}

	return &Circle{
		Name:         cf.Name,
		Enter:        enter,
		TetherSUM:    tetherSum,
		Exit:         exit,
		Balance:      balance,
		CurrencyType: cf.CurrencyType,
	}, expectedError
}

func NewCircleForm() *CircleForm {
	return &CircleForm{ID: primitive.NewObjectID(), CurrencyType: "rub"}
}

type DayStorage struct {
	storage *mongo.Collection
}

func (d *DayStorage) GetDay() *Day {

	res := d.storage.FindOne(context.TODO(), bson.M{})

	var day *Day

	fmt.Println("get day", day)

	if err := res.Decode(&day); err != nil {
		fmt.Println(err)
		return nil
	}

	return day
}

func (d *DayStorage) SaveDay(day *Day) {

	filter := bson.M{
		"_id": day.ID,
	}

	fmt.Println("save day", day)

	update := bson.M{
		"$set": day,
	}

	d.storage.UpdateOne(context.TODO(), filter, update)
}

func (d *DayStorage) CreateDay() *Day {
	day := NewDay()
	fmt.Println("create day", day)
	res, err := d.storage.InsertOne(context.TODO(), day)
	if err != nil {
		fmt.Println(res)
		return nil
	}

	return day
}

func (d *DayStorage) CloseDay() {
	day := d.GetDay()

	fmt.Println("close day", day)
	filter := bson.M{
		"_id": day.ID,
	}

	d.storage.DeleteOne(context.TODO(), filter)
}

func NewDayStorage(client *mongo.Client) *DayStorage {
	db := client.Database(databaseName)
	return &DayStorage{storage: db.Collection("circles")}
}

func (u *UserSession) SetState(s string) {
	u.State = s
}

func (u *UserSession) SetDialogID(id int) {
	u.DialogID = id
}

func NewUserSession(chatID int) *UserSession {
	return &UserSession{ID: primitive.NewObjectID(), ChatID: chatID, State: sessionMenu, Circle: NewCircleForm()}
}

func NewUserSessionStorage(client *mongo.Client) *UserSessionStorage {
	db := client.Database(databaseName)
	return &UserSessionStorage{storage: db.Collection("user_session")}
}

func (s *UserSessionStorage) SaveUserSession(session *UserSession) {
	filter := bson.M{
		"_id": session.ID,
	}
	update := bson.M{
		"$set": session,
	}
	res, err := s.storage.UpdateOne(context.TODO(), filter, update)

	fmt.Println("[SAVE USER SESSION]", res, err)
}

func (s *UserSessionStorage) CreateUserSession(chatID int) *UserSession {
	newUser := NewUserSession(chatID)

	s.storage.InsertOne(context.TODO(), newUser)

	return s.GetUserSession(chatID)
}

func (s *UserSessionStorage) GetUserSession(chatID int) *UserSession {

	filter := bson.M{
		"chat_id": chatID,
	}

	var userSessoin *UserSession

	res := s.storage.FindOne(context.TODO(), filter)

	if err := res.Decode(&userSessoin); err != nil {
		return s.CreateUserSession(chatID)
	}

	return userSessoin
}

func (s *UserSessionStorage) UpdateDialogMsgID(chatID int, newID int) error {

	userSession := s.GetUserSession(chatID)

	if userSession == nil {
		return nil
	}

	userSession.DialogID = newID

	filter := bson.M{
		"chat_id": chatID,
	}

	_, err := s.storage.UpdateOne(context.TODO(), filter, bson.M{"dialog_id": newID})

	return err
}
