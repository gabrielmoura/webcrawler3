package repoDb

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strings"
	"sync"
	"webCrawler3/config"
	"webCrawler3/infra/log"

	"webCrawler3/infra/db"
	"webCrawler3/infra/pb"

	"github.com/jackc/pgx/v5"
)

var (

	//go:embed schema.sql
	schema    string
	DbQueries *db.Queries
	Pool      *pgxpool.Pool

	addQueue = sync.Mutex{}

	ErrMaxDepthReached = errors.New("max depth reached")
)

// Initialize sets up the database connection and schema
func Initialize(ctx context.Context, cfg config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.PostgresURI)
	if err != nil {
		log.Logger.Panic("Error parsing Postgres URI", zap.Error(err))
	}
	poolConfig.MaxConns = cfg.GetDBMaxConns()
	poolConfig.MinConns = cfg.GetDBMinConns()

	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// ...
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if cfg.SeedDb {
		// Create tables
		if _, err := pool.Exec(ctx, schema); err != nil {
			return nil, err
		}
	}
	DbQueries = db.New(pool)
	Pool = pool
	return pool, nil
}

// AddToQueue inserts a new URL into the processing queue
func addToQueue(ctx context.Context, url string, depth, delay int) (db.Queue, error) {
	queueEntry, err := DbQueries.QueueInsert(ctx, db.QueueInsertParams{
		Url:    url,
		Depth:  int32(depth),
		Status: db.StatusPending,
		Delay:  int32(delay),
	})
	return queueEntry, err
}

func AddToQueueIfNotExists(ctx context.Context, url string, depth, delay int) (db.Queue, error) {
	addQueue.Lock()
	defer addQueue.Unlock()

	if depth > config.Conf.MaxDepth {
		return db.Queue{}, ErrMaxDepthReached
	}

	queueEntry, err := DbQueries.QueueGetAllByStatusAndUrl(ctx, db.QueueGetAllByStatusAndUrlParams{
		Status: db.StatusPending,
		Url:    url,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return addToQueue(ctx, url, depth, delay)
		}
		return db.Queue{}, err
	}
	return queueEntry, nil
}

// RetrieveFromQueue fetches pending URLs from the queue
func RetrieveFromQueue(ctx context.Context, limit, offset int) ([]db.Queue, error) {
	// Set default limit if not provided
	if limit == 0 {
		limit = 10
	}

	params := &db.QueueGetAllByStatusParams{
		Status: db.StatusPending,
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	var queueItems []db.Queue

	err := runWithTX(ctx, Pool, func(tx *db.Queries) error {
		items, err := DbQueries.QueueGetAllByStatus(ctx, *params)
		if err != nil {
			log.Logger.Error("Error retrieving queue items", zap.Error(err))
			return err
		}
		for _, queueItem := range items {
			err := DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
				ID:     queueItem.ID,
				Status: db.StatusProcessing,
			})
			if err != nil {
				log.Logger.Error("Error marking as processing", zap.Error(err))
				return err
			}
		}
		queueItems = items

		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(queueItems) == 0 {
		return nil, pgx.ErrNoRows
	}

	return queueItems, nil
}
func MarkAsVisited(ctx context.Context, id uuid.UUID) error {
	return DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
		ID:     id,
		Status: db.StatusVisited,
	})
}
func MarkAsError(ctx context.Context, id uuid.UUID) error {
	return DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
		ID:     id,
		Status: db.StatusError,
	})
}
func MarkAsSkipped(ctx context.Context, id uuid.UUID) error {
	addQueue.Lock()
	defer addQueue.Unlock()
	return DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
		ID:     id,
		Status: db.StatusSkipped,
	})
}
func MarkAsProcessing(ctx context.Context, id uuid.UUID) error {
	return DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
		ID:     id,
		Status: db.StatusProcessing,
	})
}

// SavePage inserts a new page record into the database
func SavePage(ctx context.Context, page *pb.Page) (db.Page, error) {
	metaData, _ := json.Marshal(page.Meta)
	wordList, _ := json.Marshal(page.Words)
	return DbQueries.InsertPage(ctx, db.InsertPageParams{
		Url:         page.Url,
		Title:       page.Title,
		Description: page.Description,
		Meta:        metaData,
		Words:       wordList,
		Visited:     true,
		Links:       page.Links,
	})
}

