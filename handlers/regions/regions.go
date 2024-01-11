package regions

import (
	"encoding/json"
	"net/http"

	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	dis *distilleries.Client
}

func NewHandler(dis *distilleries.Client) *Handler {
	return &Handler{
		dis: dis,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /regions request")
	regions, err := h.dis.GetRegions()
	if err != nil {
		log.Errorf("%+v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(regions)
}
