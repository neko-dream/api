{{ template "header" . }}
  <div class="container">
      <div class="header">
          {{if .CompanyLogo}}
          <img src="{{.CompanyLogo}}" alt="{{.AppName}}" class="logo">
          {{else}}
          <h2>{{.AppName}}</h2>
          {{end}}
      </div>
      <div class="content">
          <h1>メールアドレスの確認</h1>

          {{if .RecipientName}}
          <p>{{.RecipientName}}様</p>
          {{else}}
          <p>こんにちは</p>
          {{end}}

          <p>{{.AppName}}にご登録いただきありがとうございます。以下のボタンをクリックして、メールアドレスの確認を完了してください。</p>

          <a href="{{.VerificationURL}}" class="button">メールアドレスを確認</a>

          <p class="expiry-notice">※このリンクは{{.ExpiryHours}}時間後に期限切れとなります</p>

          <p>もしボタンがクリックできない場合は、以下のURLをブラウザにコピー＆ペーストしてください：</p>
          <p style="word-break: break-all; font-size: 14px; color: #555;">{{.VerificationURL}}</p>

          <div class="help-text">
              <p>このメールに心当たりがない場合は、無視していただいて構いません。おそらく誰かがメールアドレスを間違えて入力したものと思われます。</p>
              <p>ご不明な点がございましたら、<a href="mailto:{{.ContactEmail}}">{{.ContactEmail}}</a>までお問い合わせください。</p>
          </div>
      </div>
{{ template "footer" . }}
