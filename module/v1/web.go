package v1

import (
	"Predator/pkg/config"
	"Predator/pkg/utils"
	"Predator/pkg/xhttp"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/twmb/murmur3"
	"golang.org/x/text/encoding/simplifiedchinese"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

type CheckDatas struct {
	Body     []byte
	Headers  string
	Cookie   string
	Title    string
	IconHash string
}

func WebTitle(info *config.HostInfo) error {
	if config.ScanType == "webscan" {
		WebScan(info)
		return nil
	}
	err, CheckData := GOWebTitle(info)
	info.InfoStr = InfoCheck(info.Url, CheckData)

	if !config.NoPoc && err == nil {
		WebScan(info)
	} else {
		errLog := fmt.Sprintf("[-] webtitle %v %v", info.Url, err)
		utils.LogError(errLog)
	}
	return err
}

func GOWebTitle(info *config.HostInfo) (err error, datas []CheckDatas) {
	if info.Url == "" {
		switch info.Ports {
		case "80":
			info.Url = fmt.Sprintf("http://%s", info.Host)
		case "443":
			info.Url = fmt.Sprintf("https://%s", info.Host)
		default:
			host := fmt.Sprintf("%s:%s", info.Host, info.Ports)
			protocol := getProtocol(host, config.Timeout)
			info.Url = fmt.Sprintf("%s://%s:%s", protocol, info.Host, info.Ports)
		}
	} else {
		if !strings.Contains(info.Url, "://") {
			host := strings.Split(info.Url, "/")[0]
			protocol := getProtocol(host, config.Timeout)
			info.Url = fmt.Sprintf("%s://%s", protocol, info.Url)
		}
	}

	err, result, datas := geturl(info, 1, datas)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return
	}

	if strings.Contains(result, "://") {
		info.Url = result
		err, result, datas = geturl(info, 3, datas)
		if err != nil {
			return
		}
	}

	if result == "https" && !strings.HasPrefix(info.Url, "https://") {
		info.Url = strings.Replace(info.Url, "http://", "https://", 1)
		err, result, datas = geturl(info, 1, datas)
		if strings.Contains(result, "://") {
			info.Url = result
			err, _, datas = geturl(info, 3, datas)
			if err != nil {
				return
			}
		}
	}
	// 是否访问图标
	err, _, datas = geturl(info, 2, datas)
	if err != nil {
		return
	}
	return
}

