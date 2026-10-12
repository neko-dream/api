{{ define "header" }}
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }}</title>
    <style>
        body {
            font-family: 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333333;
            background-color: #f9f9f9;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 1px solid #eeeeee;
        }
        .logo {
            max-width: 150px;
            margin-bottom: 15px;
        }
        .content {
            padding: 30px 0;
        }
        .footer {
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid #eeeeee;
            font-size: 12px;
            color: #777777;
        }
        h1 {
            color: #2b72e4;
            margin-top: 0;
            font-size: 24px;
        }
        .button {
            display: block;
            width: 200px;
            margin: 30px auto;
            padding: 12px 20px;
            background-color: #2b72e4;
            color: #ffffff !important;
            text-decoration: none;
            text-align: center;
            border-radius: 4px;
            font-weight: bold;
            font-size: 16px;
        }
        .verification-code {
            background-color: #f5f5f5;
            padding: 15px;
            border-radius: 4px;
            margin: 20px 0;
            text-align: center;
            font-family: monospace;
            font-size: 24px;
            letter-spacing: 2px;
        }
        .expiry-notice {
            font-size: 14px;
            color: #777;
            text-align: center;
            margin-top: 10px;
        }
        .help-text {
            font-size: 14px;
            color: #555;
            margin-top: 30px;
            border-top: 1px dashed #eee;
            padding-top: 20px;
        }
    </style>
</head>
<body>
{{ end }}
