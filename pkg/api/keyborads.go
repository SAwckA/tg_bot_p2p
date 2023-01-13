package api

import (
	"fmt"
	"tg-bot-p2p/pkg/repository"
)

const (
	empty = "❌"
	saved = "✅"
)

func NullKeyboard() Keyboard {
	return Keyboard{}
}

func MenuKeyboard() Keyboard {
	return Keyboard{
		{
			Button{Text: "Новый день", CallbackData: "new day"},
		},
	}
}

func DayKeyboard() Keyboard {
	return Keyboard{
		{
			Button{Text: "Новый круг", CallbackData: "circle menu"},
		},
		{
			Button{Text: "Дневной плюс", CallbackData: "day plus"},
		},
		{
			Button{Text: "Закрыть день", CallbackData: "close day"},
		},
	}
}

func CircleKeyboard(form *repository.CircleForm) Keyboard {
	nameState := form.Name
	if len(nameState) < 1 {
		nameState = empty
	} else {
		nameState = saved
	}
	enterState := form.Enter
	if len(enterState) < 1 {
		enterState = empty
	} else {
		enterState = saved
	}
	tetherSumState := form.TetherSum
	if len(tetherSumState) < 1 {
		tetherSumState = empty
	} else {
		tetherSumState = saved
	}
	exitState := form.Exit
	if len(exitState) < 1 {
		exitState = empty
	} else {
		exitState = saved
	}
	balanceState := form.Balance
	if len(balanceState) < 1 {
		balanceState = empty
	} else {
		balanceState = saved
	}
	currencyState := form.CurrencyType
	if len(currencyState) < 1 {
		currencyState = empty
	} else {
		currencyState = saved
	}
	profitState := form.Profit
	if len(profitState) < 1 {
		profitState = empty
	} else {
		profitState = saved
	}
	return Keyboard{
		{
			Button{Text: fmt.Sprintf("%s Название связки", nameState), CallbackData: "pair name"},
		},
		{
			Button{Text: fmt.Sprintf("%s Сумма входа", enterState), CallbackData: "enter sum"},
		},
		{
			Button{Text: fmt.Sprintf("%s Сумма полученного тезера", tetherSumState), CallbackData: "tether sum"},
		},
		{
			Button{Text: fmt.Sprintf("%s Сумма выхода", exitState), CallbackData: "exit sum"},
		},
		{
			Button{Text: fmt.Sprintf("%s Остаток", balanceState), CallbackData: "balance"},
		},
		{
			Button{Text: "Закрыть круг", CallbackData: "close circle"},
		},
	}
}

func CancelKeyboard(state string) Keyboard {
	return Keyboard{
		{Button{Text: "Отмена", CallbackData: state}},
	}
}