func geturl(info *config.HostInfo, flag int, CheckData []CheckDatas) (error, string, []CheckDatas) {
	Url := info.Url

	if flag == 2 {
		URL, err := url.Parse(Url)
		if err == nil {
			Url = fmt.Sprintf("%s://%s/favicon.ico", URL.Scheme, URL.Host)
		} else {
			Url += "/favicon.ico"
		}
	}
	iconHash := calcIconHash(Url)

	req, err := xhttp.NewRequest("GET", Url, nil)
	if err != nil {
		return err, "", CheckData
	}

	req.SetHeader("User-agent", uarand.GetRandom())
	req.SetHeader("Accept", config.Accept)
	req.SetHeader("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.SetHeader("Connection", "close")
	if config.Cookie != "" {
		req.SetHeader("Cookie", config.Cookie)
	}

	var client *xhttp.Client
	if flag == 1 {
		client = xhttp.DefaultClient()
	} else {
		client = xhttp.DefaultRedirectClient()
	}

	ctx := context.Background()
	resp, err := client.Do(ctx, req)
	if err != nil {
		return err, "https", CheckData
	}

	bodyReader := bytes.NewReader(resp.GetBody())
	bodyReadCloser := ioutil.NopCloser(bodyReader)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(bodyReadCloser)

	var title string
	body, err := getRespBody(resp)
	title = getTitle(body)
	if err != nil {
		return err, "https", CheckData
	}
	var cookieStrs []string
	cookies := resp.GetCookies()
	for _, cookie := range cookies {
		cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	cookieStr := strings.Join(cookieStrs, "; ")
	CheckData = append(CheckData, CheckDatas{body, fmt.Sprintf("%s", resp.GetHeaders()), cookieStr, title, iconHash})
	var reurl string
	if flag != 2 {
		if !utf8.Valid(body) {
			body, _ = simplifiedchinese.GBK.NewDecoder().Bytes(body)
		}

		length := resp.GetHeaders().Get("Content-Length")
		if length == "" {
			length = fmt.Sprintf("%v", len(body))
		}
		redirURL, err1 := resp.GetLocation()
		if err1 == nil {
			reurl = redirURL.String()
		}
		result := fmt.Sprintf("[*] WebScan: %-25v code:%-5v len:%-5v title:%-5v", resp.Request.GetUrl(), resp.GetStatus(), length, title)
		if reurl != "" {
			result += fmt.Sprintf(" 跳转url:%s", reurl)
		}
		utils.LogSuccess(result)
	}
	if reurl != "" {
		return nil, reurl, CheckData
	}
	if resp.GetStatus() == 400 && !strings.HasPrefix(info.Url, "https") {
		return nil, "https", CheckData
	}
	return nil, "", CheckData
}

func getRespBody(oResp *xhttp.Response) ([]byte, error) {
	var body []byte
	if oResp.GetHeaders().Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(bytes.NewReader(oResp.Body))
		if err != nil {
			return nil, err
		}
		defer func(gr *gzip.Reader) {
			err := gr.Close()
			if err != nil {

			}
		}(gr)
		for {
			buf := make([]byte, 1024)
			n, err := gr.Read(buf)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if n == 0 {
				break
			}
			body = append(body, buf...)
		}
	} else {
		raw, err := io.ReadAll(bytes.NewReader(oResp.Body))
		if err != nil {
			return nil, err
		}
		body = raw
	}
	return body, nil
}

func getTitle(body []byte) (title string) {
	re := regexp.MustCompile("(?ims)<title.*?>(.*?)</title>")
	find := re.FindSubmatch(body)
	if len(find) > 1 {
		title = string(find[1])
		title = strings.TrimSpace(title)
		title = strings.Replace(title, "\n", "", -1)
		title = strings.Replace(title, "\r", "", -1)
		title = strings.Replace(title, "&nbsp;", " ", -1)
		if len(title) > 100 {
			title = title[:100]
		}
		if title == "" {
			title = "\"\"" //空格
		}
	} else {
		title = "None" //没有title
	}
	return
}

func getProtocol(host string, Timeout int64) (protocol string) {
	protocol = "http"
	//如果端口是80或443,跳过Protocol判断
	if strings.HasSuffix(host, ":80") || !strings.Contains(host, ":") {
		return
	} else if strings.HasSuffix(host, ":443") {
		protocol = "https"
		return
	}

	socksconn, err := xhttp.WrapperTcpWithTimeout("tcp", host, time.Duration(Timeout)*time.Second)
	if err != nil {
		return
	}
	conn := tls.Client(socksconn, &tls.Config{MinVersion: tls.VersionTLS10, InsecureSkipVerify: true})
	defer func() {
		if conn != nil {
			defer func() {
				if err := recover(); err != nil {
					utils.LogError(err)
				}
			}()
			err := conn.Close()
			if err != nil {
				return
			}
		}
	}()
	err = conn.SetDeadline(time.Now().Add(time.Duration(Timeout) * time.Second))
	if err != nil {
		return ""
	}
	err = conn.Handshake()
	if err == nil || strings.Contains(err.Error(), "handshake failure") {
		protocol = "https"
	}
	return protocol
}

func calcIconHash(url string) string {
	content, err := fromURLGetContent(url)

	if err != nil {
		return ""
	}

	return mmh3Hash32(standBase64(content))
}

func fromURLGetContent(requrl string) (content []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(5) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // param
		},
	}

	req, err := http.NewRequest("GET", requrl, nil) //nolint
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", uarand.GetRandom())
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func mmh3Hash32(raw []byte) string {
	var h32 hash.Hash32 = murmur3.New32()
	_, err := h32.Write(raw)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func standBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')

	return buffer.Bytes()
}

func InfoCheck(Url string, CheckData []CheckDatas) []string {
	var matched bool
	var infoname []string
	// Debug发现CheckData为0
	for _, data := range CheckData {
		for _, rule := range RuleDatas {
			if rule.Type == "body" {
				matched, _ = regexp.MatchString(rule.Rule, string(data.Body))
			} else if rule.Type == "icon_hash" {
				matched, _ = regexp.MatchString(rule.Rule, data.IconHash)
			} else if rule.Type == "title" {
				matched, _ = regexp.MatchString(rule.Rule, data.Title)
			} else if rule.Type == "cookie" {
				matched, _ = regexp.MatchString(rule.Rule, data.Cookie)
			} else {
				matched, _ = regexp.MatchString(rule.Rule, data.Headers)
			}
			if matched == true {
				infoname = utils.RemoveDuplicateOfStringArray(append(infoname, rule.Name))
				break
			}
		}
	}

	infoname = utils.RemoveDuplicateElement(infoname)

	if len(infoname) > 0 {
		result := fmt.Sprintf("[+] InfoScan: %-25v %s ", Url, infoname)
		utils.LogSuccess(result)
		return infoname
	}
	return []string{""}
}
