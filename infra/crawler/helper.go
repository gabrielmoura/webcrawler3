package crawler

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
	"webCrawler3/config"
	"webCrawler3/infra/db/repoDb"
	"webCrawler3/infra/log"
)

var (
	invalidSchemaErr = errors.New("invalid schema")
	ErrLocalLink     = errors.New("local link")
	ErrEmptyPath     = errors.New("empty path")
	ErrDenySuffix    = errors.New("deny suffix")
)

// isDenyPostfix checks if the link has a deny postfix
func isDenyPostfix(url string, denySuffixes []string) bool {
	for _, denySuffix := range denySuffixes {
		if strings.HasSuffix(strings.ToLower(url), denySuffix) {
			return true
		}
	}
	return false
}

// isAllowedSchema checks if the link has an acceptable schema
func isAllowedSchema(link string, acceptableSchema []string) bool {
	nLink, err := url.Parse(link)
	if err != nil {
		log.Logger.Debug("error parsing link in checking schema", zap.String("Link", link), zap.Error(err))
		return false
	}
	for _, schema := range acceptableSchema {
		if nLink.Scheme == schema {
			return true
		}
	}
	return false
}

func isAllowedMIME(contentType string) bool {
	for _, allowedMIME := range config.Conf.CoreConfig.AcceptableMimeTypes.MimeTypes {
		if strings.Contains(contentType, allowedMIME) {
			return true
		}
	}
	return false
}

// checkTLD checks if the link has an acceptable TLD
func checkTLD(link string) bool {
	if len(config.Conf.Filter.Tlds) > 0 {
		linkUrl, err := url.Parse(link)
		if err != nil {
			return false
		}
		for _, tld := range config.Conf.Filter.Tlds {
			if strings.HasSuffix(linkUrl.Hostname(), tld) {
				return true
			}
		}
		return false
	}
	return true
}

func handleAddToQueue(links []string, depth, delay int) {
	for _, link := range links {
		nLink := strings.TrimSpace(link)
		if checkTLD(nLink) && isAllowedSchema(nLink, config.Conf.CoreConfig.AcceptableSchema.Schemas) {
			_, err := repoDb.AddToQueueIfNotExists(context.Background(), nLink, depth, delay)
			if err != nil {
				if errors.Is(err, repoDb.ErrMaxDepthReached) {
					return
				}
				log.Logger.Error("error adding link to queue", zap.String("Link", link), zap.Error(err))
				return
			}
		}

	}
}

// isLocalLink verifica se o link é local,
// caso definido para ignorar não adiciona a fila
func isLocalLink(link *url.URL) bool {
	if !config.Conf.Filter.IgnoreLocal {
		return false
	}
	return link.Host == "localhost" || link.Host == "127.0.0.1"
}

func prepareLink(link string) (*url.URL, error) {
	linkUrl, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if isLocalLink(linkUrl) {
		return nil, ErrLocalLink
	}

	if linkUrl.Scheme == "" {
		return nil, invalidSchemaErr
	}
	q, _ := url.ParseQuery(linkUrl.RawQuery)
	q.Del("utm_source")
	q.Del("utm_medium")
	q.Del("utm_campaign")
	q.Del("utm_term")
	q.Del("utm_content")
	q.Del("#")
	linkUrl.RawQuery = q.Encode()

	if isDenyPostfix(linkUrl.Path, config.Conf.CoreConfig.DenySuffixes.Suffixes) {
		return nil, ErrDenySuffix
	}

	return linkUrl, nil
}
func prepareParentLink(parentLink, link string) (*url.URL, error) {

	// Remove o primeiro caractere se for uma barra ou ponto
	if strings.HasPrefix(link, "/") || strings.HasPrefix(link, ".") {
		link = link[1:]
	}

	// Se o link for vazio, retorna erro
	nURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if nURL.Path == "" {
		return nil, ErrEmptyPath
	}

	if nURL.IsAbs() {
		log.Logger.Debug("New URL", zap.String("URL", nURL.String()), zap.Bool("absolute", true))
		return nURL, nil
	}

	// Se o link for relativo, adiciona o host e schema do link pai
	pURL, err := url.Parse(parentLink)
	if err != nil {
		return nil, err
	}

	// Se o link for relativo, adiciona o host e schema do link pai
	nURL.Host = pURL.Host
	nURL.Scheme = pURL.Scheme
	nURL.Path = strings.TrimRight(pURL.Path, "/") + link
	log.Logger.Debug("New URL", zap.String("URL", nURL.String()), zap.Bool("absolute", false))

	return nURL, nil
}
func isStatusErr(status int, url *url.URL) bool {
	if status == http.StatusOK {
		return false
	}

	i2p := strings.HasSuffix(url.Host, ".i2p")
	switch status {
	case http.StatusConflict:
		if i2p {
			handleListHelperI2P(url.String())
		}
		return true
	default:
		return status < 200 || status >= 300
	}
}

func handleListHelperI2P(url string) {
	if _, err := repoDb.AddHelperI2P(context.Background(), url); err != nil {
		log.Logger.Error("Error adding helper I2P", zap.Error(err))
	}
}
