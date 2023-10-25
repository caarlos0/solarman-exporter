package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (c *Client) CurrentData() (CurrentData, error) {
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"https://globalapi.solarmanpv.com/device/v1.0/currentData?appId=%s&language=en&=",
			c.cfg.AppID,
		),
		strings.NewReader(fmt.Sprintf(`{"deviceSn":%q}`, c.cfg.InverterSN)),
	)
	if err != nil {
		return CurrentData{}, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.c.Do(req)
	if err != nil {
		return CurrentData{}, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentData{}, err
	}

	var data CurrentData
	if err := json.Unmarshal(bts, &data); err != nil {
		return CurrentData{}, err
	}
	return data, nil
}

type CurrentData struct {
	Code           any        `json:"code"`
	Msg            any        `json:"msg"`
	Success        bool       `json:"success"`
	RequestID      string     `json:"requestId"`
	DeviceSn       string     `json:"deviceSn"`
	DeviceID       int        `json:"deviceId"`
	DeviceType     string     `json:"deviceType"`
	DeviceState    int        `json:"deviceState"`
	CollectionTime int        `json:"collectionTime"`
	DataList       []DataList `json:"dataList"`
}

type DataList struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Unit  any    `json:"unit"`
	Name  string `json:"name"`
}
