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
          <p><strong>{{.OrganizationName}}</strong>があなたをメンバーとして招待しています。</p>

          <div class="invitation-details">
              <p>以下の情報でログインすることができます。</p>
              <ul>
                  <li><strong>Email: </strong> {{.Email}}</li>
                  <li><strong>初期パスワード：</strong> {{.Password}}</li>
              </ul>
              <p class="important-note">※初回ログイン時にパスワードの変更が必要です。</p>
          </div>

          <p>下記のURLよりログインしてください。</p>
          <p style="word-break: break-all; font-size: 14px; color: #555;">{{.InvitationURL}}</p>

          <div class="help-text">
              <p>このメールに心当たりがない場合は、無視していただいて構いません。</p>
              <p>ご不明な点がございましたら、<a href="mailto:{{.ContactEmail}}">{{.ContactEmail}}</a>までお問い合わせください。</p>
          </div>
      </div>
{{ template "footer" . }}
