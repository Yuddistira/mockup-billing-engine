package usecase

import "net/http"

func (u *Usecase) ClearHandler(w http.ResponseWriter, _ *http.Request) {
	// Send empty string to remove the div
	w.Write([]byte(""))
}
