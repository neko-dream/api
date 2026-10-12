{{ template "header" . }}
<div class="container">
  <div class="header"><h2>{{.AppName}}</h2></div>
  <div class="content">
    <h1>ログインの本人確認</h1>
    <p>ログインを開始したブラウザで、以下のリンクを開いて本人確認を完了してください。</p>
    <a href="{{.VerificationURL}}" class="button">ログインを確認する</a>
    <p>リンクを開けない場合は、ログインを開始したブラウザの確認画面に以下の6桁のPINを入力してください。</p>
    <p style="font-size: 28px; font-weight: bold; letter-spacing: 0.2em;">{{.PIN}}</p>
    <p>PINの誤入力は5回までです。リンクまたはPINのどちらかで確認すると、両方とも使えなくなります。</p>
    <p class="expiry-notice">リンクとPINは{{.ExpiryMinutes}}分で期限切れになり、一度だけ使えます。再送すると以前のリンクとPINは使えなくなります。</p>
    <p style="word-break: break-all;">{{.VerificationURL}}</p>
    <p>心当たりがない場合は、リンクを開かずこのメールを無視してください。</p>
  </div>
</div>
{{ template "footer" . }}
