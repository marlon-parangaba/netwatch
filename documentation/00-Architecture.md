# Arquitetura do Sistema

## Visão Geral da Arquitetura

O NetWatch segue uma arquitetura **cliente-servidor** com separação clara entre frontend e backend.

## Componentes Principais

### 1. Frontend - Desktop Application
- **Tecnologia**: Electron + React
- **Plataformas**: Windows, Linux, macOS
- **Responsabilidades**:
  - Interface do usuário
  - Visualização de dados
  - Editor de mapas
  - Dashboard
  - Sistema de notificações local

### 2. Backend - API Server
- **Tecnologia**: Go + Fiber
- **Plataforma**: Debian 13 (VM Proxmox)
- **Responsabilidades**:
  - API RESTful
  - WebSocket para tempo real
  - Coleta SNMP (polling ativo)
  - **Receptor de SNMP Traps**
  - Processamento de dados
  - **Gestão multi-usuário (auth, roles)**
  - Gestão de alertas
  - **Cálculo de métricas derivadas de alertas**
  - Armazenamento (PostgreSQL)

### 3. Banco de Dados
- **Tecnologia**: PostgreSQL
- **Responsabilidades**:
  - Dados de configuração
  - Métricas históricas
  - Alertas e eventos
  - Mapas e disposições

## Fluxo de Dados

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   Usuário    │────▶│   Frontend   │────▶│   Backend    │
│  (Desktop)   │◀────│  (Electron)  │◀────│    (Go)      │
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                                                 │
                            ┌────────────────────┼────────────────────┐
                            │                    │                    │
                            ▼                    ▼                    ▼
                     ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
                     │ PostgreSQL   │     │   Polling    │     │   Trap       │
                     │  (Dados)     │     │   SNMP       │     │   Receiver   │
                     └──────────────┘     └──────────────┘     └──────────────┘
                                                 │                    │
                                                 └────────┬───────────┘
                                                          ▼
                                                 ┌──────────────┐
                                                 │  Dispositivos│
                                                 │    de Rede   │
                                                 └──────────────┘
```

### Modos de Coleta SNMP

1. **Polling (Ativo)**: Backend faz requisições periódicas aos dispositivos
2. **Traps (Passivo)**: Dispositivos enviam notificações ao backend quando eventos ocorrem

## Retenção de Dados

| Tipo de Dado | Retention | Estratégia |
|--------------|-----------|------------|
| **Métricas de Tráfego** | 1 ano | Particionamento por mês |
| **Alertas** | 6 meses | Rolling window (sobreposição) |
| **Eventos de Downtime** | Vitalício | Para cálculo de métricas derivadas |

### Exemplo: Métricas Derivadas de Alertas

Quando um dispositivo fica offline (ping down):
1. Registrar evento `down` com timestamp
2. Quando voltar, registrar evento `up` com timestamp
3. Calcular duração do downtime
4. **Métricas derivadas**: Tempo de duração de baterias/UPS, MTBF, MTTR

## Estrutura de Diretórios (Backend)

```
netwatch-backend/
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── config/               # Configuração
│   ├── database/             # Conexão e migrations
│   ├── handlers/              # Handlers HTTP
│   ├── models/                # Modelos de dados
│   ├── repository/            # Acesso a dados
│   ├── services/              # Lógica de negócio
│   ├── snmp/                  # Cliente SNMP
│   └── websocket/             # WebSocket handlers
├── pkg/
│   └── utils/                 # Utilitários
├── migrations/                # Migrations SQL
├── config.yaml                 # Configuração
└── go.mod
```

## Estrutura de Diretórios (Frontend)

```
netwatch-frontend/
├── src/
│   ├── main/                  # Electron main process
│   ├── renderer/              # React app
│   │   ├── components/        # Componentes React
│   │   ├── pages/             # Páginas
│   │   ├── hooks/             # Custom hooks
│   │   ├── services/          # API calls
│   │   ├── stores/            # Estado global
│   │   └── styles/            # Estilos
│   └── preload/               # Preload scripts
├── electron-builder.yml       # Configuração de build
└── package.json
```

## API Endpoints (Planejado)

### Autenticação
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `POST /api/auth/refresh` - Refresh token

### Dispositivos
- `GET /api/devices` - Listar dispositivos
- `GET /api/devices/:id` - Detalhes do dispositivo
- `POST /api/devices` - Criar dispositivo
- `PUT /api/devices/:id` - Atualizar dispositivo
- `DELETE /api/devices/:id` - Remover dispositivo
- `POST /api/devices/discover` - Discovery de rede

### Métricas
- `GET /api/metrics/:deviceId` - Obter métricas
- `GET /api/metrics/:deviceId/history` - Histórico de métricas
- `POST /api/metrics/:deviceId/poll` - Forçar coleta

### Alertas
- `GET /api/alerts` - Listar alertas
- `POST /api/alerts` - Criar regra de alerta
- `PUT /api/alerts/:id` - Atualizar regra
- `GET /api/alerts/history` - Histórico de alertas

### Mapas
- `GET /api/maps` - Listar mapas
- `GET /api/maps/:id` - Obter mapa
- `POST /api/maps` - Criar mapa
- `PUT /api/maps/:id` - Atualizar mapa
- `PUT /api/maps/:id/positions` - Atualizar posições

### Dashboard
- `GET /api/dashboard/stats` - Estatísticas gerais
- `GET /api/dashboard/widgets` - Widgets configurados

## WebSocket Events

- `device:status` - Mudança de status de dispositivo
- `metrics:update` - Novas métricas
- `alert:triggered` - Alerta disparado
- `map:update` - Atualização de mapa

## Segurança

- **Autenticação multi-usuário** com JWT
- Roles: Admin, Operator, Viewer
- Senha hasheada com bcrypt
- HTTPS obrigatório para API
- Rate limiting
- Validação de entrada

---

*Última atualização: 2026-03-20*
