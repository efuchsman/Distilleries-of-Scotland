package regions

import (
	"encoding/json"
	"net/http"
	"strings"

	apiresponses "github.com/efuchsman/distilleries_of_scotland/internal/apiresponses"
	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	"github.com/gorilla/mux"
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

func (h *Handler) GetRegions(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) GetRegion(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /regions/:region_name request")

	vars := mux.Vars(r)
	regionName := vars["region_name"]
	regionName = strings.ToLower(regionName)

	fields := log.Fields{"Region Name": regionName}

	region, err := h.dis.GetRegionByName(regionName)
	if err != nil {
		log.WithFields(fields).Errorf("%+v", err)
		apiresponses.NotFound404(w, "region")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(region)
}
