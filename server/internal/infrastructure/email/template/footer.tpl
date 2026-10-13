{{ define "footer" }}

    <div class="footer">
            <p>© {{.AppName}} {{.Year}}</p>
            <p><a href="{{.WebsiteURL}}">{{.WebsiteURL}}</a></p>
            <p>このメールは自動送信されています。返信はできませんのでご了承ください。</p>
    </div>
  </div>
</body>
</html>

{{ end }}
