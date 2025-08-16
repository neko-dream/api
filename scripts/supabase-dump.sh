#!/bin/bash

if [ $# -lt 2 ]; then
    echo "エラー: Supabase接続文字列とターゲット環境を指定してください"
    echo ""
    echo "使い方: $0 <SUPABASE_DSN> <dev|prd> [オプション]"
    echo ""
    echo "オプション:"
    echo "  -t, --table SOURCE[:TARGET]   特定のテーブルをダンプ（複数指定可）"
    echo "  -o, --output FILE            出力ファイル名（デフォルト: supabase_to_{env}_dump_{timestamp}.sql）"
    echo "  -c, --copy                   COPY形式で出力（高速だが注意が必要）"
    echo ""
    echo "例:"
    echo "  $0 'postgresql://user:pass@host:5432/db' dev                    # 全データダンプ（INSERT形式）"
    echo "  $0 'postgresql://user:pass@host:5432/db' dev -c                 # COPY形式で高速ダンプ"
    echo "  $0 'postgresql://user:pass@host:5432/db' dev -t users          # usersテーブルのみ"
    echo "  $0 'postgresql://user:pass@host:5432/db' dev -t users:members  # usersをmembersとして"
    echo ""
    echo "注意: Supabaseのpublicスキーマから、RDSのkotohiro_{env}データベースへのダンプを想定"
    exit 1
fi

DSN="$1"
TARGET_ENV="$2"
shift 2

# ターゲット環境の検証
case "$TARGET_ENV" in
    "dev")
        TARGET_DB="kotohiro_dev"
        ;;
    "prd"|"prod")
        TARGET_DB="kotohiro_prd"
        TARGET_ENV="prd"
        ;;
    *)
        echo "エラー: 無効な環境 '$TARGET_ENV'"
        echo "有効な環境: dev, prd"
        exit 1
        ;;
esac

# デフォルト値
OUTPUT_FILE=""
declare -a TABLE_LIST  # 連想配列ではなく通常の配列を使用
declare -a TARGET_LIST
DUMP_ALL=true
USE_COPY=false

# オプション解析
while [[ $# -gt 0 ]]; do
    case $1 in
        -c|--copy)
            USE_COPY=true
            shift
            ;;
        -t|--table)
            DUMP_ALL=false
            if [[ "$2" =~ ^([^:]+):(.+)$ ]]; then
                TABLE_LIST+=("${BASH_REMATCH[1]}")
                TARGET_LIST+=("${BASH_REMATCH[2]}")
            else
                TABLE_LIST+=("$2")
                TARGET_LIST+=("$2")
            fi
            shift 2
            ;;
        -o|--output)
            OUTPUT_FILE="$2"
            shift 2
            ;;
        *)
            echo "不明なオプション: $1"
            exit 1
            ;;
    esac
done

# 出力ファイル名の設定
if [ -z "$OUTPUT_FILE" ]; then
    OUTPUT_FILE="supabase_to_${TARGET_ENV}_dump_$(date +%Y%m%d_%H%M%S).sql"
fi

echo "Supabase (public) → RDS ($TARGET_DB) データダンプ"
echo "出力ファイル: $OUTPUT_FILE"
echo ""

# ダンプ実行
if [ "$DUMP_ALL" = true ]; then
    # 全テーブルのデータダンプ
    echo "全テーブルのデータをダンプ中..."

    # pg_dumpでダンプして、public.を削除
    if [ "$USE_COPY" = true ]; then
        echo "COPY形式でダンプ中（高速）..."
        pg_dump "$DSN" \
            --data-only \
            --no-owner \
            --no-privileges \
            --disable-triggers \
            --schema='public' \
            --exclude-table-data='migrations' \
            --exclude-table-data='schema_migrations' > "$OUTPUT_FILE"
    else
        echo "INSERT形式でダンプ中（安全）..."
        pg_dump "$DSN" \
            --data-only \
            --no-owner \
            --no-privileges \
            --disable-triggers \
            --schema='public' \
            --exclude-table-data='migrations' \
            --exclude-table-data='schema_migrations' \
            --inserts \
            --column-inserts > "$OUTPUT_FILE"
    fi

else
    # 特定テーブルのダンプ
    > "$OUTPUT_FILE"  # ファイルをクリア

    echo "-- Supabase (public) to RDS ($TARGET_DB) Data Dump - $(date)" >> "$OUTPUT_FILE"
    echo "-- Source: Supabase public schema" >> "$OUTPUT_FILE"
    echo "-- Target: RDS $TARGET_DB database" >> "$OUTPUT_FILE"
    echo "-- Tables: ${TABLE_LIST[@]}" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"
    echo "\\connect $TARGET_DB" >> "$OUTPUT_FILE"
    echo "" >> "$OUTPUT_FILE"

    # 配列のインデックスでループ
    for i in "${!TABLE_LIST[@]}"; do
        source_table="${TABLE_LIST[$i]}"
        target_table="${TARGET_LIST[$i]}"
        echo "ダンプ中: $source_table → $target_table"

        if [ "$source_table" != "$target_table" ]; then
            # テーブル名変換が必要な場合（現在はINSERT形式のみサポート）
            echo "" >> "$OUTPUT_FILE"
            echo "-- Source: $source_table → Target: $target_table" >> "$OUTPUT_FILE"
            echo "" >> "$OUTPUT_FILE"

            # テーブル名変換時は常にINSERT形式を使用
            pg_dump "$DSN" \
                --data-only \
                --no-owner \
                --no-privileges \
                --disable-triggers \
                --table="public.$source_table" \
                --inserts \
                --column-inserts | \
            sed "s/INSERT INTO $source_table/INSERT INTO $target_table/g" >> "$OUTPUT_FILE"

            echo "" >> "$OUTPUT_FILE"
        else
            # 同じテーブル名の場合
            if [ "$USE_COPY" = true ]; then
                pg_dump "$DSN" \
                    --data-only \
                    --no-owner \
                    --no-privileges \
                    --disable-triggers \
                    --table="public.$source_table" >> "$OUTPUT_FILE"
            else
                pg_dump "$DSN" \
                    --data-only \
                    --no-owner \
                    --no-privileges \
                    --disable-triggers \
                    --table="public.$source_table" \
                    --inserts \
                    --column-inserts >> "$OUTPUT_FILE"
            fi
        fi
    done
fi

if [ $? -eq 0 ]; then
    echo ""
    echo "✅ ダンプ完了: $OUTPUT_FILE"
    echo ""
    echo "========================================="
    echo "RDS ($TARGET_DB) へのインポート手順:"
    echo "========================================="
    echo ""
    echo "1. まずRDSに接続:"
    echo "   ./scripts/db-connect $TARGET_ENV"
    echo ""
    echo "2. 必要に応じて既存データをクリア:"
    echo "   TRUNCATE TABLE table_name CASCADE;"
    echo ""
    echo "3. データをインポート:"
    echo "   \\i $OUTPUT_FILE"
    echo ""
    echo "または、コマンドラインから直接実行:"
    echo "   ./scripts/db-connect $TARGET_ENV -f $OUTPUT_FILE"
    echo ""
    echo "注意事項:"
    echo "  - インポート前にRDS側のテーブル構造が準備されていることを確認"
    echo "  - 外部キー制約がある場合は、適切な順序でインポート"
    echo "  - 大量データの場合は、トランザクションを分割することを検討"
else
    echo ""
    echo "❌ エラー: ダンプに失敗しました"
fi
