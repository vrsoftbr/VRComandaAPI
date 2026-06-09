# VRComandaAPI

API para operação de comandas. Roda em container Docker e se comunica com o MongoDB do VRIntegrador e um banco SQLite local para persistência de lançamentos.

## Visão Geral

- **Framework:** Go + Gin
- **Banco local:** SQLite (lançamentos)
- **Banco compartilhado:** MongoDB via VRIntegrador (`pdv`)
- **Porta:** `28232`

## Pré-requisitos

O VRIntegradorMaster precisa estar rodando antes de subir a VRComandaAPI, pois é ele que fornece o MongoDB com os dados do PDV.

Para subir o VRIntegrador localmente:

```bash
cd ../VRIntegradorMaster
docker compose -f docker-compose.local.yml up -d
```

Isso vai criar os containers e a rede Docker `vrintegradormaster_vrintegrador-net`, que a VRComandaAPI usa para se comunicar com o MongoDB.

## Subindo com Docker

Você deve estar na `main` dos dois projetos antes de rodar.

```bash
docker compose up -d --build
```

### Verificando

```bash
docker logs vrcomandaapi --tail 20
```

A inicialização correta aparece assim (sem erros de MongoDB):

```
INFO VRComandaAPI starting port=:28232
INFO request method=GET path=/health status=200 ...
```

## Ajustes necessários por ambiente

Os valores do MongoDB estão hardcoded no `docker-compose.yml` (não como variável de ambiente, porque o `.env` sobrepõe os valores quando usado via Docker). Na prática é pra precisar ajustar só a porta mesmo, porque o resto já tá certinho.

### 1. Porta do MongoDB

Abra o `docker-compose.yml` e confira a porta na `MONGO_URI`:

```yaml
MONGO_URI: "mongodb://root:example@vrintegrador-db:26017/..."
```

Para confirmar qual porta o seu MongoDB está usando:

```bash
docker ps
```

A porta aparece no formato `26017->26017/tcp`. O número da **esquerda** é o que vai na URI. Se for diferente de `26017`, ajuste no `docker-compose.yml`.

### 2. Rede Docker

A VRComandaAPI se conecta à rede externa criada pelo VRIntegradorMaster. O nome dessa rede é gerado a partir do nome da pasta onde o VRIntegradorMaster foi clonado. Para confirmar o nome da rede no seu ambiente:

```bash
docker network ls
```

Procure algo que termine com `_vrintegrador-net`. Se for diferente de `vrintegradormaster_vrintegrador-net`, ajuste no `docker-compose.yml`:

```yaml
networks:
  vr_network:
    external: true
    name: vrintegradormaster_vrintegrador-net  # <- ajustar aqui se necessário
```

> O host `vrintegrador-db` não precisa ser alterado — ele está fixo como `container_name` no `docker-compose.local.yml` do VRIntegradorMaster, então enquanto todo mundo usar esse compose ele nunca muda.

## Conectando o Front (VRComanda)

No campo **"Endereço IP do servidor"** do front, preencha:

```
http://localhost:28232
```

## Testando

**Health check:**

```bash
curl http://localhost:28232/health
```

**Lojas cadastradas:**

```bash
curl http://localhost:28232/api/v1/lojas
```

**Swagger UI:**

```
http://localhost:28232/swagger/index.html
```

## Configuração do MongoDB (referência)

| Container         | Porta | Usuário | Senha     | Banco |
| ----------------- | ----- | ------- | --------- | ----- |
| `vrintegrador-db` | 26017 | `root`  | `example` | `pdv` |

## Variáveis de ambiente (.env)

O arquivo `.env` é usado apenas para rodar a API **fora do Docker** (desenvolvimento local direto com `go run .`). Copie o exemplo:

```bash
cp .env.example .env
```

| Variável         | Valor                                                                            | Descrição                         |
| ---------------- | -------------------------------------------------------------------------------- | --------------------------------- |
| `HTTP_PORT`      | `:28232`                                                                         | Porta HTTP da API                 |
| `MONGO_URI`      | `mongodb://root:example@localhost:26017/?directConnection=true&authSource=admin` | URI para rodar fora do Docker     |
| `MONGO_DATABASE` | `pdv`                                                                            | Nome do banco de dados no MongoDB |
| `SQLITE_PATH`    | `./data/pdv.db`                                                                  | Caminho do arquivo SQLite         |
