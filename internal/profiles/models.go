package profiles

import "time"

type Profile struct {
	ID        int       `json:"id"`      // Идентификатор профиля
	UserID    int       `json:"user_id"` // Идентификатор пользователя (внешний ключ)
	Avatar    string    `json:"avatar"`  // URL изображения профиля
	About     string    `json:"about"`   // Информация о пользователе
	Friends   []int     `json:"friends"` // Список друзей, представленных как массив строк
	Wallet    int       `json:"wallet"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"` // Время создания профиля
	UpdatedAt time.Time `json:"updated_at"` // Время последнего обновления профиля
}
