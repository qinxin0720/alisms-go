#阿里云短信平台SDK Go语言实现

##DEMO

```go
package main

import (
	"github.com/qinxin0720/alisms-go/SmsClient"
)

const (
	accessKeyId     = "yourAccessKeyId"
	secretAccessKey = "yourAccessKeySecret"
)

func main() {
	sc, err := SmsClient.NewSMSClient(accessKeyId, secretAccessKey)
	if err != nil {
		return
	}
	statusCode := sc.SendSMS(SmsClient.Params{"1500000000", "阿里云短信", "SMS_000000", `{"code":"12345“}`})
	if statusCode == http.StatusOK {
		log.Println("发送成功")
	} else {
		log.Println("发送失败")
	}
}
```

