### Prerequisites
1. Docker Community Edition
2. Docker Compose
3. yarn
4. golang
5. mage
6. [plugin-validator](https://github.com/grafana/plugin-validator)

### Installation
```
yarn install
yarn build
mage -v
```

for signing a new version:
```
export GRAFANA_API_KEY=<replace_with_key>
npx @grafana/toolkit plugin:sign --rootUrls http://localhost:3000,http://127.0.0.1:80
```

zip and test:
```
mv dist itay-s3-datasource
zip itay-s3-datasource-1.0.1.zip itay-s3-datasource -r
~/go/bin/plugincheck ~/<path-to-this-repo>/aws-s-3/itay-s3-datasource-1.0.1.zip
```