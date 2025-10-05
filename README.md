# poc-api-cosmos-dynamo-go

API de exemplo em Go para integração com Azure CosmosDB e DynamoDB.

## Como executar CosmosDB

O fluxo recomendado é rodar tudo via Docker Compose:

```sh
docker-compose up --build
```

Isso irá:
- Subir o CosmosDB emulador
- Executar o seed automaticamente (preenchendo o arquivo `uuids.txt` na raiz do projeto)
- Subir a API Go em `http://localhost:8888`
- Subir o Jaeger com a IU acessível em `http://localhost:16686`

## Como executar DynamoDB

O fluxo recomendado é rodar tudo via Docker Compose:

```sh
docker-compose -f docker-compose.dynamo.yml up --build
```

Isso irá:
- Subir o DynamoDB local
- Executar o seed automaticamente (preenchendo o arquivo `uuids.txt` na raiz do projeto)
- Subir a API Go em `http://localhost:8888`
- Subir o Jaeger com a IU acessível em `http://localhost:16686`

## Pré-requisitos

- Docker

## Endpoints

- `GET /health` — Healthcheck
- `GET /items/:id` — Busca um relacionamento pelo ID (retorna 204 se não encontrado)

## Observações

- O projeto ignora a validação de certificados TLS para facilitar testes locais. Não usar em produção!
- O campo `id` não é retornado pela API, apenas os dados de relacionamento.
- Retorno 204 caso o item não seja encontrado.
