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
	dayStorage     *repository.DayStorage
}

type ChatCtx struct {
	*api.Update
	session *repository.UserSession
}

func NewChatCtx(upd api.Update) *ChatCtx {
	return &ChatCtx{Update: &upd}
}

func NewBot(storage *repository.UserSessionStorage, tg *api.TelegramAPI, dayStorage *repository.DayStorage) *Bot {
	return &Bot{sessionStorage: storage, tg: tg, dayStorage: dayStorage}
}

func (b *Bot) setSessionID(chatID int) *repository.UserSession {
	userSession := b.sessionStorage.GetUserSession(chatID)

	if userSession == nil {
		userSession = b.sessionStorage.CreateUserSession(chatID)
	}

	return userSession
}

func (b *Bot) StartHandling(c chan api.Update) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("HANDLING RECOVER:", r)
	// 	}
	// }()

	for upd := range c {

		ctx := NewChatCtx(upd)

		if upd.CallbackQuery != nil {
			ctx.session = b.sessionStorage.GetUserSession(upd.CallbackQuery.Message.Chat.ID)
			go b.handleCallback(ctx)
			continue
		}

		if upd.Message != nil {
			if upd.Message.IsCommand() {
				ctx.session = b.sessionStorage.GetUserSession(upd.Message.Chat.ID)
				go b.handleCommand(ctx)
				continue
			}

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
		if day := b.dayStorage.GetDay(); day != nil {
			b.tg.SendMessage(ctx.Chat.ID, createDayMessage(day.Date, len(day.Circles), day.DayPlus), api.DayKeyboard())
			return
		}
		t := time.Now()
		msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
		b.tg.SendMessage(ctx.Chat.ID, msg, api.MenuKeyboard())

	case "send_channel":
		go b.tg.SendToChannel("asd", api.NullKeyboard())
	}
}

