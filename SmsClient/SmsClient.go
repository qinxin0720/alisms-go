package SmsClient

import (
    "errors"
)

const (
    DYSMSAPI_ENDPOINT = "http://dysmsapi.aliyuncs.com"
)

type Params struct {
    PhoneNumbers,
    SignName,
    TemplateCode,
    TemplateParam string
}

type SMSClient struct {
    accessKeyId,
    secretAccessKey string
    dysmsapiClient *dysmsapiClient
}

func NewSMSClient(accessKeyId, secretAccessKey string) (*SMSClient, error) {
    if accessKeyId == "" {
        return nil, errors.New("accessKeyId is empty")
    }
    if secretAccessKey == "" {
        return nil, errors.New("secretAccessKey is empty")
    }
    dsmsc, err := newDysmsapiClient(accessKeyId, secretAccessKey, DYSMSAPI_ENDPOINT)
    if err != nil {
        return nil, err
    }
    return &SMSClient{
        accessKeyId:     accessKeyId,
        secretAccessKey: secretAccessKey,
        dysmsapiClient:  dsmsc,
    }, nil
}

func (sc *SMSClient) SendSMS(params Params) int {
    statusCode, _ := sc.dysmsapiClient.SendSms(params, sc.accessKeyId, sc.secretAccessKey)
    return statusCode
}
