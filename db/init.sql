-- UUIDを自動生成するための拡張機能を有効化
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ユーザーテーブル
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 大会テーブル（current_entriesで現在の申込数を管理）
CREATE TABLE tournaments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    capacity INT NOT NULL,
    current_entries INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- エントリー履歴テーブル
CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tournament_id UUID NOT NULL REFERENCES tournaments(id),
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- 同じユーザーが同じ大会に二重で申し込めないようにする制約
    UNIQUE (tournament_id, user_id) 
);