package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/baraa-almasri/shortsninja/config"
	"github.com/baraa-almasri/shortsninja/models"
	"net/http"
	"strings"
)

// RequestDataManager manages data from requests
type RequestDataManager struct {
	conf      *config.Config
	dbManager models.Database
}

// NewRequestDataManager returns a new RequestDataManager instance
func NewRequestDataManager(config *config.Config, dbManager models.Database) *RequestDataManager {
	return &RequestDataManager{
		conf:      config,
		dbManager: dbManager,
	}
}

// GetURLDataFromRequestData returns a URLData instance with the needed data
func (reqData *RequestDataManager) GetURLDataFromRequestData(req *http.Request) *models.URLData {
	ip := reqData.GetIP(req)
	return &models.URLData{
		IP:            ip,
		VisitLocation: reqData.getIPLocation(ip),
		UserAgent:     req.Header.Get("User-Agent"),
	}
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func (reqData *RequestDataManager) GetIP(req *http.Request) string {
	forwarded := req.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded[:strings.Index(forwarded, ":")]
	}
	return req.RemoteAddr[:strings.Index(req.RemoteAddr, ":")]
}

// getIPLocation return a string of the IP's location using ipinfo.io
func (reqData *RequestDataManager) getIPLocation(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://ipinfo.io/%s?token=%s", ip, reqData.conf.IPInfoIoToken))
	if err != nil {
		return "NULL/NULL"
	}

	defer resp.Body.Close()

	ipData := make(map[string]interface{})
	err = json.NewDecoder(resp.Body).Decode(&ipData)
	if err != nil {
		return "NULL/NULL"
	}

	return fmt.Sprintf("%s/%s", ipData["region"], ipData["country"])
}
