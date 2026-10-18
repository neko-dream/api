package event

// Notifier イベントが保存されたことをイベントプロセッサーに知らせる。
// DB を定期的に見に行かないことで、Aurora Serverless v2 がアイドル時に自動一時停止できるようにする。
type Notifier struct {
	ch chan struct{}
}

func NewNotifier() *Notifier {
	return &Notifier{ch: make(chan struct{}, 1)}
}

// Notify 処理待ちのイベントがあることを知らせる。既に知らせている場合は何もしない
func (n *Notifier) Notify() {
	select {
	case n.ch <- struct{}{}:
	default:
	}
}

// C 通知を受け取るチャネル
func (n *Notifier) C() <-chan struct{} {
	return n.ch
}
