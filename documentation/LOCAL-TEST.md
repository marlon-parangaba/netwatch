# Como Testar o Backend Localmente

## Opção 1: Docker Compose (Recomendado)

### Pré-requisitos
- Docker Desktop instalado no Windows

### Passos

```bash
# 1. Na raiz do projeto, copie o .env.example para .env
cp .env.example .env

# 2. Inicie o PostgreSQL e Backend
docker-compose up -d postgres netwatch-backend

# 3. (Opcional) Inicie o Adminer para visualizar o banco
docker-compose --profile dev up -d adminer
```

**Acesso:**
- API: http://localhost:8080
- Health Check: http://localhost:8080/health
- Adminer: http://localhost:8081

### Verificar se está rodando
```bash
docker-compose ps
curl http://localhost:8080/health
```

### Parar
```bash
docker-compose down
```

---

## Opção 2: Go Direto (Sem Docker)

### Pré-requisitos
- Go 1.23+ instalado
- PostgreSQL 16 instalado localmente (ou via Docker)

### Passos

```bash
# 1. Crie o banco de dados no PostgreSQL
psql -U postgres -c "CREATE DATABASE netwatch;"

# 2. Entre no diretório do backend
cd netwatch-backend

# 3. Copie config example e ajuste
cp config.yaml.example config.yaml
# Edite o config.yaml se necessário

# 4. Baixe as dependências
go mod download

# 5. Execute
go run cmd/server/main.go
```

### Testar endpoints
```bash
# Health check
curl http://localhost:8080/health

# Login (depois de criar usuário)
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@netwatch.com","password":"admin123"}'
```

---

## Criar Usuário Admin Inicial

O backend deve criar um usuário admin na primeira execução, ou você pode inserir manualmente:

```sql
-- Via psql ou Adminer (http://localhost:8081)
INSERT INTO users (name, email, password_hash, role, active)
VALUES (
  'Admin',
  'admin@netwatch.com',
  -- Hash bcrypt de 'admin123' (você pode gerar um novo)
  '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
  'admin',
  true
);
```

Ou gere o hash no Go:
```bash
cd netwatch-backend
go run -e 'package main; import ("fmt"; "golang.org/x/crypto/bcrypt"); func main() { h, _ := bcrypt.GenerateFromPassword([]byte("admin123"), 10); fmt.Println(string(h)) }'
```

---

## Endpoints para Testar

### Autenticação
```bash
# Login
POST http://localhost:8080/api/auth/login
Body: {"email":"admin@netwatch.com","password":"admin123"}

# Refresh token
POST http://localhost:8080/api/auth/refresh
Body: {"refresh_token":"..."}
```

### Dispositivos
```bash
# Listar (precisa token)
GET http://localhost:8080/api/devices
Header: Authorization: Bearer <token>

# Criar
POST http://localhost:8080/api/devices
Body: {
  "name": "Roteador Teste",
  "ip_address": "192.168.1.1",
  "snmp_community": "public",
  "type": "mikrotik"
}
```

### Discovery
```bash
# Iniciar discovery
POST http://localhost:8080/api/devices/discover
Body: {"cidr":"192.168.1.0/24","snmp_community":"public"}
```

---

## Troubleshooting

### PostgreSQL não conecta
```bash
# Verifique se o PostgreSQL está rodando
# No Windows, inicie o serviço PostgreSQL
```

### Porta já em uso
```bash
# Verificar qual processo está usando a porta
netstat -ano | findstr :8080
netstat -ano | findstr :5432

# No .env, mude as portas se necessário
API_PORT=8081
POSTGRES_PORT=5433
```

### Logs do container
```bash
docker-compose logs -f netwatch-backend
```

---

## Comandos Úteis

```bash
# Rebuild do backend (após mudanças)
docker-compose build netwatch-backend
docker-compose up -d netwatch-backend

# Ver logs
docker-compose logs -f

# Reset completo (remove banco)
docker-compose down -v
docker-compose up -d
```
