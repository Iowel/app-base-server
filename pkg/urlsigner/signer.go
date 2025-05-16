package urlsigner

import (
	"fmt"
	"strings"
	"time"

	goalone "github.com/bwmarrin/go-alone"
)

// Параметры подписи
type Signer struct {
	Secret []byte
}

// Генерация токена (или подписи) на основе входных данных и секретного ключа
func (s *Signer) GenerateTokenFromString(data string) string {
	// Строка, к которой затем будет добавлена хеш-подпись
	var urlToSign string

	// Создается новый объект (или экземпляр подписи)
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	// Функция возвращает байтовый срез, содержащий подпись
	tokenBytes := crypt.Sign([]byte(urlToSign))
	// Преобразует срез байт с токеном в строку
	token := string(tokenBytes)

	// Возвращает полученную строку-токен для проверки целостности данных или для аутентификации запроса
	return token
}

// Проверка что ссылка по которой переходят пользователи не была изменена
// Aнализирует URL, извлекает из него подписанную часть (например, параметр hash), и проверяет,
// соответствует ли она тому, что должен генерировать сервер. Если подпись корректна, возвращается true, иначе — false
func (s *Signer) VerifyToken(token string) bool {
	// Создается новый объект (или экземпляр подписи)
	crypt := goalone.New(s.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

// Проверка срока действия токена
func (s *Signer) Expired(token string, minutesUntilExpire int) bool {
	// Создается новый объект (или экземпляр подписи)
	crypt := goalone.New(s.Secret, goalone.Timestamp)

	ts := crypt.Parse([]byte(token))
	return time.Since(ts.Timestamp) > time.Duration(minutesUntilExpire)*time.Minute
}
