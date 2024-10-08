package crawler

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"regexp"
	"strings"
	"webCrawler3/config"
	"webCrawler3/infra/log"
	"webCrawler3/infra/pb"
)

var InvalidMeta = errors.New("invalid meta")

func ignoreTags(data []byte) []byte {
	htmlRegex := regexp.MustCompile("(?s)<(script|style|noscript|link|meta)[^>]*?>.*?</(script|style|noscript|link|meta)>")
	return htmlRegex.ReplaceAll(data, []byte(""))
}
func extractOnlyText(data []byte) []byte {
	tagsRegex := regexp.MustCompile("<([^>]*)>")
	return tagsRegex.ReplaceAll(data, []byte(""))
}
func removeSpecialChars(data []byte) []byte {
	wordRegex := regexp.MustCompile("[^\\pL\\pN\\pZ'-]+")
	return wordRegex.ReplaceAll(data, []byte(" "))
}

// countWordsInText Extrai e conta a frequência de palavras do conteúdo HTML, ignorando palavras irrelevantes comuns.
func countWordsInText(data []byte) (config.MapWordsType, error) {
	log.Logger.Debug("Word Count")
	// Etapa 1: Ignorar determinadas tags HTML
	parcialPlainText := ignoreTags(data)

	// Etapa 2: Remover tags HTML
	plainText := extractOnlyText(parcialPlainText)

	// Etapa 3: Normalizar texto
	normalizedText := bytes.ToLower(plainText)

	// Etapa 4: Remova caracteres especiais e divida em palavras
	noSpecialCh := removeSpecialChars(normalizedText)

	// Etapa 5: Divida o texto em palavras
	words := bytes.Split(noSpecialCh, []byte(" "))

	// Etapa 6: Conte a frequência das palavras (ignorando palavras comuns)
	wordCounts := make(config.MapWordsType)
	for _, wordBytes := range words {
		word := string(bytes.TrimSpace(wordBytes))

		// Pule palavras curtas e palavras de parada comuns
		if len(word) < 2 || containsStopWord(word, "") {
			continue
		}
		wordCounts[word]++
		log.Logger.Debug("Word: ", zap.Uint32(word, wordCounts[word]))
	}

	return wordCounts, nil
}

// ContainsMap Verifica se uma palavra está em uma lista de stop words comuns,
func containsStopWord(item string, language string) bool {
	wordMap := config.Conf.CoreConfig.CommonStopWords.LanguageStopWords

	// Check for the specified language first, if provided
	if language != "" {
		if stopWords := findStopWordsByLanguage(wordMap, language); stopWords != nil {
			if contains(stopWords, item) {
				return true
			}
		}
	}

	// Check across all stop words if the language check didn't yield a result
	return containsAllStopWords(wordMap, item)
}

// findStopWordsByLanguage retrieves stop words for the specified language
func findStopWordsByLanguage(wordMap []*config.LanguageStopWords, language string) []string {
	for _, entry := range wordMap {
		if entry.LanguageCode == language {
			return entry.StopWords
		}
	}
	return nil
}

// contains checks if the item exists in the provided slice of stop words
func contains(stopWords []string, item string) bool {
	for _, word := range stopWords {
		if word == item {
			return true
		}
	}
	return false
}

// containsAllStopWords checks if the item exists in all stop words across all languages
func containsAllStopWords(wordMap []*config.LanguageStopWords, item string) bool {
	for _, entry := range wordMap {
		if contains(entry.StopWords, item) {
			return true
		}
	}
	return false
}

func extractData(n *html.Node, dataPage *pb.Page) error {

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				extractTitle(n, dataPage)
			case "meta":
				extractMeta(n, dataPage)
			case "script":
				extractJSONLD(n, dataPage)
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return nil
}

func extractDescription(n *html.Node) string {
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == "description" {
			for _, a := range n.Attr {
				if a.Key == "content" {
					return a.Val
				}
			}
		}
	}
	return ""
}

func extractOG(n *html.Node) (map[string]string, error) {
	ogData := make(map[string]string)
	if n.Data != "meta" {
		return nil, InvalidMeta
	}
	for _, a := range n.Attr {
		if a.Key == "property" && len(a.Val) > 3 && a.Val[:3] == "og:" {
			for _, b := range n.Attr {
				if b.Key == "content" {
					ogData[a.Val] = b.Val
				}
			}
		}
	}

	return ogData, nil
}

func extractKeywords(n *html.Node) ([]string, error) {
	if n.Data != "meta" {
		return nil, InvalidMeta
	}
	for _, a := range n.Attr {
		if a.Key == "name" && a.Val == "keywords" {
			for _, a := range n.Attr {
				if a.Key == "content" {
					return strings.Split(a.Val, ","), nil
				}
			}
		}
	}
	return nil, nil
}

func extractManifest(n *html.Node) string {
	if n.Data == "link" {
		var isManifest bool
		for _, a := range n.Attr {
			if a.Key == "rel" && a.Val == "manifest" {
				isManifest = true
			}
			if a.Key == "href" && isManifest {
				return a.Val
			}
		}
	}
	return ""
}

func extractMeta(n *html.Node, dataPage *pb.Page) {
	if n.Data == "meta" {
		if dataPage.Meta == nil {
			dataPage.Meta = metaNull
		}
		description := extractDescription(n)
		if description != "" {
			dataPage.Description = description
		}

		ogData, err := extractOG(n)
		if err == nil {
			for k, v := range ogData {
				dataPage.Meta.OG[k] = v
			}
		}

		keywords, err := extractKeywords(n)
		if keywords != nil && err == nil {
			dataPage.Meta.Keywords = keywords
		}
	}

	manifest := extractManifest(n)
	if manifest != "" {
		dataPage.Meta.Manifest = manifest
	}

	if dataPage.Meta.Manifest == "" && len(dataPage.Meta.OG) == 0 && len(dataPage.Meta.Keywords) == 0 {
		dataPage.Meta = nil
	}
}

func extractTitle(n *html.Node, dataPage *pb.Page) {
	if n.Data == "title" && n.FirstChild != nil {
		dataPage.Title = n.FirstChild.Data
	}
}

// extractLinks Extrai links de um documento HTML.
func extractLinks(parentLink string, n *html.Node) ([]string, error) {
	var links []string

	var extract func(*html.Node)
	extract = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					urlE, err := prepareLink(a.Val)
					if err != nil {
						if errors.Is(invalidSchemaErr, err) {
							preparedLink, err := prepareParentLink(parentLink, a.Val)
							if err != nil {
								continue
							}
							urlE = preparedLink
						}
						log.Logger.Debug(fmt.Sprintf("Error preparing link: %s", err))
						continue
					}
					links = append(links, urlE.String())
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extract(c)
		}
	}

	extract(n)
	return links, nil
}
func extractJSONLD(n *html.Node, dataPage *pb.Page) {
	for _, a := range n.Attr {
		if a.Key == "type" && a.Val == "application/ld+json" {
			if dataPage.Meta == nil {
				dataPage.Meta = metaNull
			}
			if n.FirstChild != nil {
				content := n.FirstChild.Data
				dataPage.Meta.Ld = content
			}
		}
	}
}

var metaNull = &pb.PageMeta{
	OG:       make(map[string]string),
	Keywords: []string{},
	Manifest: "",
	Ld:       "",
}
