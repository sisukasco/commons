package http_utils

import (
	"log"
	"net/http"
)

type AuthConnection interface {
	GetLoggedInUserID(r *http.Request) (string, error)
}

func HandleCall(ff func(http.ResponseWriter, *http.Request) (interface{}, error),
	w http.ResponseWriter, r *http.Request) {

	response, err := ff(w, r)
	if err != nil {
		SendErrorResponse(err, w, r)
		return
	}
	if response != nil {
		err = SendJSON(w, http.StatusOK, response)
		if err != nil {
			log.Printf("Error sending response %+v", err)
		}
	}

}

func HandleCallWithUser(ff func(http.ResponseWriter, *http.Request, string) (interface{}, error),
	w http.ResponseWriter, r *http.Request, auth AuthConnection) {

	userID, err := auth.GetLoggedInUserID(r)
	if err != nil || len(userID) < 1 {
		log.Printf("HandleCallWithUser UserID is empty!")
		SendErrorResponse(UnauthorizedError("Not logged in"), w, r)
		return
	}

	response, err := ff(w, r, userID)
	if err != nil {
		SendErrorResponse(err, w, r)
		return
	}
	err = SendJSON(w, http.StatusOK, response)
	if err != nil {
		log.Printf("Error sending response %+v", err)
	}

}
