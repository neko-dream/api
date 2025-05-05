curl -L "http://127.0.0.1:3658/export/openapi/674502/0" -o api/target/apidog.openapi.yaml

jq '(..|objects|select(IN(paths; ["securitySchemes", "apikey-header-SessionId"])))."securitySchemes"."apikey-header-SessionId" |= {"type":"apiKey","in":"cookie","name":"SessionId"}' ./api/target/apidog.openapi.yaml > ./tmp.openapi.yaml
mv ./tmp.openapi.yaml ./api/target/apidog.openapi.yaml
rm -f ./tmp.openapi.yaml
sed -i.back -e 's/apikey-header-SessionId/ApiKeyAuth/g' ./api/target/apidog.openapi.yaml
rm -f ./api/target/apidog.openapi.yaml.back

cd manage
npm run build
cd ../

oasc -f ./api/target/apidog.openapi.yaml  -f ./api/target/manage.openapi.yaml -o ./static/oas/openapi.yaml --format yaml
ogen --package oas --target internal/presentation/oas --clean ./static/oas/openapi.yaml --convenient-errors=on

sqlc generate
oapi-codegen -config oapi.yaml ./api/analysis.openapi.json

find . -name "*.go" | grep -v "vendor/\|.git/\|_test.go" | xargs -n 1 -t otelinji -template "./internal/infrastructure/telemetry/otelinji.template" -w -filename &> /dev/null

cd admin-ui
npm run build
cd ../

