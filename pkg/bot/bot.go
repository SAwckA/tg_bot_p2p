package bot

import (
	"fmt"
	"tg-bot-p2p/pkg/api"
	"tg-bot-p2p/pkg/repository"
	"time"
)

type Bot struct {
	sessionStorage *repository.UserSessionStorage
	tg             *api.TelegramAPI
	day            *repository.Day
}

type ChatCtx struct {
	*api.Update
	session *repository.UserSession
}

func NewChatCtx(upd api.Update) *ChatCtx {
	return &ChatCtx{Update: &upd}
}

func NewBot(storage *repository.UserSessionStorage, tg *api.TelegramAPI) *Bot {
	return &Bot{sessionStorage: storage, tg: tg}
}

func (b *Bot) setSessionID(chatID int) *repository.UserSession {
	userSession := b.sessionStorage.GetUserSession(chatID)

	if userSession == nil {
		userSession = b.sessionStorage.CreateUserSession(chatID)
	}

	return userSession
}

func (b *Bot) StartHandling(c chan api.Update) {
	for upd := range c {

		ctx := NewChatCtx(upd)

		if upd.CallbackQuery != nil {
			ctx.session = b.sessionStorage.GetUserSession(upd.CallbackQuery.Message.Chat.ID)
			go b.handleCallback(ctx)
			continue
		}

		if upd.Message.IsCommand() {
			ctx.session = b.sessionStorage.GetUserSession(upd.Message.Chat.ID)
			go b.handleCommand(ctx)
			continue
		}

		if upd.Message != nil {
			ctx.session = b.sessionStorage.GetUserSession(upd.Message.Chat.ID)
			go b.handleMessage(ctx)
			continue
		}
	}
}

func (b *Bot) handleCommand(ctx *ChatCtx) {
	fmt.Println("Command:", ctx.Message.Command())
	switch ctx.Message.Command() {
	case "start":
		t := time.Now()
		msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
		b.tg.SendMessage(ctx.Chat.ID, msg, api.MenuKeyboard())

	case "send_channel":
		go b.tg.SendToChannel("asd", api.NullKeyboard())
	}
}

func (b *Bot) handleMessage(ctx *ChatCtx) {

	if b.day == nil {
		return
	}

	switch ctx.session.State {

	//Заполнение круга
	case "enter pair name":
		ctx.session.SetState("circle menu")
		ctx.session.Circle.Name = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter enter sum":
		ctx.session.SetState("circle menu")
		ctx.session.Circle.Enter = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter tether sum":
		ctx.session.SetState("circle menu")
		ctx.session.Circle.TetherSum = ctx.Message.Text
		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter exit sum":
		ctx.session.SetState("circle menu")
		ctx.session.Circle.Exit = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter balance":
		ctx.session.SetState("circle menu")
		ctx.session.Circle.Balance = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	//Конец заполнения
	default:
		b.tg.DeleteMessage(ctx.Chat.ID, ctx.MessageID)
		return
	}

	b.tg.DeleteMessage(ctx.Chat.ID, ctx.MessageID)
}

func (b *Bot) handleCallback(ctx *ChatCtx) {
	fmt.Println(ctx.session.State)
	var editedMessage *api.Message

	switch ctx.CallbackQuery.Data {

	case "day menu":
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(b.day.Date, len(b.day.Circles), b.day.DayPlus), api.DayKeyboard())

	case "new day":
		ctx.session.SetState("day menu")
		if b.day != nil {
			fmt.Println("DAY ALREADY OPEN")
			b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(b.day.Date, len(b.day.Circles), b.day.DayPlus), api.DayKeyboard())
			return
		}
		fmt.Println("CREATE NEW DAY")
		b.day = repository.NewDay()
		msg := fmt.Sprintf("День %s открыт", b.day.Date)
		b.tg.SendToChannel(msg, api.NullKeyboard())
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(b.day.Date, len(b.day.Circles), b.day.DayPlus), api.DayKeyboard())

	//КРУГ
	case "circle menu":
		ctx.session.SetState("circle menu")
		c, _ := ctx.session.Circle.CalculateCircle()
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID,
			createPairForm(c, b.day.Date), api.CircleKeyboard(ctx.session.Circle))
	//Заполнение круга
	case "pair name":
		ctx.session.SetState("enter pair name")
		editedMessage = b.EditByCtx(ctx, "Введите название связки", api.CancelKeyboard("circle menu"))

	case "enter sum":
		ctx.session.SetState("enter enter sum")
		editedMessage = b.EditByCtx(ctx, "Введите сумму входа", api.CancelKeyboard("circle menu"))

	case "tether sum":
		ctx.session.SetState("enter tether sum")
		editedMessage = b.EditByCtx(ctx, "Введите сумму полученного тезера", api.CancelKeyboard("circle menu"))

	case "exit sum":
		ctx.session.SetState("enter exit sum")
		editedMessage = b.EditByCtx(ctx, "Введите сумму выхода", api.CancelKeyboard("circle menu"))

	case "balance":
		ctx.session.SetState("enter balance")
		editedMessage = b.EditByCtx(ctx, "Введите остаток", api.CancelKeyboard("circle menu"))

	case "balance currency":
		currencyType := ctx.session.Circle.CurrencyType
		if currencyType == "rub" {
			ctx.session.Circle.CurrencyType = "usdt"
		}
		if currencyType == "usdt" {
			ctx.session.Circle.CurrencyType = "rub"
		}
		c, _ := ctx.session.Circle.CalculateCircle()
		fmt.Println("CHANGE CURRENCY:", ctx.session.Circle.CurrencyType)
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createPairForm(c, b.day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "close circle":
		ctx.session.SetState("circle menu")

		c, err := ctx.session.Circle.CalculateCircle()
		if err != nil {
			editedMessage = b.EditByCtx(ctx, "Форма не заполнена или заполнена неправильно", api.CancelKeyboard("circle menu"))
			return
		}
		b.day.AddCircle(c)
		go b.tg.SendToChannel(fmt.Sprintf("Круг от %s\n%s", ctx.CallbackQuery.From.FirstName, createPairForm(c, b.day.Date)), api.NullKeyboard())
		editedMessage = b.EditByCtx(ctx, createDayMessage(b.day.Date, len(b.day.Circles), b.day.DayPlus), api.DayKeyboard())
		ctx.session.Circle = repository.NewCircleForm()

		///////////
	case "close day":
		t := time.Now()
		if b.day == nil {
			msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
			editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, msg, api.MenuKeyboard())
			return
		}
		b.tg.SendToChannel(createDayCloseMessage(b.day), api.NullKeyboard())
		b.day = nil
		msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, msg, api.MenuKeyboard())
		/////////////

	case "day plus":
		fmt.Println("day plus")

		var dayPlus float64

		dayPlus = 0.0

		for _, c := range b.day.Circles {
			profit := c.Exit - c.Enter
			dayPlus = dayPlus + profit
		}

		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID,
			fmt.Sprintf("Дневной плюс: %.2f", dayPlus), api.CancelKeyboard("day menu"))

	}

	ctx.session.DialogID = editedMessage.MessageID

}

func (b *Bot) EditByCtx(ctx *ChatCtx, msg string, keyboard api.Keyboard) *api.Message {
	return b.tg.EditMessage(ctx.CallbackQuery.Message.Chat.ID, ctx.CallbackQuery.Message.MessageID, msg, keyboard)
}
