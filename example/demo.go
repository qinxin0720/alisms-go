package main

import (
    "github.com/qinxin0720/alisms-go/SmsClient"
    "net/http"
    "log"
)

const (
    accessKeyID     = "yourAccessKeyId"
    secretAccessKey = "yourAccessKeySecret"
)

func main() {
    sc, err := SmsClient.NewSMSClient(accessKeyID, secretAccessKey)
    if err != nil {
        return
    }
    statusCode, _, _ := sc.SendSMS(SmsClient.Params{"1500000000", "阿里云短信", "SMS_000000", `{"code":"12345“}`})
    if statusCode == http.StatusOK {
        log.Println("发送成功")
    } else {
        log.Println("发送失败")
    }
}