func (b *Bot) handleMessage(ctx *ChatCtx) {

	day := b.dayStorage.GetDay()

	if day == nil {
		b.tg.DeleteMessage(ctx.Chat.ID, ctx.MessageID)
		msg := b.tg.SendMessage(ctx.Chat.ID, "День уже закрыт", api.NullKeyboard())
		go func() {
			time.Sleep(time.Second * 2)
			b.tg.DeleteMessage(ctx.Chat.ID, msg.MessageID)
		}()
		return
	}

	switch ctx.session.State {

	//Заполнение круга
	case "enter pair name":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		ctx.session.Circle.Name = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter enter sum":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		ctx.session.Circle.Enter = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter tether sum":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		ctx.session.Circle.TetherSum = ctx.Message.Text
		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter exit sum":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		ctx.session.Circle.Exit = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "enter balance":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		ctx.session.Circle.Balance = ctx.Message.Text

		c, _ := ctx.session.Circle.CalculateCircle()
		b.tg.EditMessage(
			ctx.Message.Chat.ID,
			ctx.session.DialogID,
			createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	//Конец заполнения
	default:
		b.tg.DeleteMessage(ctx.Chat.ID, ctx.MessageID)
		return
	}

	if day != nil {
		b.dayStorage.SaveDay(day)
	}
	b.tg.DeleteMessage(ctx.Chat.ID, ctx.MessageID)
	b.sessionStorage.SaveUserSession(ctx.session)
}

func (b *Bot) handleCallback(ctx *ChatCtx) {
	fmt.Println(ctx.session.State)
	var editedMessage *api.Message

	switch ctx.CallbackQuery.Data {

	case "new day":
		day := b.dayStorage.GetDay()
		ctx.session.SetState("day menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		if day != nil {
			fmt.Println("DAY ALREADY OPEN")
			b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(day.Date, len(day.Circles), day.DayPlus), api.DayKeyboard())
			return
		}

		fmt.Println("CREATE NEW DAY")
		day = b.dayStorage.CreateDay()
		msg := fmt.Sprintf("День %s открыт", day.Date)
		b.tg.SendToChannel(msg, api.NullKeyboard())
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(day.Date, len(day.Circles), day.DayPlus), api.DayKeyboard())
		return
	}

	day := b.dayStorage.GetDay()

	if day == nil {
		t := time.Now()
		msg := createMainMessage(t.Format("02-01-2006"))
		b.tg.SendMessage(ctx.CallbackQuery.Chat.ID, msg, api.MenuKeyboard())
		return
	}

	switch ctx.CallbackQuery.Data {

	case "day menu":
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createDayMessage(day.Date, len(day.Circles), day.DayPlus), api.DayKeyboard())

	//КРУГ
	case "circle menu":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)
		c, _ := ctx.session.Circle.CalculateCircle()
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID,
			createPairForm(c, day.Date), api.CircleKeyboard(ctx.session.Circle))
	//Заполнение круга
	case "pair name":
		ctx.session.SetState("enter pair name")
		b.sessionStorage.SaveUserSession(ctx.session)
		editedMessage = b.EditByCtx(ctx, "Введите название связки", api.CancelKeyboard("circle menu"))

	case "enter sum":
		ctx.session.SetState("enter enter sum")
		b.sessionStorage.SaveUserSession(ctx.session)
		editedMessage = b.EditByCtx(ctx, "Введите сумму входа", api.CancelKeyboard("circle menu"))

	case "tether sum":
		ctx.session.SetState("enter tether sum")
		b.sessionStorage.SaveUserSession(ctx.session)
		editedMessage = b.EditByCtx(ctx, "Введите сумму полученного тезера", api.CancelKeyboard("circle menu"))

	case "exit sum":
		ctx.session.SetState("enter exit sum")
		b.sessionStorage.SaveUserSession(ctx.session)
		editedMessage = b.EditByCtx(ctx, "Введите сумму выхода", api.CancelKeyboard("circle menu"))

	case "balance":
		ctx.session.SetState("enter balance")
		b.sessionStorage.SaveUserSession(ctx.session)
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
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, createPairForm(c, day.Date),
			api.CircleKeyboard(ctx.session.Circle))

	case "close circle":
		ctx.session.SetState("circle menu")
		b.sessionStorage.SaveUserSession(ctx.session)

		c, err := ctx.session.Circle.CalculateCircle()
		if err != nil {
			editedMessage = b.EditByCtx(ctx, "Форма не заполнена или заполнена неправильно", api.CancelKeyboard("circle menu"))
			return
		}
		day.AddCircle(c)
		go b.tg.SendToChannel(fmt.Sprintf("Круг от %s\n%s", ctx.CallbackQuery.From.FirstName, createPairForm(c, day.Date)), api.NullKeyboard())
		editedMessage = b.EditByCtx(ctx, createDayMessage(day.Date, len(day.Circles), day.DayPlus), api.DayKeyboard())
		ctx.session.Circle = repository.NewCircleForm()

		///////////
	case "close day":
		t := time.Now()
		if day == nil {
			msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
			editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, msg, api.MenuKeyboard())
			return
		}
		b.tg.SendToChannel(createDayCloseMessage(day), api.NullKeyboard())
		b.dayStorage.CloseDay()
		msg := createMainMessage(fmt.Sprintf("%s", t.Format("02-01-2006")))
		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID, msg, api.MenuKeyboard())
		/////////////

	case "day plus":
		fmt.Println("day plus")

		var dayPlus float64

		dayPlus = 0.0

		for _, c := range day.Circles {
			profit := c.Exit - c.Enter
			dayPlus = dayPlus + profit
		}

		editedMessage = b.tg.EditMessage(ctx.CallbackQuery.Chat.ID, ctx.CallbackQuery.MessageID,
			fmt.Sprintf("Дневной плюс: %.2f", dayPlus), api.CancelKeyboard("day menu"))

	}

	ctx.session.DialogID = editedMessage.MessageID

	if day != nil {
		b.dayStorage.SaveDay(day)
	}
	b.sessionStorage.SaveUserSession(ctx.session)
}

func (b *Bot) EditByCtx(ctx *ChatCtx, msg string, keyboard api.Keyboard) *api.Message {
	return b.tg.EditMessage(ctx.CallbackQuery.Message.Chat.ID, ctx.CallbackQuery.Message.MessageID, msg, keyboard)
}
