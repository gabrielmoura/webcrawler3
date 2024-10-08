package crawler

import (
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"time"
	"webCrawler3/config"
	"webCrawler3/infra/log"
)

func proxyClient() *http.Client {
	urlProxy, err := url.Parse(config.Conf.Proxy.ProxyURL)
	if err != nil {
		log.Logger.Panic("Error parsing proxy URL", zap.Error(err))
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(urlProxy),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second, // Definir um timeout de 5 segundos
	}

	return client
}
func httpClient() *http.Client {
	if config.Conf.Proxy.Enabled {
		return proxyClient()
	} else {
		return &http.Client{
			Timeout: 5 * time.Second, // Definir um timeout de 5 segundos
			Transport: &http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 10,
			},
		}
	}
}

func httpRequest(pageUrl string) (resp *http.Response, err error) {
	client := httpClient()
	req, err := http.NewRequest("GET", pageUrl, nil)
	if err != nil {
		log.Logger.Error("Error creating request", zap.Error(err))
		return
	}
	req.Header.Set("User-Agent", config.Conf.UserAgent)
	resp, err = client.Do(req)
	return
}
