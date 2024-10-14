# WebCrawler 3

## Descrição

Ser um crawler leve, altamente configurável, que seja possível rodar em qualquer computador relativamente antigo e que
também consiga rodar em paralelo com outros computadores, usando o mínimo de pacote de terceiros e compatível com I2P,
Tor, ZeroNet, Lokinet e Hyperboria.

## Tecnologias Utilizadas

- Go
- PostgreSQL
- TypeScript
- React

## Funcionalidades

- Funcionalidade 1: Descrição da funcionalidade 1.
- Funcionalidade 2: Descrição da funcionalidade 2.
- Funcionalidade 3: Descrição da funcionalidade 3.

## Requisitos

- Go 1.23
- Node.js 18.x ou superior
- PostgreSQL 16.x ou superior

## Instalação

1. Clone o repositório:
   ```sh
   git clone https://github.com/gabrielmoura/webcrawler3.git
   cd webcrawler3
   ```

2. Instale as dependências do backend:
   ```sh
   go mod tidy
   ```

3. Instale as dependências do frontend:
   ```sh
   make build-front
   ```
4. Compilação
   ```sh
   make build
   ```
5. Configure o banco de dados:
    - Crie o banco de dados.
    - Configure as variáveis de ambiente necessárias.

## Uso

1. Inicie o backend:
   ```sh
   ./wc3
   ```
2. Acesse o projeto no navegador:
   ```
   http://localhost:9090
   ```

## Informativo

O que este projeto não é:

- Um produto comercial.
- Um crawler profissional.

E por todos os motivos acima é importante que ao usar ou compartilhar concorda com os termos a seguir:

Este programa é distribuído como está e não cabe ao desenvolvedor a responsabilidade de mantê-lo (corrigir bugs,
melhorias), como também não cabe a responsabilidade pelo uso indevido ou que quebre as regras de outros sites.

Caso encontre uma maneira melhor de fazer as coisas envie um pull-request, caso tenha uma sugestão ou bug abra uma
issue.

## Contribuição

1. Faça um fork do projeto.
2. Crie uma branch para sua feature (`git checkout -b feature/nome-da-feature`).
3. Commit suas mudanças (`git commit -m 'Adiciona nova feature'`).
4. Faça o push para a branch (`git push origin feature/nome-da-feature`).
5. Abra um Pull Request.

## Licença

Este projeto está licenciado sob a Licença AGPL-3 - veja o arquivo `LICENSE` para mais detalhes.

## Contato

- Nome: Gabriel Moura
- Email: gmouradev96[at]gmail.com