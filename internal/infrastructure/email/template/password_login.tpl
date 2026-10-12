{{ template "header" . }}
<div class="container">
  <div class="header"><h2>{{.AppName}}</h2></div>
  <div class="content">
    <h1>ログインの確認</h1>
    <p>ログインを開始したブラウザで、以下のリンクを開いて本人確認を完了してください。</p>
    <a href="{{.VerificationURL}}" class="button">ログインを確認する</a>
    <p>リンクを開けない場合は、ログインを開始したブラウザの確認画面に以下の6桁の確認コードを入力してください。</p>
    <p style="font-size: 28px; font-weight: bold; letter-spacing: 0.2em;">{{.PIN}}</p>
    <p>確認コードの誤入力は5回までです。リンクまたは確認コードのどちらかで確認すると、両方とも使えなくなります。</p>
    <p class="expiry-notice">リンクと確認コードは{{.ExpiryMinutes}}分で期限切れになり、一度だけ使えます。再送すると以前のリンクと確認コードは使えなくなります。</p>
    <p>心当たりがない場合は、リンクを開かずこのメールを無視してください。</p>
  </div>
</div>
{{ template "footer" . }}
