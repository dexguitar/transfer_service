package rest

import (
	"encoding/json"
	"net/http"

	"github.com/dexguitar/p2p_service/internal/usecase"
)

type TransferRequest struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

type Handler struct {
	p2pUsecase usecase.P2PTransferUsecase
}

func NewHandler(u usecase.P2PTransferUsecase) *Handler {
	return &Handler{p2pUsecase: u}
}

func (h *Handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.p2pUsecase.Transfer(req.Sender, req.Receiver, req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"transaction_id": id,
		"status":         "created",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
