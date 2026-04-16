package models

// クライアントから送られてくるエントリー要求のJSON型
type EntryRequest struct {
	TournamentID string `json:"tournament_id"`
	UserID       string `json:"user_id"`
}

// 大会情報を返すためのJSON型
type Tournament struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	Capacity       int    `json:"capacity"`
	CurrentEntries int    `json:"current_entries"`
	Status         string `json:"status"`
}
