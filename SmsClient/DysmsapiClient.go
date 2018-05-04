package SmsClient

import (
    "crypto/hmac"
    "crypto/md5"
    "crypto/sha1"
    "encoding/base64"
    "errors"
    "fmt"
    "io"
    "math"
    "math/rand"
    "net/http"
    "net/url"
    "os"
    "reflect"
    "strconv"
    "strings"
    "time"
)

type keepAliveAgent struct {
    defaultPort    int
    KeepAlive      bool
    KeepAliveMsecs int
    protocol       string
}

func newKeepAliveAgent(defaultPort int, KeepAlive bool, KeepAliveMsecs int, protocol string) *keepAliveAgent {
    return &keepAliveAgent{
        defaultPort:    defaultPort,
        KeepAlive:      KeepAlive,
        KeepAliveMsecs: KeepAliveMsecs,
        protocol:       protocol,
    }
}

const (
    apiVersion = "2017-05-25"
)

type dysmsapiClient struct {
    accessKeyId,
    secretAccessKey,
    endpoint string
}

func newDysmsapiClient(accessKeyId, secretAccessKey, endpoint string) (*dysmsapiClient, error) {
    var ep string
    if accessKeyId == "" {
        return nil, errors.New("accessKeyId is empty")
    }
    if secretAccessKey == "" {
        return nil, errors.New("secretAccessKey is empty")
    }
    if endpoint == "" {
        return nil, errors.New("endpoint is empty")
    }
    if endpoint[len(endpoint)-1] == 0x2F { //0x2F 为 ASCII /
        ep = endpoint[:len(endpoint)-1]
    } else {
        ep = endpoint
    }
    return &dysmsapiClient{
        accessKeyId:     accessKeyId,
        secretAccessKey: secretAccessKey,
        endpoint:        ep,
    }, nil
}

func (dsc *dysmsapiClient) SendSms(params Params, accessKeyId, secretAccessKey string) (int, error) {
    if params.PhoneNumbers == "" {
        return http.StatusBadRequest, errors.New("parameter \"PhoneNumbers\" is required")
    }
    if params.SignName == "" {
        return http.StatusBadRequest, errors.New("parameter \"SignName\" is required")
    }
    if params.TemplateCode == "" {
        return http.StatusBadRequest, errors.New("parameter \"TemplateCode\" is required")
    }
    if accessKeyId == "" {
        return http.StatusBadRequest, errors.New("parameter \"accessKeyId\" is required")
    }
    statusCode, err := request("SendSms", params, accessKeyId, secretAccessKey, dsc.endpoint)
    return statusCode, err
}

type mapList struct {
    l          []string
    normalized map[string]interface{}
}

func request(action string, param Params, accessKeyId, secretAccessKey, endpoint string) (int, error) {
    defaults := buildParams(accessKeyId)
    ml := mapList{make([]string, 0, 24), make(map[string]interface{})}
    ml.l = append(ml.l, "AccessKeyId")
    ml.l = append(ml.l, "Action")
    ml.l = append(ml.l, "Format")
    ml.l = append(ml.l, "PhoneNumbers")
    ml.l = append(ml.l, "SignName")
    ml.l = append(ml.l, "SignatureMethod")
    ml.l = append(ml.l, "SignatureNonce")
    ml.l = append(ml.l, "SignatureVersion")
    ml.l = append(ml.l, "TemplateCode")
    ml.l = append(ml.l, "TemplateParam")
    ml.l = append(ml.l, "Timestamp")
    ml.l = append(ml.l, "Version")

    var params struct {
        Action,
        Format,
        SignatureMethod,
        SignatureNonce,
        SignatureVersion,
        Timestamp,
        AccessKeyId,
        Version,
        PhoneNumbers,
        SignName,
        TemplateCode,
        TemplateParam string
    }
    params.Action = action
    params.Format = defaults.Format
    params.SignatureMethod = defaults.SignatureMethod
    params.SignatureNonce = defaults.SignatureNonce
    params.SignatureVersion = defaults.SignatureVersion
    params.Timestamp = defaults.Timestamp
    params.AccessKeyId = defaults.AccessKeyId
    params.Version = defaults.Version
    params.PhoneNumbers = param.PhoneNumbers
    params.SignName = param.SignName
    params.TemplateCode = param.TemplateCode
    params.TemplateParam = param.TemplateParam

    ml.normalized = normalize(params)
    for k, v := range ml.normalized {
        ml.normalized[k] = url.QueryEscape(v.(string))
    }
    canonicalized := strings.Replace(canonicalize(&ml), "+", "%20", -1)

    stringToSign := "GET&" + url.QueryEscape("/") + "&" + url.QueryEscape(canonicalized)
    key := secretAccessKey + "&"
    mac := hmac.New(sha1.New, []byte(key))
    mac.Write([]byte(stringToSign))
    signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
    ml.l = append(ml.l, "Signature")
    ml.normalized["Signature"] = url.QueryEscape(signature)
    urls := endpoint + "/?" + canonicalize(&ml)
    req, err := http.NewRequest("GET", urls, nil)
    if err != nil {
        return http.StatusBadRequest, err
    }
    q := req.URL.Query()
    req.URL.RawQuery = q.Encode()
    c := http.Client{}
    var resp *http.Response
    resp, err = c.Do(req)
    defer resp.Body.Close()
    return resp.StatusCode, err
}

type buildParam struct {
    Format,
    SignatureMethod,
    SignatureNonce,
    SignatureVersion,
    Timestamp,
    AccessKeyId,
    Version string
}

func buildParams(accessKeyId string) *buildParam {
    return &buildParam{
        "JSON", "HMAC-SHA1", makeNonce(), "1.0", timeStap(), accessKeyId, apiVersion,
    }
}

func makeNonce() string {
    var counter = 0
    var last float64
    machine, _ := os.Hostname()
    pid := os.Getpid()
    val := math.Floor(float64(rand.Float64() * 1000000000000))
    if val == last {
        counter++
    } else {
        counter = 0
    }
    last = val
    uid := machine + strconv.Itoa(pid) + strconv.FormatFloat(val, 'f', -1, 64) + strconv.Itoa(counter)
    m := md5.New()
    io.WriteString(m, uid)
    return fmt.Sprintf("%x", m.Sum(nil))
}

func timeStap() string {
    return time.Now().UTC().Format("2006-01-02T15:04:05Z")
}

func normalize(obj interface{}) map[string]interface{} {
    t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)

    var data = make(map[string]interface{})
    for i := 0; i < t.NumField(); i++ {
        data[t.Field(i).Name] = v.Field(i).Interface()
    }
    return data
}

func canonicalize(ml *mapList) string {
    var params string
    for _, v := range ml.l {
        params += v
        params += "="
        params += ml.normalized[v].(string)
        params += "&"
    }
    params = params[:len(params)-1]
    return params
}
