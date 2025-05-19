# gofiberlog

运行：./fiberlog
配置信息：config.yaml

请求示例：
curl --request POST \
  --url http://127.0.0.1:3000/api/v1/log \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Basic bXlsb2c6bXlsb2dhYmM=' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 181' \
  --header 'Content-Type: application/json' \
  --header 'Host: 127.0.0.1:3000' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{"p":"/tmp/order_2025-05-19.log","msg":"{\"request\":{\"name\":\"test\"},\"response\":{\"code\":10000,\"status\":\"success\",\"data\":{\"list\":[{\"name\":\"test\",\"age\":12}]}}}"}'