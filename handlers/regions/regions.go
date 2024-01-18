package regions

import (
	"net/http"
	"strings"

	apiresponses "github.com/efuchsman/distilleries_of_scotland/internal/apiresponses"
	"github.com/efuchsman/distilleries_of_scotland/internal/distilleries"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	dis distilleries.Client
}

func NewHandler(dis distilleries.Client) *Handler {
	return &Handler{
		dis: dis,
	}
}

func (h *Handler) GetRegions(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /regions request")
	regions, err := h.dis.GetRegions()
	if err != nil {
		log.Errorf("%+v", err)
		apiresponses.InternalError500(w, "region", err)
		return
	}

	apiresponses.OK200(w, regions)
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

	log.Printf("Region Data: %+v", region)
	apiresponses.OK200(w, region)
}
