package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	deleteRegulation = "/dr"
)

type RegulationUsecase interface {
	DeleteRegulation(ctx context.Context, regulationID uint64) error
}

type regulationHandler struct {
	regulationUsecase RegulationUsecase
}

func NewRegulationHandler(regulationUsecase RegulationUsecase) *regulationHandler {
	return &regulationHandler{regulationUsecase: regulationUsecase}
}

func (h *regulationHandler) DeleteRegulation(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Set headers
	w.Header().Set("Content-Type", "application/json")

	// Input and Output
	var input dto.GetFullRegulationRequestDTO

	// Get JSON request
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(err)
		return
	}
	defer r.Body.Close()

	err := h.regulationUsecase.DeleteRegulation(r.Context(), input.RegulationID)
	if err != nil {
		newErrorResponse(w, http.StatusBadRequest, err.Error())
	}
}
