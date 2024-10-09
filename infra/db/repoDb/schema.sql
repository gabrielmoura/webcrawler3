-- Criação da extensão para geração de UUIDs, se ainda não existir
CREATE
EXTENSION IF NOT EXISTS "uuid-ossp";

-- Criação da tabela de páginas
CREATE TABLE IF NOT EXISTS pages
(
    id          uuid                     DEFAULT uuid_generate_v4() PRIMARY KEY, -- Identificador único
    url         TEXT,                                                            -- URL da página
    links       TEXT[],                                                          -- Links contidos na página
    title       TEXT,                                                            -- Título da página
    description TEXT,                                                            -- Descrição da página
    meta        JSONB,                                                           -- Metadados em formato JSON
    visited     BOOLEAN                  DEFAULT FALSE,                          -- Indica se a página foi visitada
    timestamp   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),                          -- Timestamp da criação
    words       JSONB                                                            -- Palavras contidas na página
);

-- Índices para otimização de consultas na tabela pages
CREATE INDEX IF NOT EXISTS idx_words_gin ON pages USING GIN (words);
CREATE INDEX IF NOT EXISTS idx_url ON pages (url);

-- Criação do tipo ENUM para o status da fila
CREATE TYPE status AS ENUM (
    'unknown',
    'pending',
    'visited',
    'error',
    'skipped',
    'processing_out',
    'processing'
);
-- Criação da tabela de fila
CREATE TABLE IF NOT EXISTS queue
(
    id        uuid                     DEFAULT uuid_generate_v4() PRIMARY KEY, -- Identificador único
    url       TEXT,                                                            -- URL da fila
    depth     INTEGER,                                                         -- Profundidade da URL na fila
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),                          -- Timestamp da criação
    status    status NOT NULL,                                                 -- Status da URL na fila
    delay     INTEGER DEFAULT 0,                                               -- Delay para processamento da URL
    error     JSONB                                                            -- Erro ocorrido no processamento da URL
);

-- Índices para otimização de consultas na tabela queue
CREATE INDEX IF NOT EXISTS idx_url_queue ON queue (url);
CREATE INDEX IF NOT EXISTS idx_status_queue ON queue (status);
