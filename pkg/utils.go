package pkg

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"

	color "github.com/fatih/color"
	"github.com/valyala/fasthttp"
)

var System = runtime.GOOS

func RandomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuwxyz123456789_")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomStringUpper(n int) string {
	letters := []rune("QWERTYUIOPASDFGHJKLZXCVBNM1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}
func RandomStringNumber(n int) string {
	letters := []rune("1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)

}
func Getuuid() string {
	resp, _ := http.Get("https://httpbin.org/uuid")
	body, _ := ioutil.ReadAll(resp.Body)
	body1 := string(body)
	return regexp.MustCompile("\"uuid\": \"(.*?)\"").FindStringSubmatch(body1)[1]
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func CreateFile(path, text string) os.File {
	File, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer File.Close()
	File.WriteString(text)
	return *File

}

func FasthttpHTTPDialer(proxy string) fasthttp.DialFunc {
	var auth string
	if strings.Contains(proxy, "@") {
		split := strings.Split(proxy, "@")
		auth = base64.StdEncoding.EncodeToString([]byte(split[0]))
		proxy = split[1]
	}
	return func(addr string) (net.Conn, error) {
		var conn net.Conn
		var err error
		conn, err = fasthttp.Dial(proxy)
		if err != nil {
			return nil, err
		}

		req := fmt.Sprintf("CONNECT %s HTTP/1.1\r\n", addr)
		if auth != "" {
			req += fmt.Sprintf("Proxy-Authorization: Basic %s\r\n", auth)

		}
		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)
		res.SkipBody = true
		if err := res.Read(bufio.NewReader(conn)); err != nil {
			conn.Close()
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			conn.Close()
			return nil, fmt.Errorf("could not connect to proxy: %s status code: %d", proxy, res.Header.StatusCode())
		}
		return conn, nil
	}

}

func Random_USER_AGENT() string {
	var Devices_menu = []string{"HUAWEI", "Xiaomi", "samsung", "OnePlus"}
	var DPIs = []string{"480", "320", "640", "515", "120", "160", "240", "800"}
	var randResolution int = Int(RandomStringNumber(2)) * 180
	var lowerResolution int = randResolution - 180
	var DEVICE_SETTINGS = []string{"Android", "Instagram", "" + Devices_menu[rand.Intn(len(Devices_menu))] + "", "" + Devices_menu[rand.Intn(len(Devices_menu))] + "-" + RandomStringUpper(4) + "", RandomStringNumber(4), RandomStringNumber(4), RandomString(4) + RandomStringNumber(4), "" + fmt.Sprintf("%d"+"x"+"%d", randResolution, lowerResolution) + "", RandomString(6), DPIs[rand.Intn(len(DPIs))]}
	var system string = DEVICE_SETTINGS[0]
	var Host string = DEVICE_SETTINGS[1]
	var manufacturer string = DEVICE_SETTINGS[2]
	var model string = DEVICE_SETTINGS[3]
	var android_version string = DEVICE_SETTINGS[4]
	var android_release string = DEVICE_SETTINGS[5]
	var cpu string = DEVICE_SETTINGS[6]
	var resolution string = DEVICE_SETTINGS[7]
	var randomL string = DEVICE_SETTINGS[8]
	var dpi string = DEVICE_SETTINGS[9]
	//var ss string = //, "133.0.0.34.124", system, "("+android_version+"/"+android_release+"; "+dpi+"dpi; "+resolution+"; "+manufacturer+"; "+model+"; "+cpu+";"+randomL+"; en_US)"
	//var UserAgent string =
	return string("" + Host + " 133.0.0.34.124 " + system + "(" + android_version + "/" + android_release + "; " + dpi + "dpi; " + resolution + "; " + manufacturer + "; " + model + "; " + cpu + ";" + randomL + "; en_US)")

	//"("+android_version+"/"+android_release+"; "+dpi+"dpi; "+resolution+"; "+manufacturer+"; "+model+"; "+cpu+";"+randomL+"; en_US)"
	//////////////////////////////////fmt.Sprintf("{Host} 133.0.0.34.124 {system} ({android_version}/{android_release}; {dpi}dpi; {resolution}; {manufacturer}; {model}; {cpu}; {randomL}; en_US)")

}
func CreateFileOnly(Path string) os.File {
	File, err := os.Create(Path)
	if err != nil {
		log.Fatalln(File)
	}
	return *File
}

func Request(method, Endpoint, Data string, httpHeaders map[string]string) string {
Retry:
	request, err := http.NewRequest(method, Endpoint, bytes.NewBufferString(Data))
	if err != nil {
		goto Retry
	}
	for headerType, headerValue := range httpHeaders {
		request.Header.Set(headerType, headerValue)
	}
	//fmt.Println(request)

	Response, err := http.DefaultClient.Do(request)
	if err != nil {
		goto Retry
	}
	defer Response.Body.Close()
	bodys, err := ioutil.ReadAll(Response.Body)
	if err != nil {
		goto Retry
	}
	return string(bodys)
}
func FastHttpRequest(method, Endpoint, Data string, httpHeaders map[string]string, proxies bool, proxy []string) string {

	Req := fasthttp.AcquireRequest()
	data := []byte(Data)
	Req.SetRequestURI(Endpoint)
	Req.Header.SetMethod(method)
	for headerType, headerValue := range httpHeaders {
		Req.Header.Set(headerType, headerValue)
	}

	Req.SetBody(data)
	if proxies {
		proxy := proxy[rand.Intn(len(proxy))]

		Requests = &fasthttp.Client{

			Dial:                          FasthttpHTTPDialer(proxy),
			MaxConnsPerHost:               64,
			ReadBufferSize:                1000000,
			NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			TLSConfig:                     &tls.Config{InsecureSkipVerify: true},
		}

	} else {

		Requests = &fasthttp.Client{

			//Dial:                          FasthttpHTTPDialer(proxy),
			TLSConfig:                     &tls.Config{InsecureSkipVerify: true},
			MaxConnsPerHost:               1000000,
			ReadBufferSize:                1000000,
			NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
			DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
			DisablePathNormalizing:        true,
			Dial: (&fasthttp.TCPDialer{
				Concurrency:      4096,
				DNSCacheDuration: time.Hour,
			}).Dial,
		}

	}
	Response := fasthttp.AcquireResponse()
	Response.ConnectionClose()
	defer fasthttp.ReleaseRequest(Req)
	defer fasthttp.ReleaseResponse(Response)
	Requests.Do(Req, Response)
	return string(Response.Body())

}
func FasthttpWebinsta(method, Endpoint, Endpoint2, Data, Cookie string, Proxies bool, proxy []string) string {
	var Body string
	Request1 := FastHttpRequest(method, "https://www.instagram.com/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"content-type":     "application/x-www-form-urlencoded",
		"user-agent":       "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"x-csrftoken":      "2cd27df1ee26",
		"x-requested-with": "XMLHttpRequest",
		"cookie":           Cookie}, Proxies, proxy)

	if Request1 != "" {
		Body = Request1
	}
	return Body
}
func FasthttpApiInsta(method, Endpoint, Endpoint2, Data, Cookie string, Proxies bool, proxy []string) string {
	var Body string
	Request1 := FastHttpRequest(method, "https://i.instagram.com/api/v1/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"User-Agent":      "Instagram 133.0.0.34.124 Android (29/10; 480dpi; 1080x2107; HUAWEI; MAR-LX1M; HWMAR; kirin710; en_US; 302733750)",
		"Accept":          "*/*",
		"Host":            "i.instagram.com",
		"Accept-Language": "en-Us",
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":          Cookie}, Proxies, proxy)
	if Request1 != "" {
		Body = Request1

	}
	return Body

}

func RequestWEBinsta(method, Endpoint, Endpoint2, Data, Cookie string) string {
	var Body string
	Request1 := Request(method, "https://www.instagram.com/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded",
		"origin":           "https://www.instagram.com",
		"referer":          "https://www.instagram.com/" + Endpoint + "/" + Endpoint2 + "/",
		"user-agent":       "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"x-instagram-ajax": "2cd27df1ee26",
		"x-csrftoken":      "2cd27df1ee26",
		"x-ig-app-id":      "missing",
		"x-ig-www-claim":   "0",
		"x-requested-with": "XMLHttpRequest",
		"cookie":           Cookie})
	if Request1 != "" {
		Body = Request1
	}
	return Body

	//fmt.Println(Request1)
	// Data = fmt.Sprintf("first_name=&email=%s&username=%s&phone_number=%s&biography=&chaining_enabled=on", RandomString(20)+"@gmail.com", "ss", "")
	// request, _ := http.NewRequest("POST", "https://www.instagram.com/accounts/edit/", bytes.NewBufferString(Data))
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// request.Header.Add("origin", "https://www.instagram.com")
	// request.Header.Add("referer", "https://www.instagram.com/accounts/edit/")
	// request.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	// request.Header.Add("x-instagram-ajax", "2cd27df1ee26")
	// request.Header.Add("x-csrftoken", "xvgrVcdCNWh80yZbneVgxbbR0fLEeN8X")
	// request.Header.Add("x-ig-app-id", "missing")
	// request.Header.Add("x-ig-www-claim", "0")
	// request.Header.Add("x-requested-with", "XMLHttpRequest")
	// request.Header.Set("cookie", "sessionid="+Cookie+"")
	// fmt.Println(request)
	// D, _ := http.DefaultClient.Do(request)
	// bodys, _ := ioutil.ReadAll(D.Body)
	// fmt.Println(string(bodys))

	// Retry:
	// 	request, err := http.NewRequest("POST", "https://www.instagram.com/accounts/edit/", bytes.NewBufferString(Data))
	// 	if err != nil {
	// 		goto Retry
	// 	}
	// 	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// 	request.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36")
	// 	request.Header.Set("x-csrftoken", "missing")
	// 	request.Header.Set("cookie", Cookie)
	// 	Response, err := http.DefaultClient.Do(request)
	// 	if err != nil {
	// 		goto Retry
	// 	}
	// 	defer Response.Body.Close()
	// 	bodys, err := ioutil.ReadAll(Response.Body)
	// 	if err != nil {
	// 		goto Retry
	// 	}
	// 	fmt.Println(string(bodys))
	// return string(bodys)
	// fmt.Println(Request1)
}

func ClearConsole() string {
	if System == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	return System

}

func LoadFile(Name string, Path string) []string {
	File, err := os.Open(Path)
	if err != nil {
		Path = Input(""+Name+" Path : ", WHITE)
		fmt.Println()
		return LoadFile(Name, Path)
	}
	var LIST []string
	scanner := bufio.NewScanner(File)
	for scanner.Scan() {
		ReplaceStrings := strings.Join(strings.Fields(scanner.Text()), "")
		ReplaceStrings = strings.Replace(ReplaceStrings, "\n", "", -1)
		ReplaceStrings = strings.Replace(ReplaceStrings, " ", "", -1)
		if ReplaceStrings != "" {
			LIST = append(LIST, ReplaceStrings)
		}

	}
	lenth := len(LIST)
	File.Close()
	if len(LIST) <= 0 {
		PPrint(RED, " ! ", RED, ""+Name+" Is Empty File ", true)
		fmt.Scanln()
		os.Exit(0)
	}
	PPrint(WHITE, " + ", GREEN, "Successfully "+color.HiGreenString("Loaded "+Name+"")+" "+color.HiCyanString("(%d)", lenth)+"", true)
	return LIST
}
func LoadFileUpdater(Path string) []string {
	File, _ := os.Open(Path)
	var LIST []string
	scanner := bufio.NewScanner(File)
	for scanner.Scan() {
		ReplaceStrings := strings.Join(strings.Fields(scanner.Text()), "")
		ReplaceStrings = strings.Replace(ReplaceStrings, "\n", "", -1)
		ReplaceStrings = strings.Replace(ReplaceStrings, " ", "", -1)
		if ReplaceStrings != "" {
			LIST = append(LIST, ReplaceStrings)
		}

	}
	File.Close()
	return LIST

}
func NetRequest(method, Endpoint, Data string, httpHeaders map[string]string, Proxies []string) string {
Re:
	Req, _ := http.NewRequest(method, Endpoint, bytes.NewBufferString(Data))
	for headerType, headerValue := range httpHeaders {
		Req.Header.Set(headerType, headerValue)
	}
	pr := Proxies[rand.Intn(len(Proxies))]
	transpost := http.Transport{
		Proxy:       http.ProxyURL(&url.URL{Host: pr}),
		DialContext: (&net.Dialer{Timeout: 10 * time.Second, KeepAlive: -1, DualStack: true}).DialContext, IdleConnTimeout: 10 * time.Second, TLSHandshakeTimeout: 10 * time.Second, ExpectContinueTimeout: 0, MaxIdleConns: 256, MaxConnsPerHost: 256, MaxIdleConnsPerHost: 256, DisableKeepAlives: false}
	transpost.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := http.Client{Transport: &transpost}
	Response, err := client.Do(Req)

	if err != nil {
		//fmt.Println(err)
		goto Re
	}
	defer Response.Body.Close()
	Body, _ := ioutil.ReadAll(Response.Body)

	return string(Body)

}
func NetApiInsta(method, Endpoint, Endpoint2, Data, Cookie string, proxies []string) string {
	var Body string
	Request1 := NetRequest(method, "https://i.instagram.com/api/v1/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"User-Agent":      "Instagram 133.0.0.34.124 Android (29/10; 480dpi; 1080x2107; HUAWEI; MAR-LX1M; HWMAR; kirin710; en_US; 302733750)",
		"Accept":          "*/*",
		"Host":            "i.instagram.com",
		"Accept-Language": "en-Us",
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":          Cookie}, proxies)
	if Request1 != "" {
		Body = Request1

	}
	return Body

}

func NetWEBinsta(method, Endpoint, Endpoint2, Data, Cookie string, proxies []string) string {
	var Body string
	Request1 := NetRequest(method, "https://www.instagram.com/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"Content-Type":     "application/x-www-form-urlencoded",
		"origin":           "https://www.instagram.com",
		"referer":          "https://www.instagram.com/" + Endpoint + "/" + Endpoint2 + "/",
		"user-agent":       "Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"x-instagram-ajax": "2cd27df1ee26",
		"x-csrftoken":      "2cd27df1ee26",
		"x-ig-app-id":      "missing",
		"x-ig-www-claim":   "0",
		"x-requested-with": "XMLHttpRequest",
		"cookie":           Cookie}, proxies)
	if Request1 != "" {
		Body = Request1
	}
	return Body
}

func RequestApiInsta(method, Endpoint, Endpoint2, Data, Cookie string) string { // -> With out proxy
	var Body string

	Request1 := Request(method, "https://i.instagram.com/api/v1/"+Endpoint+"/"+Endpoint2+"/", Data, map[string]string{
		"User-Agent":      "Instagram 218.0.0.19.108 Android (29/10; 480dpi; 1080x2107; HUAWEI; MAR-LX1M; HWMAR; kirin710; en_US; 302733750)",
		"Accept":          "*/*",
		"Host":            "i.instagram.com",
		"Accept-Language": "en-Us",
		"Content-Type":    "application/x-www-form-urlencoded; charset=UTF-8",
		"cookie":          Cookie})
	if Request1 != "" {
		Body = Request1

	}

	return Body

}
