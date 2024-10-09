package crawler

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"net/url"
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
		Timeout:   config.Conf.GetReqTimeout(), // Definir um timeout de 5 segundos
	}

	return client
}
func httpClient() *http.Client {
	if config.Conf.Proxy.Enabled {
		return proxyClient()
	} else {
		return &http.Client{
			Timeout: config.Conf.GetReqTimeout(), // Definir um timeout de 5 segundos
			Transport: &http.Transport{
				//MaxIdleConns:        2,                           // Limite de conexões ociosas
				//MaxIdleConnsPerHost: 20,                          // Limite de conexões ociosas por host
				IdleConnTimeout: config.Conf.GetReqTimeout(), // Tempo máximo para conexões ociosas
			},
		}
	}
}

func httpRequest(ctx context.Context, pageUrl string) (resp *http.Response, err error) {
	context.WithTimeout(ctx, config.Conf.GetReqTimeout())
	client := httpClient()
	req, err := http.NewRequest("GET", pageUrl, nil)
	if err != nil {
		log.Logger.Error("Error creating request", zap.Error(err))
		return
	}
	req.Header.Set("User-Agent", config.Conf.UserAgent)

	req.Header.Set("Accept", config.Conf.GetAcceptableMimeType())

	if config.Conf.HintsEnabled() {
		for key, value := range config.Conf.GetHints() {
			req.Header.Set(key, value)
		}
	}

	req.WithContext(ctx)

	resp, err = client.Do(req)
	return
}
