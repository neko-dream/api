{{ template "header" . }}
<div class="container">
  <div class="header"><h2>{{.AppName}}</h2></div>
  <div class="content">
    <h1>ログインの確認</h1>
    <p>ログインを開始したブラウザで、下のボタンを押してください。</p>
    <a href="{{.VerificationURL}}" class="button">確認する</a>
    <p>ボタンで確認できない場合は、ログイン画面に以下の6桁の番号を入力してください。</p>
    <p style="font-size: 28px; font-weight: bold; letter-spacing: 0.2em; text-align: center;">{{.PIN}}</p>
    <p>番号を5回間違えると、このメールでは確認できなくなります。</p>
    <p class="expiry-notice">このメールでの確認は{{.ExpiryMinutes}}分以内に、一度だけ行えます。再送すると以前のメールでは確認できなくなります。</p>
    <p>心当たりがない場合は、ボタンを押さず、このメールを無視してください。</p>
  </div>
</div>
{{ template "footer" . }}
