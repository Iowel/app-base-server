package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
}

// Cтруктура для передачи данных по WebSocket
type WsPayload struct {
	Action      string              `json:"action"`       // действие, которое должен выполнить сервер или клиент
	Message     string              `json:"message"`      // текст сообщения, которое отправляется по WebSocket
	UserName    string              `json:"username"`     // имя пользователя который отправил сообщение
	MessageType string              `json:"message_type"` // тип сообщения (text, image)
	UserID      int                 `json:"user_id"`
	Conn        WebSocketConnection `json:"-"` // контекст соединения с конкретным пользователем
}

// Тип ответа вебсокета в json-формате, то что будем отправлять пользователю
type WsJsonResponse struct {
	Action  string `json:"action"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}

// Когда клиент подключается по WebSocket, он сначала делает обычный HTTP-запрос (с заголовком Upgrade: websocket), и вот именно Upgrader обрабатывает этот апгрейд
// websocket.Upgrader нужен, чтобы превращать HTTP-запрос в WebSocket-соединение
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // ВСЕ соединения разрешены, независимо от Origin.  По умолчанию Upgrader проверяет заголовок Origin (то есть, с какого сайта пришёл запрос). Если он не совпадает с сервером, соединение отклоняется (это защита от кросс-доменных атак)
}

// Отслеживаем каждого подключенного юзера
// таблица подключённых клиентов, где каждый клиент представлен своим WebSocket-соединением, а значение — это имя пользователя или идентификатор
var clients = make(map[WebSocketConnection]string)

// Канал для отправки на него информации когда захотим ее получить
var wsChan = make(chan WsPayload)

func (app *application) WsEndPoint(w http.ResponseWriter, r *http.Request) {
	// Обновляем соединение (когда етот обработчик будет вызван, когда из js на фронтенде придет запрос - нам нужно обновить соединение)
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}

	app.infoLog.Println(fmt.Sprintf("Client connected from %s", r.RemoteAddr))

	var response WsJsonResponse
	response.Message = "Connected to server"

	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("Error WriteJSON:", err)
		return
	}

	// на данном этапе мы подключились, далее получаем соединение
	// conn - значение соединения с вебсокетом
	conn := WebSocketConnection{Conn: ws}

	// Добавляем ето соединение в нашу мапу clients
	clients[conn] = ""

	// нужно что-то делать в фоновом режиме чтобы прослушивать вебсокет соединения, если мы получим какую-то информацию с фронта - нам нужно будет что-то с ней сделать
	go app.ListenForWS(&conn)

}

func (app *application) ListenForWS(conn *WebSocketConnection) {
	// нужно чтобы приложение не просто закрылось если что то пойдет не так а чтобы мы могли корректно восстановиться
	defer func() {
		if r := recover(); r != nil {
			app.errorLog.Println("ERROR:", fmt.Sprintf("%v", r))
		}
	}()

	// переменная для отправки данных
	var payload WsPayload

	for {
		// читаем json и передаем данные в payload
		err := conn.ReadJSON(&payload)
		if err != nil {
			// do nothing
		} else {
			payload.Conn = *conn
			wsChan <- payload // отправляем в канал
		}
	}
}

// Ожидаем ввода в канал wsChan а затем выполнит действия в зависимости от значения которое указано в ответе wsChan
func (app *application) ListenToWsChannel() {
	// нас интересует только одно действие но мы все равно будем отправлять ответ
	var response WsJsonResponse
	for {
		// получаем все что было отправлено на канал
		e := <-wsChan
		switch e.Action {
		case "deleteUser":
			response.Action = "logout"
			response.Message = "Ваша учетная запись была удалена"
			response.UserID = e.UserID
			app.broadcastToAll(response) // отправим ето всем подключенным клиентам

		default:
		}
	}
}

func (app *application) broadcastToAll(response WsJsonResponse) {
	// смотрим список подключенных клиентов
	for client := range clients {
		// отправляем сообщение каждому подключенному клиенту
		err := client.WriteJSON(response)
		if err != nil {
			app.errorLog.Printf("Websocket err on %s: %s", response.Action, err)
			_ = client.Close()
			delete(clients, client)
		}

	}

}
