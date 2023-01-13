package bot

import (
	"fmt"
	"tg-bot-p2p/pkg/repository"
)

const (
	mainMessage  = "Главное меню\n\nТекущий день: %s"
	dayMessage   = "День %s\nКругов: %d\nДневной плюс: %.2f"
	pairFormRUB  = "Форма заполнения:\n\n%s\n\nВход по: %.2f\nВыход по: %.2f\nОстаток: %.2f\n\nПрофит: %.2f₽"
	pairFormUSTD = "Форма заполнения:\n\n%s\n\nВход по: %.2f\nВыход по: %.2f\nОстаток: %.2f\n\nПрофит: %.2f₮"
)

func createMainMessage(day string) string {
	return fmt.Sprintf(mainMessage, day)
}

func createDayMessage(day string, circles int, dayPlus float64) string {
	return fmt.Sprintf(dayMessage, day, circles, dayPlus)
}

func createPairForm(c *repository.Circle, date string) string {

	msg := `
%s

Вход по: %.0f (%.2f)
Выход по: %.0f (%.2f)
Остаток: %.0f (%.2f)

Профит: %.0f
%s`

	enterRate := c.Enter / c.TetherSUM

	exitRate := c.Exit / c.TetherSUM

	balance := c.Balance * enterRate

	profit := c.Exit - c.Enter

	return fmt.Sprintf(msg,
		c.Name,
		c.Enter, enterRate,
		c.Exit, exitRate,
		c.Balance, balance,
		profit,
		date)
}

func createDayCloseMessage(d *repository.Day) string {

	msg := `День %s

Кругов сделано: %d

Суммарный профит: %.2f
`
	var dayPlus float64

	dayPlus = 0.0

	for _, c := range d.Circles {
		profit := c.Exit - c.Enter
		dayPlus = dayPlus + profit
	}

	return fmt.Sprintf(msg, d.Date, len(d.Circles), dayPlus)
}