func SavePageAndMarkVisited(ctx context.Context, page *pb.Page, queueId uuid.UUID) error {
	return runWithTX(ctx, Pool, func(tx *db.Queries) error {
		metaData, _ := json.Marshal(page.Meta)
		wordList, _ := json.Marshal(page.Words)
		_, err := DbQueries.InsertPage(ctx, db.InsertPageParams{
			Url:         page.Url,
			Title:       page.Title,
			Description: page.Description,
			Meta:        metaData,
			Words:       wordList,
			Visited:     true,
			Links:       page.Links,
		})
		if err != nil {
			return err
		}
		return DbQueries.QueueMarkAs(ctx, db.QueueMarkAsParams{
			ID:     queueId,
			Status: db.StatusVisited,
		})
	})
}

// PageHasVisited checks if a page has been visited
func PageHasVisited(ctx context.Context, url string) (bool, error) {
	page, err := DbQueries.QueueGetAllByStatusAndUrl(ctx, db.QueueGetAllByStatusAndUrlParams{
		Status: db.StatusVisited,
		Url:    url,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return page.Status == db.StatusVisited, nil
}

func runWithTX(ctx context.Context, c *pgxpool.Pool, fn func(*db.Queries) error) error {
	tx, err := c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if errRb := tx.Rollback(ctx); errRb != nil {
			return fmt.Errorf("error on rollback: %v, original error: %w", errRb, err)
		}
		return err
	}
	return tx.Commit(ctx)
}

// PrepareHealth consolidates the health of the database
func PrepareHealth(ctx context.Context) []error {
	var errorList []error

	// 1. Retrieve all queues with a pending status and mark duplicates as skipped
	errorList = append(errorList, runWithTX(ctx, Pool, func(tx *db.Queries) error {
		// Fetch queues with the status 'PENDING'
		queues, err := tx.QueueGetAllByStatus(ctx, db.QueueGetAllByStatusParams{
			Status: db.StatusPending,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			return err // Return error if fetching queues fails
		}
		log.Logger.Debug("Processing queues", zap.Int("count", len(queues)))

		// Track visited URLs to identify duplicates
		visitedUrls := make(map[string]bool)
		for _, queue := range queues {
			if _, exists := visitedUrls[queue.Url.String]; exists {
				// Mark as skipped if the URL has already been visited
				if err := tx.QueueMarkAs(ctx, db.QueueMarkAsParams{
					ID:     queue.ID,
					Status: db.StatusSkipped,
				}); err != nil {
					return err // Return error if marking fails
				}
				log.Logger.Debug("Marked as skipped", zap.String("url", queue.Url.String))
			}
			visitedUrls[queue.Url.String] = true // Mark URL as visited
		}
		return nil // Successfully processed all queues
	}))

	// 2. Retrieve all pages and mark duplicates in the queue as visited
	errorList = append(errorList, runWithTX(ctx, Pool, func(tx *db.Queries) error {
		// Fetch all pages
		pages, err := tx.GetPages(ctx, db.GetPagesParams{
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			return err // Return error if fetching pages fails
		}
		log.Logger.Debug("Processing pages", zap.Int("count", len(pages)))

		// Track visited URLs for pages
		visitedPages := make(map[string]bool)
		for _, page := range pages {
			if _, exists := visitedPages[page.Url.String]; exists {
				// Mark as visited if the URL has already been processed
				if err := tx.QueueMarkAsByUrl(ctx, db.QueueMarkAsByUrlParams{
					Url:    page.Url.String,
					Status: db.StatusVisited,
				}); err != nil {
					return err // Return error if marking fails
				}
				log.Logger.Debug("Marked as visited", zap.String("url", page.Url.String))
			}
			visitedPages[page.Url.String] = true // Mark URL as visited
		}
		return nil // Successfully processed all pages
	}))

	// 3. Retrieve all queue marked as processing and mark as pending
	errorList = append(errorList, runWithTX(ctx, Pool, func(tx *db.Queries) error {
		// Fetch all queues marked as processing
		queues, err := tx.QueueGetAllByStatus(ctx, db.QueueGetAllByStatusParams{
			Status: db.StatusProcessing,
			Limit:  100,
			Offset: 0,
		})
		if err != nil {
			return err // Return error if fetching queues fails
		}
		log.Logger.Debug("Processing queues", zap.Int("count", len(queues)))

		// Mark queues as pending
		for _, queue := range queues {
			if err := tx.QueueMarkAs(ctx, db.QueueMarkAsParams{
				ID:     queue.ID,
				Status: db.StatusPending,
			}); err != nil {
				return err // Return error if marking fails
			}
			log.Logger.Debug("Marked as pending", zap.String("url", queue.Url.String))
		}
		return nil // Successfully processed all queues
	}))

	return errorList // Return the list of errors encountered
}
func SearchByWeightWithPagination(ctx context.Context, terms string, limit, offset int32) ([]db.Page, error) {
	// Separe os termos em palavras por espaço
	words := strings.Fields(terms)

	// Use strings.Builder para construir a query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, url, title, description, visited, timestamp")

	// Monte a cláusula WHERE
	queryBuilder.WriteString(" FROM pages WHERE ")
	for i, word := range words {
		if i > 0 {
			queryBuilder.WriteString(" OR ")
		}
		queryBuilder.WriteString(fmt.Sprintf("words ? '%s'", word))
	}

	// Monte a cláusula ORDER BY
	queryBuilder.WriteString(" ORDER BY ")
	for i, word := range words {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("COALESCE(words ->> '%s', '0')::int DESC", word))
	}

	// Adicione LIMIT e OFFSET
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset))

	// Execute a query
	query := queryBuilder.String()
	rows, err := Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapeie os resultados
	var pages []db.Page
	for rows.Next() {
		var page db.Page
		if err := rows.Scan(
			&page.ID,
			&page.Url,
			//&page.Links,
			&page.Title,
			&page.Description,
			//&page.Meta,
			&page.Visited,
			&page.Timestamp,
			//&page.Words,
		); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func SearchByWeight(ctx context.Context, terms string) ([]db.Page, error) {
	// Separe os termos em palavras por espaço
	words := strings.Fields(terms)

	// Use strings.Builder para construir a query
	var queryBuilder strings.Builder
	queryBuilder.WriteString("SELECT id, url, links, title, description, meta, visited, timestamp, words")

	// Monte a cláusula WHERE
	queryBuilder.WriteString(" FROM pages WHERE ")
	for i, word := range words {
		if i > 0 {
			queryBuilder.WriteString(" OR ")
		}
		queryBuilder.WriteString(fmt.Sprintf("words ? '%s'", word))
	}

	// Monte a cláusula ORDER BY
	queryBuilder.WriteString(" ORDER BY ")
	for i, word := range words {
		if i > 0 {
			queryBuilder.WriteString(", ")
		}
		queryBuilder.WriteString(fmt.Sprintf("COALESCE(words ->> '%s', '0')::int DESC", word))
	}

	// Execute a query
	query := queryBuilder.String()
	rows, err := Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Mapeie os resultados
	var pages []db.Page
	for rows.Next() {
		var page db.Page
		if err := rows.Scan(
			&page.ID,
			&page.Url,
			&page.Links,
			&page.Title,
			&page.Description,
			&page.Meta,
			&page.Visited,
			&page.Timestamp,
			&page.Words,
		); err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}
func AddHelperI2P(ctx context.Context, url string) (db.Queue, error) {
	queueEntry, err := DbQueries.QueueInsert(ctx, db.QueueInsertParams{
		Url:    url,
		Depth:  0,
		Status: db.StatusUnknown,
	})
	return queueEntry, err
}

func GetPagesForGraph(ctx context.Context) ([]db.GetAllPagesForGraphRow, error) {
	pages, err := DbQueries.GetAllPagesForGraph(ctx)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func GetQueue(ctx context.Context, limit, offset int32, strStatus string) ([]db.Queue, error) {
	status := db.Status(strStatus)
	if !status.Valid() {
		return nil, fmt.Errorf("invalid status")
	}

	queues, err := DbQueries.QueueGetAllByStatus(ctx, db.QueueGetAllByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return queues, nil
}
