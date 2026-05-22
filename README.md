# VRComandaAPI

API para operação de comandas. Roda em container Docker e se comunica com o MongoDB do VRIntegrador e um banco SQLite local para persistência de lançamentos.

## Visão Geral

- **Framework:** Go + Gin- **Banco local:** SQLite (lançamentos)
- **Banco compartilhado:** MongoDB via VRIntegrador
- **Porta:** `28232` ( verificar qual porta devemos usar quando for pra prod ou entrar em esteira de deploy. )

## Configuração

### 1. Configure as variáveis de ambiente

Copie o arquivo de exemplo e ajuste se necessário:

```bash
.env.example .env
```

Valores padrão do `.env`:

| Variável         | Valor padrão                                                                  | Descrição                       |
| ---------------- | ----------------------------------------------------------------------------- | ------------------------------- |
| `HTTP_PORT`      | `:28232`                                                                      | Porta HTTP da API               |
| `MONGO_URI`      | `mongodb://root:root@localhost:27017/?directConnection=true&authSource=admin` | URI de conexão com o MongoDB    |
| `MONGO_DATABASE` | `vrcomanda`                                                                   | Nome do banco de dados no Mongo |
| `SQLITE_PATH`    | `./data/vrcomanda.db`                                                         | Caminho do arquivo SQLite       |

> No container, a `MONGO_URI` aponta para `vrintegrador-db` (hostname interno do Docker, verifique o seu se deve se manter ou não). O `.env` usa `localhost` pra rodar fora do docker.

## Subindo com Docker

### 1. Suba o VRIntegrador (MongoDB)

O MongoDB é provido pelo VRIntegrador. Ele precisa estar rodando antes da VRComandaAPI. É a partir dele que a API vai ler e escrever os dados compartilhados com o PDV.;

### 2. Suba a VRComandaAPI

```bash
docker compose up -d --build
```

### 3. Verifique se subiu

```bash
docker ps
docker logs vrcomandaapi
```

## Testando

**Health check:**

```bash
curl http://localhost:28232/health
```

Resposta esperada: `200 OK`

**Swagger UI:**

```
http://localhost:28232/swagger/index.html
```

##### Rede Docker

O container se conecta à rede externa `vrintegradormaster_vrintegrador-net`, criada pelo VRIntegrador. Essa rede precisa existir antes de subir a VRComandaAPI. Deposi podemos trocar por uma variável de ambiente ou algo mais flexível. Se der erro, verifique o nome dessa rede no seu ambiente e ajuste o `docker-compose.yml` ou crie a rede manualmente:
