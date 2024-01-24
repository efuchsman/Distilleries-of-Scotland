package distilleries

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

func (h *Handler) GetRegionalDistilleries(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /regions/:region_name/distilleries request")

	vars := mux.Vars(r)
	regionName, exists := vars["region_name"]
	fields := log.Fields{"Region Name": regionName}
	if !exists || regionName == "" {
		log.WithFields(fields).Error("MISSING_ARG_REGION_NAME")
		apiresponses.BadRequest400(w, "region", "MISSING_ARG_REGION_NAME")
		return
	}
	regionName = strings.ToLower(regionName)
	_, err := h.dis.GetRegionByName(regionName)
	if err != nil {
		log.WithFields(fields).Errorf("%+v", err)
		apiresponses.NotFound404(w, "region")
		return
	}

	distilleries, err := h.dis.GetRegionalDistilleries(regionName)
	if err != nil {
		log.Errorf("%+v", err)
		apiresponses.InternalError500(w, "regional_distilleries", err)
		return
	}

	apiresponses.OK200(w, distilleries)
}
