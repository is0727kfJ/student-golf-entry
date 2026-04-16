package usecase

import (
	"errors"
	"testing"

	"github.com/is0727kfJ/student-golf-entry/internal/models"
)

// ==========================================
// 1. テスト用の「偽物の倉庫番（モック）」を作る
// ==========================================
type mockTournamentRepository struct {
	// テストのシナリオに合わせて「どんなエラーを返してほしいか」を設定できる変数
	mockCreateEntryTxError error
}

// ITournamentRepositoryのルールを満たすためのダミー
func (m *mockTournamentRepository) GetAll() ([]models.Tournament, error) {
	return nil, nil
}

// 設定されたエラーをそのまま返すだけの偽物メソッド
func (m *mockTournamentRepository) CreateEntryTx(tournamentID, userID string) error {
	return m.mockCreateEntryTxError
}

// ==========================================
// 2. 自動テストの本体（テーブル駆動テスト）
// ==========================================
func TestApplyEntry(t *testing.T) {
	// テストのシナリオをテーブル形式で定義する（上から順番に実行される）
	tests := []struct {
		name          string              // テスト名
		req           models.EntryRequest // 入力データ
		mockRepoError error               // 偽の倉庫番が返すエラー（設定）
		wantErr       bool                // エラーになることを期待するか？
		expectedErr   string              // 期待するエラーメッセージ
	}{
		{
			name:          "✅ 正常系：正しくエントリーできる",
			req:           models.EntryRequest{TournamentID: "t-1", UserID: "u-1"},
			mockRepoError: nil, // エラーなし（成功）
			wantErr:       false,
		},
		{
			name:          "❌ 異常系：入力IDが空っぽの場合は弾かれる",
			req:           models.EntryRequest{TournamentID: "", UserID: ""},
			mockRepoError: nil,
			wantErr:       true,
			expectedErr:   "invalid_input",
		},
		{
			name:          "❌ 異常系：定員オーバーの場合は弾かれる",
			req:           models.EntryRequest{TournamentID: "t-1", UserID: "u-1"},
			mockRepoError: errors.New("capacity_full_or_not_found"), // 倉庫番が満員を伝えてきた設定
			wantErr:       true,
			expectedErr:   "capacity_full_or_not_found",
		},
	}

	// 上で作ったシナリオ（テーブル）を上から順番に実行していく
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. テストのシナリオに合わせて「偽物の倉庫番」を作る
			mockRepo := &mockTournamentRepository{
				mockCreateEntryTxError: tt.mockRepoError,
			}

			// 2. 脳みそ（Usecase）に「偽物の倉庫番」を渡して初期化する（依存性の注入）
			u := NewTournamentUsecase(mockRepo)

			// 3. 実際にエントリー処理を実行
			err := u.ApplyEntry(tt.req)

			// 4. 結果の答え合わせ
			if (err != nil) != tt.wantErr {
				t.Errorf("期待するエラーの有無が違います。 結果 = %v, 期待 = %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.expectedErr {
				t.Errorf("エラーメッセージが違います。 結果 = %v, 期待 = %v", err.Error(), tt.expectedErr)
			}
		})
	}
}
