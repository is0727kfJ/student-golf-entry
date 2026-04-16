package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/is0727kfJ/student-golf-entry/internal/handler"
	"github.com/is0727kfJ/student-golf-entry/internal/repository"
	"github.com/is0727kfJ/student-golf-entry/internal/usecase"
	_ "github.com/lib/pq"
)

func main() {
	// 1. DB接続
	connStr := "host=localhost port=5432 user=root password=password dbname=golf_entry sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB設定エラー: ", err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		log.Fatal("DB接続失敗: ", err)
	}

	// 2. 依存関係の注入（Dependency Injection）
	repo := repository.NewTournamentRepository(db) // 倉庫番を作る
	usecase := usecase.NewTournamentUsecase(repo)  // 脳みそを作り、倉庫番を渡す
	h := handler.NewTournamentHandler(usecase)     // 窓口を作り、脳みそを渡す

	// 3. ルーティング
	http.HandleFunc("GET /api/health", h.HealthCheck)
	http.HandleFunc("GET /api/tournaments", h.GetTournaments)
	http.HandleFunc("POST /api/entries", h.CreateEntry)

	// 4. サーバー起動
	fmt.Println("🚀 レイヤードアーキテクチャでサーバーを起動しました: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
