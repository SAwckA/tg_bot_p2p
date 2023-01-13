package repository

import (
	"strconv"
	"time"
)

const (
	sessionMenu = "menu"
)

type Day struct {
	Date    string
	Circles []Circle
	DayPlus float64
}

func NewDay() *Day {
	t := time.Now()
	date := t.Format("02-01-2006")

	return &Day{Date: date}
}

func (d *Day) AddCircle(circle *Circle) {
	d.Circles = append(d.Circles, *circle)
}

type Circle struct {
	Name         string
	Enter        float64
	Exit         float64
	TetherSUM    float64
	Balance      float64
	CurrencyType string
	Profit       float64
}

type UserSessionStorage struct {
	storage []*UserSession
}

type UserSession struct {
	ChatID   int
	State    string
	DialogID int
	Circle   *CircleForm
}

type CircleForm struct {
	Name         string
	Enter        string
	TetherSum    string
	Exit         string
	Balance      string
	CurrencyType string
	Profit       string
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
	return &CircleForm{CurrencyType: "rub"}
}

func (u *UserSession) SetState(s string) {
	u.State = s
}

func (u *UserSession) SetDialogID(id int) {
	u.DialogID = id
}

func NewUserSession(chatID int) *UserSession {
	return &UserSession{ChatID: chatID, State: sessionMenu, Circle: NewCircleForm()}
}

func NewUserSessionStorage() *UserSessionStorage {
	return &UserSessionStorage{}
}

func (s *UserSessionStorage) CreateUserSession(chatID int) *UserSession {
	newUser := NewUserSession(chatID)
	s.storage = append(s.storage, newUser)

	return s.GetUserSession(chatID)
}

func (s *UserSessionStorage) GetUserSession(chatID int) *UserSession {

	for _, userSession := range s.storage {
		if userSession.ChatID == chatID {
			return userSession
		}
	}

	return s.CreateUserSession(chatID)
}

func (s *UserSessionStorage) UpdateDialogMsgID(chatID int, newID int) error {

	userSession := s.GetUserSession(chatID)

	if userSession == nil {
		return nil
	}

	userSession.DialogID = newID

	return nil
}
