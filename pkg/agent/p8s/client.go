package p8s

import (
	"fmt"
	"time"

	"github.com/jaehoonkim/sentinel/pkg/agent/httpclient"
	"github.com/jaehoonkim/sentinel/pkg/agent/log"
)

const defaultQueryTimeout = 10 * time.Second

type Client struct {
	client *httpclient.HttpClient
}

func NewClient(url string) (*Client, error) {
	client, err := httpclient.NewHttpClient(url, false, 0, 0)
	if err != nil {
		return nil, err
	}
	client.SetDisableKeepAlives()
	return &Client{client: client}, nil
}

func (c *Client) ApiRequest(apiVersion, apiName string, queryParams map[string]interface{}) (string, error) {
	var data string
	var warnings []string
	var err error

	switch apiVersion {
	case "v1":
		switch apiName {
		case "query":
			data, warnings, err = c.Query(queryParams)
		case "query_range":
			data, warnings, err = c.QueryRange(queryParams)
		case "alerts":
			data, err = c.Alerts()
		case "rules":
			data, err = c.Rules()
		case "alertmanagers":
			data, err = c.AlertManagers()
		case "targets":
			data, err = c.Targets()
		case "targets/metadata":
			data, err = c.TargetsMetadata(queryParams)
		default:
			return "", fmt.Errorf("unknown api name(%s)", apiName)
		}
	default:
		return "", fmt.Errorf("unknown api version(%s)", apiVersion)
	}

	if len(warnings) > 0 {
		log.Warnf("Prometheus API(%s) Warnings : %v\n", apiName, warnings)
	}

	return data, err
}
