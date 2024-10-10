package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
	"webCrawler3/infra/db/repoDb"
)

func HandlerGraph(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		graphHandler(w, r)
	default:
		responseWithMethodNotAllowed(w)
	}
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	ctxTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	data, err := returnGraph(ctxTimeout)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	//w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func returnGraphWithDirection(ctx context.Context) ([]byte, error) {
	// Obtém as páginas do banco de dados para construir o gráfico
	pages, err := repoDb.GetPagesForGraph(ctx)
	if err != nil {
		return nil, err // Retorna erro se não for possível obter as páginas
	}

	// Inicializa a estrutura do gráfico
	graph := Graph{
		Nodes: []Node{},
		Edges: []Edge{},
	}

	// Mapa para associar domínios a IDs de nós e IDs de páginas a IDs incrementais
	nodeMap := make(map[string]int)
	pageIDMap := make(map[uuid.UUID]int)
	nextID := 1 // Inicializa o ID incremental

	// Processa cada página para criar nós e mapear page.ID a um ID incremental
	for _, page := range pages {
		// Verifica se o título e a URL são válidos
		if page.Title.Valid && page.Url.Valid {
			parsedUrl, err := url.Parse(page.Url.String)
			if err == nil {
				domain := parsedUrl.Host

				// Se o domínio já existir, adiciona os links ao nó existente
				if nodeID, exists := nodeMap[domain]; exists {
					for i := range graph.Nodes {
						if graph.Nodes[i].ID == nodeID {
							graph.Nodes[i].Links = append(graph.Nodes[i].Links, page.Links...)
							break
						}
					}
				} else {
					// Cria um novo nó para o domínio com o ID incremental
					graph.Nodes = append(graph.Nodes, Node{ID: nextID, Label: parsedUrl.Host, Links: page.Links})
					nodeMap[domain] = nextID // Mapeia o domínio ao novo ID de nó
					nextID++                 // Incrementa o próximo ID
				}

				// Mapeia page.ID (UUID) ao ID incremental
				pageIDMap[page.ID] = nodeMap[domain]
			}
		}
	}

	// Conjunto para acompanhar arestas já adicionadas
	edgesSet := make(map[string]struct{})

	// Processa cada página novamente para criar arestas
	for _, page := range pages {
		if page.Title.Valid && page.Url.Valid {
			for _, link := range page.Links {
				parsedLink, err := url.Parse(link)
				if err == nil {
					domainLink := parsedLink.Host
					if targetNodeID, exists := nodeMap[domainLink]; exists {
						if sourceNodeID, exists := pageIDMap[page.ID]; exists {
							edgeKey := fmt.Sprintf("%d->%d", sourceNodeID, targetNodeID)
							if _, exists := edgesSet[edgeKey]; !exists {
								graph.Edges = append(graph.Edges, Edge{From: sourceNodeID, To: targetNodeID})
								edgesSet[edgeKey] = struct{}{} // Marca a aresta como adicionada
							}
						}
					}
				}
			}
		}
	}

	// Serializa o gráfico em formato JSON
	result, err := json.Marshal(graph)
	if err != nil {
		return nil, err // Retorna erro se ocorrer falha na serialização
	}

	return result, nil // Retorna o gráfico serializado
}
func returnGraph(ctx context.Context) ([]byte, error) {
	// Obtém as páginas do banco de dados para construir o gráfico
	pages, err := repoDb.GetPagesForGraph(ctx)
	if err != nil {
		return nil, err // Retorna erro se não for possível obter as páginas
	}

	// Inicializa a estrutura do gráfico
	graph := Graph{
		Nodes: []Node{},
		Edges: []Edge{},
	}

	// Mapa para associar domínios a IDs de nós e IDs de páginas a IDs incrementais
	nodeMap := make(map[string]int)
	pageIDMap := make(map[uuid.UUID]int)
	nextID := 1 // Inicializa o ID incremental

	// Processa cada página para criar nós e mapear page.ID a um ID incremental
	for _, page := range pages {
		// Verifica se o título e a URL são válidos
		if page.Title.Valid && page.Url.Valid {
			parsedUrl, err := url.Parse(page.Url.String)
			if err == nil {
				domain := parsedUrl.Host

				// Se o domínio já existir, adiciona os links ao nó existente
				if nodeID, exists := nodeMap[domain]; exists {
					for i := range graph.Nodes {
						if graph.Nodes[i].ID == nodeID {
							graph.Nodes[i].Links = append(graph.Nodes[i].Links, page.Links...)
							break
						}
					}
				} else {
					// Cria um novo nó para o domínio com o ID incremental
					graph.Nodes = append(graph.Nodes, Node{ID: nextID, Label: parsedUrl.Host, Links: page.Links})
					nodeMap[domain] = nextID // Mapeia o domínio ao novo ID de nó
					nextID++                 // Incrementa o próximo ID
				}

				// Mapeia page.ID (UUID) ao ID incremental
				pageIDMap[page.ID] = nodeMap[domain]
			}
		}
	}

	// Conjunto para acompanhar arestas já adicionadas
	edgesSet := make(map[string]struct{})

	// Processa cada página novamente para criar arestas
	for _, page := range pages {
		if page.Title.Valid && page.Url.Valid {
			for _, link := range page.Links {
				parsedLink, err := url.Parse(link)
				if err == nil {
					domainLink := parsedLink.Host
					if targetNodeID, exists := nodeMap[domainLink]; exists {
						if sourceNodeID, exists := pageIDMap[page.ID]; exists {
							// Cria uma chave que considera a ordem dos nós
							edgeKey := fmt.Sprintf("%d_%d", Min(sourceNodeID, targetNodeID), Max(sourceNodeID, targetNodeID))
							if _, exists := edgesSet[edgeKey]; !exists {
								graph.Edges = append(graph.Edges, Edge{From: sourceNodeID, To: targetNodeID})
								edgesSet[edgeKey] = struct{}{} // Marca a aresta como adicionada
							}
						}
					}
				}
			}
		}
	}

	// Serializa o gráfico em formato JSON
	result, err := json.Marshal(graph)
	if err != nil {
		return nil, err // Retorna erro se ocorrer falha na serialização
	}

	return result, nil // Retorna o gráfico serializado
}
