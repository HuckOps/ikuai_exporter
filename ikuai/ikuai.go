package ikuai

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60"
const contentType = "application/json;charset=UTF-8"

func md5Sum(input string) string {
	hasher := md5.New()
	hasher.Write([]byte(input))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

type ikuaiClient struct {
	IP       string
	UserName string
	Password string
	Cookies  []*http.Cookie
}

func NewClient(ip string, userName string, password string) *ikuaiClient {
	return &ikuaiClient{
		IP:       ip,
		UserName: userName,
		Password: password,
	}
}

func (c *ikuaiClient) Login() (err error) {
	loginData := map[string]interface{}{
		"username": c.UserName,
		"passwd":   md5Sum(c.Password),
		"pass": base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf("salt_11%s", c.Password))),
		"remember_password": nil,
	}
	jsonData, _ := json.Marshal(loginData)
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/Action/login", c.IP), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	data := map[string]interface{}{}
	json.NewDecoder(resp.Body).Decode(&data)
	if int(data["Result"].(float64)) != 10000 {
		return errors.New("Auth failed")
	}
	cookies := resp.Cookies()
	c.Cookies = []*http.Cookie{
		{
			Name:  "username",
			Value: c.UserName,
		},
		{
			Name:  "login",
			Value: "1",
		},
	}
	for _, cookie := range cookies {
		if cookie.Name == "sess_key" {
			c.Cookies = append(c.Cookies, &http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}
	return nil
}

type CallBody struct {
	Action   string                 `json:"action"`
	FuncName string                 `json:"func_name"`
	Param    map[string]interface{} `json:"param"`
}

func (c *ikuaiClient) Call(callBody CallBody, result interface{}) (err error) {
	url := fmt.Sprintf("http://%s/Action/call", c.IP)
	jsonData, _ := json.Marshal(callBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	for _, cookie := range c.Cookies {
		req.AddCookie(cookie)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", contentType)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(result)
	return
}

type Response[T any] struct {
	Result int    `json:"Result"`
	ErrMsg string `json:"ErrMsg"`
	Data   T      `json:"Data"`
}

type Sysstat struct {
	CPU    []string `json:"cpu"`
	Memory Memory   `json:"memory"`
}

// Memory represents memory usage details.
type Memory struct {
	Total     int    `json:"total"`
	Available int    `json:"available"`
	Free      int    `json:"free"`
	Cached    int    `json:"cached"`
	Buffers   int    `json:"buffers"`
	Used      string `json:"used"`
}

type SysStatResult struct {
	CPUPercent  float64
	Buffer      float64
	Available   float64
	Cache       float64
	Free        float64
	Total       float64
	MemoryUsage float64
}

func (c *ikuaiClient) GetSysstat() (sysStat SysStatResult) {
	body := CallBody{
		Action:   "show",
		FuncName: "sysstat",
		Param: map[string]interface{}{
			"TYPE": "verinfo,cpu,memory,stream,cputemp",
		},
	}

	result := Response[Sysstat]{}
	if err := c.Call(body, &result); err != nil {
		log.Fatal(err)
		return
	}

	// cpu usage
	cpuPercent := result.Data.CPU[0]
	cpuPercent = strings.Replace(cpuPercent, "%", "", -1)
	cpuPercentFloat, err := strconv.ParseFloat(cpuPercent, 2)
	if err != nil {
		log.Fatal(err)
	}
	sysStat.CPUPercent = cpuPercentFloat

	// memory info
	sysStat.Available = float64(result.Data.Memory.Available)
	sysStat.Buffer = float64(result.Data.Memory.Buffers)
	sysStat.Cache = float64(result.Data.Memory.Cached)
	sysStat.Free = float64(result.Data.Memory.Free)
	sysStat.Total = float64(result.Data.Memory.Total)
	memoryUsage := result.Data.Memory.Used
	memoryUsage = strings.Replace(memoryUsage, "%", "", -1)
	sysStat.MemoryUsage, err = strconv.ParseFloat(memoryUsage, 2)
	if err != nil {
		log.Fatal(err)
	}
	return
}

type IfaceStatus struct {
	IfaceCheck  []IfaceCheck  `json:"iface_check"`
	IfaceStream []IfaceStream `json:"iface_stream"`
}

// IfaceCheck represents an interface check entry.
type IfaceCheck struct {
	ID              int    `json:"id"`
	Interface       string `json:"interface"`
	ParentInterface string `json:"parent_interface"`
	IpAddr          string `json:"ip_addr"`
	Gateway         string `json:"gateway"`
	Internet        string `json:"internet"`
	UpdateTime      string `json:"updatetime"`
	AutoSwitch      string `json:"auto_switch"`
	Result          string `json:"result"`
	Errmsg          string `json:"errmsg"`
	Comment         string `json:"comment"`
}

// IfaceStream represents an interface stream entry.
type IfaceStream struct {
	Interface   string `json:"interface"`
	Comment     string `json:"comment"`
	IpAddr      string `json:"ip_addr"`
	ConnectNum  string `json:"connect_num"`
	Upload      int    `json:"upload"`
	Download    int    `json:"download"`
	TotalUp     int    `json:"total_up"`
	TotalDown   int    `json:"total_down"`
	UpDropped   int    `json:"updropped"`
	DownDropped int    `json:"downdropped"`
	UpPacked    int    `json:"uppacked"`
	DownPacked  int    `json:"downpacked"`
}

func (c *ikuaiClient) GetIface() IfaceStatus {
	body := CallBody{
		Action:   "show",
		FuncName: "monitor_iface",
		Param: map[string]interface{}{
			"TYPE": "iface_check,iface_stream",
		},
	}
	result := Response[IfaceStatus]{}
	if err := c.Call(body, &result); err != nil {
		log.Fatal(err)
		return IfaceStatus{}
	}
	return result.Data
}

// Data holds the data part of the response.
type LanIP struct {
	Items []LanIPItem `json:"data"`
	Total int         `json:"total"`
}

// Item represents an individual data item in the 'data' array.
type LanIPItem struct {
	WebID        int    `json:"webid"`
	IPAddr       string `json:"ip_addr"`
	DownRate     string `json:"downrate"`
	TotalUp      int    `json:"total_up"`
	TotalDown    int    `json:"total_down"`
	UpRate       string `json:"uprate"`
	Signal       string `json:"signal"`
	PPPType      string `json:"ppptype"`
	Hostname     string `json:"hostname"`
	ID           int    `json:"id"`
	LinkAddr     string `json:"link_addr"`
	BSSID        string `json:"bssid"`
	IPAddrInt    int    `json:"ip_addr_int"`
	ConnectNum   int    `json:"connect_num"`
	Upload       int    `json:"upload"`
	Download     int    `json:"download"`
	AuthType     int    `json:"auth_type"`
	ClientType   string `json:"client_type"`
	ClientDevice string `json:"client_device"`
	Timestamp    int    `json:"timestamp"`
	APName       string `json:"apname"`
	ACGID        int    `json:"ac_gid"`
	APMAC        string `json:"apmac"`
	Username     string `json:"username"`
	MAC          string `json:"mac"`
	Reject       int    `json:"reject"`
	SSID         string `json:"ssid"`
	Frequencies  string `json:"frequencies"`
	DTalkName    string `json:"dtalk_name"`
	Comment      string `json:"comment"`
}

func (c *ikuaiClient) GetLanIPs() []LanIPItem {
	body := CallBody{
		Action:   "show",
		FuncName: "monitor_lanip",
		Param: map[string]interface{}{
			"TYPE": "data,total",
		},
	}
	result := Response[LanIP]{}
	if err := c.Call(body, &result); err != nil {
		log.Fatal(err)
		return []LanIPItem{}
	}
	return result.Data.Items
}
