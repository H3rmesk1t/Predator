package ymlpoc

import (
	"Predator/module/v1/ymlpoc/cel"
	"Predator/module/v1/ymlpoc/lib"
	yml_structs "Predator/module/v1/ymlpoc/structs"
	"Predator/pkg/config"
	"Predator/pkg/xhttp"
	"bufio"
	"fmt"
	"github.com/google/cel-go/checker/decls"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	BodyBufPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 1024)
		},
	}
	BodyPool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}
	VariableMapPool = sync.Pool{
		New: func() interface{} {
			return make(map[string]interface{})
		},
	}
)

type RequestFuncType func(ruleName string, rule yml_structs.Rule) error

func ExecutePoc(oReq *xhttp.Request, target string, poc *yml_structs.Poc) (isVul bool, err error) {
	isVul = false

	var (
		milliseconds int64
		tcpudpType   = ""

		request       *xhttp.Request
		response      *xhttp.Response
		oProtoRequest *yml_structs.Request
		protoRequest  *yml_structs.Request
		protoResponse *yml_structs.Response
		variableMap   = VariableMapPool.Get().(map[string]interface{})
		requestFunc   cel.RequestFuncType
	)

	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(r.(error), "Run Poc[%s] error", poc.Name)
			isVul = false
		}
	}()
	// 回收
	defer func() {
		if protoRequest != nil {
			lib.PutUrlType(protoRequest.Url)
			lib.PutRequest(protoRequest)

		}
		if oProtoRequest != nil {
			lib.PutUrlType(oProtoRequest.Url)
			lib.PutRequest(oProtoRequest)

		}
		if protoResponse != nil {
			lib.PutUrlType(protoResponse.Url)
			if protoResponse.Conn != nil {
				lib.PutAddrType(protoResponse.Conn.Source)
				lib.PutAddrType(protoResponse.Conn.Destination)
				lib.PutConnectInfo(protoResponse.Conn)
			}
			lib.PutResponse(protoResponse)
		}

		for _, v := range variableMap {
			switch v.(type) {
			case *yml_structs.Reverse:
				cel.PutReverse(v)
			default:
			}
		}
		VariableMapPool.Put(variableMap)
	}()

	// 初始赋值, 设置原始请求变量
	oProtoRequest, _ = lib.ParseHttpRequest(oReq)
	variableMap["request"] = oProtoRequest

	// 判断transport，如果不合法则跳过
	transport := poc.Transport
	if transport == "tcp" || transport == "udp" {
		if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
			return
		}
	} else {
		_, err = url.ParseRequestURI(target)
		if err != nil {
			return
		}
	}

	// 初始化cel-go环境，并在函数返回时回收
	c := cel.NewEnvOption()
	defer cel.PutCustomLib(c)

	env, err := cel.NewEnv(&c)
	if err != nil {
		return false, err
	}

	// 定义渲染函数
	render := func(v string) string {
		for k1, v1 := range variableMap {
			_, isMap := v1.(map[string]string)
			if isMap {
				continue
			}
			v1Value := fmt.Sprintf("%v", v1)
			t := "{{" + k1 + "}}"
			if !strings.Contains(v, t) {
				continue
			}
			v = strings.ReplaceAll(v, t, v1Value)
		}
		return v
	}
	ReCreateEnv := func() error {
		env, err = cel.NewEnv(&c)
		if err != nil {
			return err
		}
		return nil
	}

	// 定义evaluateUpdateVariableMap
	evaluateUpdateVariableMap := func(set yaml.MapSlice) {
		for _, item := range set {
			k, expression := item.Key.(string), item.Value.(string)

			if expression == "newReverse()" && config.DnsLog {
				variableMap[k] = cel.YmlNewReverse()
				continue
			}

			// 需要重新生成一遍环境，否则之前增加的变量定义不生效
			if err := ReCreateEnv(); err != nil {
			}

			out, err := cel.Evaluate(env, expression, variableMap)
			if err != nil {
				continue
			}

			// 设置variableMap并且更新CompileOption
			switch value := out.Value().(type) {
			case *yml_structs.UrlType:
				variableMap[k] = cel.UrlTypeToString(value)
				c.UpdateCompileOption(k, cel.UrlTypeType)
			case *yml_structs.Reverse:
				variableMap[k] = value
				c.UpdateCompileOption(k, cel.ReverseType)
			case int64:
				variableMap[k] = int(value)
				c.UpdateCompileOption(k, decls.Int)
			case map[string]string:
				variableMap[k] = value
				c.UpdateCompileOption(k, cel.StrStrMapType)
			default:
				variableMap[k] = value
				c.UpdateCompileOption(k, decls.String)
			}
		}
		if err := ReCreateEnv(); err != nil {
		}
	}

	// 处理set
	evaluateUpdateVariableMap(poc.Set)

	// 处理payload
	for _, setMapVal := range poc.Payloads.Payloads {
		setMap := setMapVal.Value.(yaml.MapSlice)
		evaluateUpdateVariableMap(setMap)
	}

	// 处理transport: http时情况
	HttpRequestInvoke := func(rule yml_structs.Rule) error {
		if rule.Request.Raw != "" {
			req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rule.Request.Raw)))
			if err != nil {
				return fmt.Errorf("%s raw request format error", err)
			}
			rule.Request.Method = req.Method
			rule.Request.Path = req.RequestURI
			rule.Request.Headers = make(map[string]string)

			for k := range req.Header {
				rule.Request.Headers[k] = req.Header.Get(k)
			}
			// 拿到 body
			i := strings.Index(rule.Request.Raw, "\n\n")
			n := 2
			if i < 0 {
				i = strings.Index(rule.Request.Raw, "\r\n\r\n")
				n = 4
			}
			rule.Request.Body = rule.Request.Raw[i+n:]
		}

		var (
			ok               bool
			err              error
			ruleReq          = rule.Request
			rawHeaderBuilder strings.Builder
		)

		// 渲染请求头，请求路径和请求体
		for k, v := range ruleReq.Headers {
			ruleReq.Headers[k] = render(v)
		}
		ruleReq.Path = render(strings.TrimSpace(ruleReq.Path))
		ruleReq.Body = render(strings.TrimSpace(ruleReq.Body))

		// 尝试获取缓存
		if request, protoRequest, protoResponse, ok = lib.GetHttpRequestCache(&ruleReq); !ok || !rule.Request.Cache {
			// 获取protoRequest
			protoRequest, err = lib.ParseHttpRequest(oReq)
			if err != nil {
				return err
			}

			// 处理Path
			if strings.HasPrefix(ruleReq.Path, "/") {
				protoRequest.Url.Path = strings.Trim(oReq.GetUrl().Path, "/") + "/" + ruleReq.Path[1:]
			} else if strings.HasPrefix(ruleReq.Path, "^") {
				protoRequest.Url.Path = "/" + ruleReq.Path[1:]
			}

			if !strings.HasPrefix(protoRequest.Url.Path, "/") {
				protoRequest.Url.Path = "/" + protoRequest.Url.Path
			}

			// 某些poc没有区分path和query，需要处理
			protoRequest.Url.Path = strings.ReplaceAll(protoRequest.Url.Path, " ", "%20")
			protoRequest.Url.Path = strings.ReplaceAll(protoRequest.Url.Path, "+", "%20")

			// 克隆请求对象
			request, err = xhttp.NewRequest(ruleReq.Method, fmt.Sprintf("%s://%s%s", protoRequest.Url.Scheme, protoRequest.Url.Host, protoRequest.Url.Path), strings.NewReader(ruleReq.Body))
			if err != nil {
				return err
			}

			// 处理请求头
			request.RawRequest.Header = oReq.RawRequest.Header.Clone()
			for k, v := range ruleReq.Headers {
				request.SetHeader(k, v)
				rawHeaderBuilder.WriteString(k)
				rawHeaderBuilder.WriteString(": ")
				rawHeaderBuilder.WriteString(v)
				rawHeaderBuilder.WriteString("\n")
			}

			protoRequest.RawHeader = []byte(strings.Trim(rawHeaderBuilder.String(), "\n"))

			// 额外处理protoRequest.Raw
			protoRequest.Raw, _ = httputil.DumpRequestOut(request.RawRequest, true)

			// 发起请求
			response, milliseconds, err = lib.DoRequest(request, ruleReq.FollowRedirects)
			if err != nil {
				return err
			}

			// 获取protoResponse
			protoResponse, err = lib.ParseHttpResponse(response, milliseconds)
			if err != nil {
				return err
			}

			// 设置缓存
			lib.SetHttpRequestCache(&ruleReq, request, protoRequest, protoResponse)

		} else {
		}

		return nil
	}

	// 处理transport: tcp/udp时情况
	TCPUDPRequestInvoke := func(rule yml_structs.Rule) error {
		if rule.Request.Raw != "" {
			req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rule.Request.Raw)))
			if err != nil {
				return fmt.Errorf("%s raw request format error", err)
			}
			rule.Request.Method = req.Method
			rule.Request.Path = req.RequestURI
			rule.Request.Headers = make(map[string]string)
			for k := range req.Header {
				rule.Request.Headers[k] = req.Header.Get(k)
			}
			// 拿到 body
			i := strings.Index(rule.Request.Raw, "\n\n")
			n := 2
			if i < 0 {
				i = strings.Index(rule.Request.Raw, "\r\n\r\n")
				n = 4
			}
			rule.Request.Body = rule.Request.Raw[i+n:]
		}

		var (
			buffer = BodyBufPool.Get().([]byte)

			content      = rule.Request.Content
			connectionID = rule.Request.ConnectionID
			conn         net.Conn
			connCache    *net.Conn
			responseRaw  []byte
			readTimeout  int

			ok  bool
			err error
		)
		defer BodyBufPool.Put(buffer)

		// 获取response缓存
		if responseRaw, protoResponse, ok = lib.GetTcpUdpResponseCache(rule.Request.Content); !ok || !rule.Request.Cache {
			responseRaw = BodyPool.Get().([]byte)
			defer BodyPool.Put(responseRaw)

			// 获取connectionID缓存
			if connCache, ok = lib.GetTcpUdpConnectionCache(connectionID); !ok {
				// 处理timeout
				readTimeout, err = strconv.Atoi(rule.Request.ReadTimeout)
				if err != nil {
					return err
				}

				// 发起连接
				conn, err = net.Dial(tcpudpType, target)
				if err != nil {
					return err
				}

				// 设置读取超时
				err := conn.SetReadDeadline(time.Now().Add(time.Duration(readTimeout) * time.Second))
				if err != nil {
					return err
				}

				// 设置连接缓存
				lib.SetTcpUdpConnectionCache(connectionID, &conn)
			} else {
				conn = *connCache
			}

			// 获取protoRequest
			protoRequest, _ = lib.ParseTCPUDPRequest([]byte(content))

			// 发送数据
			_, err = conn.Write([]byte(content))
			if err != nil {
				return err
			}

			// 接收数据
			for {
				n, err := conn.Read(buffer)
				if err != nil {
					if err == io.EOF {
					} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					} else {
						return err
					}
					break
				}
				responseRaw = append(responseRaw, buffer[:n]...)
			}

			// 获取protoResponse
			protoResponse, _ = lib.ParseTCPUDPResponse(responseRaw, &conn, tcpudpType)

			// 设置响应缓存
			lib.SetTcpUdpResponseCache(content, responseRaw, protoResponse)

		}
		return nil
	}

	// reqeusts总处理
	RequestInvoke := func(requestFunc cel.RequestFuncType, ruleName string, rule yml_structs.Rule) (bool, error) {
		var (
			flag bool
			ok   bool
			err  error
		)
		err = requestFunc(rule)
		if err != nil {
			return false, err
		}

		variableMap["request"] = protoRequest
		variableMap["response"] = protoResponse

		// 执行表达式
		out, err := cel.Evaluate(env, rule.Expression, variableMap)

		if err != nil {
			return false, err
		}

		// 判断表达式结果
		flag, ok = out.Value().(bool)
		if !ok {
			flag = false
		}

		// 处理output
		evaluateUpdateVariableMap(rule.Output)

		return flag, nil
	}

	// 判断transport类型，设置requestInvoke
	if poc.Transport == "tcp" {
		tcpudpType = "tcp"
		requestFunc = TCPUDPRequestInvoke
	} else if poc.Transport == "udp" {
		tcpudpType = "udp"
		requestFunc = TCPUDPRequestInvoke
	} else {
		requestFunc = HttpRequestInvoke
	}

	ruleSlice := poc.Rules
	// 提前定义名为ruleName的函数
	for _, ruleItem := range ruleSlice {
		c.DefineRuleFunction(requestFunc, ruleItem.Key, ruleItem.Value, RequestInvoke)
	}

	// ? 最后再生成一遍环境，否则之前增加的变量定义不生效
	if err := ReCreateEnv(); err != nil {

	}

	// 执行rule 并判断poc总体表达式结果
	successVal, err := cel.Evaluate(env, poc.Expression, variableMap)
	if err != nil {
		return false, err
	}

	isVul, ok := successVal.Value().(bool)
	if !ok {
		isVul = false
	}

	return isVul, nil
}
