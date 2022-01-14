package handler

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"chatapp/chat"
	"chatapp/factory"
	"chatapp/response"
)

func HandleChat(f factory.Factory, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIdString, ok := vars["userId"]
		if !ok {
			l.Errorf("HandleChat: could not read 'targetId' from path params")
			response.Error{Error: "invalid request"}.ClientError(w)
			return
		}

		userId, err := strconv.ParseInt(userIdString, 10, 64)
		if err != nil {
			l.Errorf("HandleChat: invalid value for 'targetId'")
			response.Error{Error: "invalid request"}.ClientError(w)
			return
		}

		up := websocket.Upgrader{
			Error: func(wr http.ResponseWriter, rr *http.Request, status int, reason error) {
				response.Error{Error: reason.Error()}.ClientError(wr)
			},
		}

		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			l.Errorf("HandleChat: unable to upgrade session: %s", err)
			//response.Error{Error: "unexpected error happened"}.ServerError(w)
			return
		}

		hub := f.GetChatHub()
		client := chat.NewChatClient(hub, ws, userId)
		for _, cl := range hub.GetClients() {
			if cl.GetId() == userId {
				return
			}
		}

		go client.Reader()
		go client.Writer()
		hub.Register(client)
	}
}

func GetClients(f factory.Factory, l *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, client := range f.GetChatHub().GetClients() {
			l.Infoln("UserId: ", client.GetId())
		}
	}
}
