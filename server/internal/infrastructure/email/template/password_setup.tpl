{{template "header" .}}
<div class="container"><div class="content">
<h1 style="font-size:24px;color:#222;margin-bottom:24px">パスワードの設定</h1>
<p>{{.RecipientName}}さん</p>
<p>ことひろのアカウントが作成されました。下のボタンから、ご自身でパスワードを設定してください。</p>
<p style="margin:32px 0"><a href="{{.SetupURL}}" style="display:inline-block;background:#1677ff;color:#fff;text-decoration:none;padding:14px 24px;border-radius:8px">パスワードを設定する</a></p>
<p>このリンクは{{.ExpiryHours}}時間有効で、一度だけ使えます。期限切れの場合は、アカウントを作成した管理者にメールの再送を依頼してください。</p>
<p>設定後はログイン画面でメールアドレスとパスワードを入力してください。ログインのたびに、メールに届く6桁の番号で確認します。</p>
<p style="font-size:13px;color:#666">心当たりがない場合は、パスワードを設定せず、このメールを破棄してください。</p>
</div>
{{template "footer" .}}
