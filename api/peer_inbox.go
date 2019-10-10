package api

import (
	"encoding/base64"
	"github.com/evilsocket/islazy/log"
	"github.com/go-chi/chi"
	"io/ioutil"
	"net/http"
)

// /api/v1/inbox
func (api *API) PeerGetInbox(w http.ResponseWriter, r *http.Request) {
	obj, err := api.Client.Inbox()
	if err != nil {
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	JSON(w, http.StatusOK, obj)
}

// POST /api/v1/unit/<fingerprint>/inbox
func (api *API) PeerSendMessageTo(w http.ResponseWriter, r *http.Request) {
	messageBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error("error reading request body: %v", err)
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	fingerprint := chi.URLParam(r, "fingerprint")

	log.Info("signing new message for %s ...", fingerprint)

	signature, err := api.Keys.SignMessage(messageBody)
	if err != nil {
		log.Error("%v", err)
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	msg := Message{
		Signature: base64.StdEncoding.EncodeToString(signature),
		Data:      base64.StdEncoding.EncodeToString(messageBody),
	}

	log.Debug("%v", msg)

	if err := api.Client.SendMessageTo(fingerprint, msg); err != nil {
		log.Error("%v", err)
		ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
	})
}