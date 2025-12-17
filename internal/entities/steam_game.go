package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SteamResponse — корневой объект ответа.
// JSON: { "items": [...] }
type SteamResponse struct {
	Items []SteamItem `json:"items"`
}

// Metascore - кастомный тип для обработки metascore, который может быть int или string
type Metascore string

// UnmarshalJSON обрабатывает metascore, который может быть int (0) или string
func (m *Metascore) UnmarshalJSON(data []byte) error {
	// Пробуем распарсить как число
	var num int
	if err := json.Unmarshal(data, &num); err == nil {
		// Если это 0, значит метаскор отсутствует
		if num == 0 {
			*m = ""
			return nil
		}
		*m = Metascore(strconv.Itoa(num))
		return nil
	}

	// Пробуем распарсить как строку
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		// Если пустая строка или "0", значит метаскор отсутствует
		if str == "" || str == "0" {
			*m = ""
			return nil
		}
		*m = Metascore(str)
		return nil
	}

	// Если ничего не подошло, оставляем пустым
	*m = ""
	return nil
}

// String возвращает строковое представление метаскора
func (m Metascore) String() string {
	return string(m)
}

// IsEmpty проверяет, пустой ли метаскор
func (m Metascore) IsEmpty() bool {
	return m == "" || m == "0"
}

// SteamItem — один элемент (игра/дополнение).
type SteamItem struct {
	Type              string     `json:"type"`
	Name              string     `json:"name"`
	ID                int        `json:"id"`
	Price             *PriceInfo `json:"price,omitempty"` // может отсутствовать → указатель
	TinyImage         string     `json:"tiny_image"`
	Metascore         Metascore  `json:"metascore"` // может быть int (0) или string
	Platforms         Platforms  `json:"platforms"`
	StreamingVideo    bool       `json:"streamingvideo"`
	ControllerSupport string     `json:"controller_support,omitempty"` // не у всех есть
}

// PriceInfo — информация о цене.
type PriceInfo struct {
	Currency string `json:"currency"`
	Initial  int    `json:"initial"` // в центах (999 = $9.99)
	Final    int    `json:"final"`   // в центах
}

// Platforms — поддерживаемые ОС.
type Platforms struct {
	Windows bool `json:"windows"`
	Mac     bool `json:"mac"`
	Linux   bool `json:"linux"`
}

// String возвращает человекочитаемое представление игры для Telegram.
func (s SteamItem) String() string {
	// Цена
	price := "бесплатно"
	if s.Price != nil {
		// Форматируем как $9.99 (не 999 центов!)
		price = fmt.Sprintf("%s %.2f", s.Price.Currency, float64(s.Price.Final)/100)
	}

	// Платформы (эмодзи)
	var platforms string
	if s.Platforms.Windows {
		platforms += "🖥️"
	}
	if s.Platforms.Mac {
		platforms += "🍎"
	}
	if s.Platforms.Linux {
		platforms += "🐧"
	}
	if platforms == "" {
		platforms = "—"
	}

	// Metascore (если есть)
	metascore := ""
	if !s.Metascore.IsEmpty() {
		metascore = fmt.Sprintf(" ⭐ %s", s.Metascore.String())
	}

	// Controller support (если есть)
	controller := ""
	if s.ControllerSupport != "" {
		controller = fmt.Sprintf(" 🎮 %s", s.ControllerSupport)
	}

	// Формируем строку (Markdown/HTML-friendly)
	return fmt.Sprintf(
		"🎮 *%s*\n"+
			"💰 %s\n"+
			"📊%s\n"+
			"💻 %s\n"+
			"🔗 [Store](https://store.steampowered.com/app/%d/)%s",
		s.Name,
		price,
		metascore,
		platforms,
		s.ID,
		controller,
	)
}

// RegionalPriceInfo represents price information for a specific region
type RegionalPriceInfo struct {
	CountryCode  string
	CountryFlag  string
	Item         *SteamItem
	ConvertedRub float64 // Converted price to rubles if available
}

// MultiRegionPriceData holds pricing information across multiple regions
type MultiRegionPriceData struct {
	ID       int
	GameName string
	Regions  []*RegionalPriceInfo
}
