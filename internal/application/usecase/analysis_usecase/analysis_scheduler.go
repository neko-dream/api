package analysis_usecase

import (
	"context"

	"github.com/neko-dream/api/internal/domain/model/shared"
	"github.com/neko-dream/api/internal/domain/model/talksession"
)

// ScheduleMode は分析リクエストのスケジュール方式を表す
type ScheduleMode int

const (
	// ScheduleDebounced はデバウンスモード。最後のリクエストから一定時間待ってから分析を実行する。
	// 投票時に使用し、連続する投票を1回の分析にまとめる。
	ScheduleDebounced ScheduleMode = iota

	// ScheduleImmediate は即時モード。同時実行制限と重複排除は適用されるが、デバウンスなしで実行する。
	// トグルや編集による分析有効化時に使用。
	ScheduleImmediate
)

// AnalysisScheduler は分析リクエストをデバウンス・重複排除・同時実行制限付きでスケジュールする
type AnalysisScheduler interface {
	// ScheduleAnalysis は指定されたトークセッションの分析をスケジュールする。
	// mode=ScheduleDebounced: デバウンス期間後に実行（期間中に新リクエストが来たらリセット）
	// mode=ScheduleImmediate: 即座にディスパッチ（同時実行制限・重複排除は適用）
	ScheduleAnalysis(ctx context.Context, talkSessionID shared.UUID[talksession.TalkSession], mode ScheduleMode)

	// Start はスケジューラーのバックグラウンド処理を開始する。アプリケーション起動時に1回呼ぶ。
	Start(ctx context.Context)

	// Stop はスケジューラーをgracefulにシャットダウンする。実行中の分析の完了を待つ。
	Stop()
}
