package handlers

import "net/http"

func (h *handlers) PingDataBase(w http.ResponseWriter, r *http.Request) {
	err := h.nodeService.Ping(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
