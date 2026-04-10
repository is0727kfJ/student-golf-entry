package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// クライアントから送られてくるJSONデータを受け取るための「型」
type EntryRequest struct {
	TournamentID string `json:"tournament_id"`
	UserID       string `json:"user_id"`
}

func main() {
	connStr := "host=localhost port=5432 user=root password=password dbname=golf_entry sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB設定エラー: ", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("DB接続失敗: ", err)
	}
	fmt.Println("🎉 データベースへの接続に成功しました！")

	// 既存のヘルスチェックAPI
	http.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(`{"status": "ok", "message": "学生ゴルフ選手権APIは正常に稼働しています🏌️‍♂️"}`))
	})

	// 🌟 新規追加：エントリー受付API（ここがポートフォリオの核心部です！）
	http.HandleFunc("POST /api/entries", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 1. 送られてきたJSONを読み込む
		var req EntryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error": "リクエストの形式が正しくありません"}`, http.StatusBadRequest)
			return
		}

		// 2. トランザクション（一連の安全な処理）の開始
		tx, err := db.Begin()
		if err != nil {
			http.Error(w, `{"error": "サーバーエラーが発生しました"}`, http.StatusInternalServerError)
			return
		}
		// 途中でエラーが起きたら、すべてのDB操作を無かったこと（Rollback）にする設定
		defer tx.Rollback()

		// 3. 【排他制御】現在の申込数が「定員未満」の場合のみ、申込数を+1する
		// ※この1行で、同時に1000アクセス来ても絶対に定員をオーバーしないように防ぎます
		updateQuery := `
			UPDATE tournaments 
			SET current_entries = current_entries + 1 
			WHERE id = $1 AND current_entries < capacity
		`
		res, err := tx.Exec(updateQuery, req.TournamentID)
		if err != nil {
			http.Error(w, `{"error": "データベース更新エラー"}`, http.StatusInternalServerError)
			return
		}

		// 4. 更新された行数を確認する
		rowsAffected, _ := res.RowsAffected()
		if rowsAffected == 0 {
			// 更新された行が0 = 既に満員だった、あるいは大会が存在しない
			w.WriteHeader(http.StatusConflict) // 409 Conflict（競合）を返す
			w.Write([]byte(`{"error": "申し訳ありません。この大会はすでに定員に達しています。"}`))
			return
		}

		// 5. 枠が確保できたので、エントリー履歴を保存する
		insertQuery := `
			INSERT INTO entries (tournament_id, user_id, status) 
			VALUES ($1, $2, 'RESERVED')
		`
		_, err = tx.Exec(insertQuery, req.TournamentID, req.UserID)
		if err != nil {
			// ※二重申し込み（UNIQUE制約違反）などのエラーはここに入ります
			http.Error(w, `{"error": "エントリーに失敗しました。既に申し込んでいる可能性があります。"}`, http.StatusBadRequest)
			return
		}

		// 6. すべて成功したので、DBの変更を確定（Commit）する
		if err := tx.Commit(); err != nil {
			http.Error(w, `{"error": "コミットエラー"}`, http.StatusInternalServerError)
			return
		}

		// 7. 成功レスポンスを返す
		w.WriteHeader(http.StatusCreated) // 201 Created
		w.Write([]byte(`{"status": "success", "message": "エントリーが完了しました！"}`))
	})

	fmt.Println("🚀 サーバーを起動しました: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
