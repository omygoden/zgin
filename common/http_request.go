package common

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func HttpRequest(urlStr, method string, header map[string]string, param interface{}, proxyUrl string, timeOut int, args ...interface{}) (string, error) {
	_, err := url.ParseRequestURI(urlStr)
	if err != nil {
		return "", errors.New("地址格式有误：" + urlStr)
	}

	if timeOut < 0 {
		timeOut = 10
	}
	//初始化客户端参数
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxyUrl != "" {
		log.Println("代理ip:", proxyUrl)
		uri, err := url.Parse(proxyUrl)
		if err != nil {
			log.Println("proxy url err:", err.Error())
			return "", nil
		}
		transport.Proxy = http.ProxyURL(uri)
	} else {
		transport.Proxy = nil
	}

	client := http.Client{
		Timeout:   time.Duration(timeOut+5) * time.Second,
		Transport: transport,
	}

	var isRedirect = true
	if len(args) > 0 {
		isRedirect = args[0].(bool)
	}
	if isRedirect {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	var postParams io.Reader
	if param != nil && method == http.MethodPost {
		postParams = formatParamsByHeader(header, param)
	} else {
		postParams = nil
	}

	//创建请求体
	req, _ := http.NewRequest(method, urlStr, postParams)
	//设置header头
	if header != nil {
		headers := header
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(res), nil
}

//根据头部content-type不同，对参数进行不同处理
func formatParamsByHeader(header map[string]string, params interface{}) io.Reader {
	if v, ok := params.(string); ok {
		return strings.NewReader(v)
	}
	if v, ok := params.([]byte); ok {
		return bytes.NewReader(v)
	}
	contentType := strings.Split(header["Content-Type"], ";")[0]
	switch contentType {
	case "application/x-www-form-urlencoded":
		var strBuff = bytes.NewBufferString("")
		for k, v := range params.(map[string]interface{}) {
			strBuff.WriteString(fmt.Sprintf("%s=%v&", k, v))
		}
		str := strings.TrimRight(strBuff.String(), "&")
		return strings.NewReader(str)
	case "application/json":
		p, _ := json.Marshal(params)
		return bytes.NewReader(p)
	default:
		log.Println("header参数有误，参数格式化失败", header)
		return nil
	}
}
