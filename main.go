package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/slack-go/slack"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/verifing-requests", handleInteraction)

	middleware := NewSecretsVerifierMiddleware(mux)
	log.Fatal(http.ListenAndServe(":80", middleware))
}

func handleInteraction(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[START]handleInteraction")

	api := slack.New(os.Getenv("SLACK_BOT_OAUTH_TOKEN"))

	var payload *slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch payload.Type {
	case slack.InteractionTypeMessageAction:
		api.PostMessage(payload.Channel.ID,
			slack.MsgOptionText("メッセージアクションが実行されました！", false))
	}

	fmt.Println("[END]handleInteraction")
}
