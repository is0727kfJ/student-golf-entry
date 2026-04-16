package handler

import (
	"encoding/json"
	"net/http"

	"github.com/is0727kfJ/student-golf-entry/internal/models"
	"github.com/is0727kfJ/student-golf-entry/internal/usecase"
)

type TournamentHandler struct {
	usecase usecase.ITournamentUsecase // 脳みそ（インターフェース）を保持する
}

func NewTournamentHandler(u usecase.ITournamentUsecase) *TournamentHandler {
	return &TournamentHandler{usecase: u}
}

func (h *TournamentHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"status": "ok", "message": "APIは正常に稼働しています"}`))
}

func (h *TournamentHandler) GetTournaments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	tournaments, err := h.usecase.GetTournaments()
	if err != nil {
		http.Error(w, `{"error": "データの取得に失敗しました"}`, http.StatusInternalServerError)
		return
	}

	if tournaments == nil {
		tournaments = []models.Tournament{}
	}
	json.NewEncoder(w).Encode(tournaments)
}

func (h *TournamentHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var req models.EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "リクエスト形式エラー"}`, http.StatusBadRequest)
		return
	}

	err := h.usecase.ApplyEntry(req)
	if err != nil {
		switch err.Error() {
		case "invalid_input":
			http.Error(w, `{"error": "IDが入力されていません"}`, http.StatusBadRequest)
		case "capacity_full_or_not_found":
			http.Error(w, `{"error": "定員オーバー、または大会が存在しません"}`, http.StatusConflict)
		case "duplicate_entry":
			http.Error(w, `{"error": "既にエントリー済みです"}`, http.StatusBadRequest)
		default:
			http.Error(w, `{"error": "サーバーエラー"}`, http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status": "success", "message": "エントリー完了"}`))
}
