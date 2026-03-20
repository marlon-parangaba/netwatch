# PROMPT: Implementar Fase 3 - Frontend Desktop (NetWatch)

## 1. Visão Geral do Projeto

O **NetWatch** é uma ferramenta de monitoramento SNMP para ISPs, mesclando funcionalidades do Zabbix e The Dude da Mikrotik.

- **Repositório**: https://github.com/marlon-parangaba/netwatch
- **Backend**: Go + Fiber + PostgreSQL (já implementado)
- **Frontend**: Electron + React + TypeScript (a implementar)
- **Plataformas**: Windows, Linux, macOS

## 2. Stack Tecnológica Confirmada

```json
{
  "framework": "Electron + React",
  "language": "TypeScript",
  "state_management": "Zustand",
  "routing": "React Router v6",
  "http_client": "Axios",
  "styling": "Tailwind CSS",
  "charts": "Recharts",
  "maps": "React Flow (@xyflow/react)",
  "build": "electron-builder"
}
```

## 3. Estrutura do Backend (para integração)

### Estrutura de diretórios do backend
```
netwatch-backend/
├── cmd/server/main.go
├── internal/
│   ├── config/config.go
│   ├── database/
│   ├── handlers/
│   │   ├── auth_handler.go      # POST /api/auth/login, /logout, /refresh
│   │   ├── device_handler.go    # CRUD /api/devices
│   │   ├── discovery_handler.go # POST /api/devices/discover
│   │   └── metric_handler.go   # GET /api/metrics/:deviceId
│   ├── models/                  # User, Device, Metric, Alert
│   ├── repository/
│   ├── services/
│   │   ├── auth_service.go
│   │   ├── device_service.go
│   │   └── poller_service.go   # SNMP polling + traps
│   └── snmp/
│       ├── client.go
│       ├── discovery.go
│       ├── poller.go
│       └── trap_receiver.go
├── config.yaml
└── go.mod
```

### Modelos de Dados (Schema PostgreSQL)

#### User
```typescript
interface User {
  id: string;           // UUID
  name: string;
  email: string;
  role: 'admin' | 'operator' | 'viewer';
  active: boolean;
  last_login_at: string; // ISO timestamp
  created_at: string;
  updated_at: string;
}
```

#### Device
```typescript
interface Device {
  id: string;
  name: string;
  hostname: string | null;
  ip_address: string;
  type: 'mikrotik' | 'cisco' | 'huawei' | 'juniper' | 'ubiquiti' | 'generic';
  status: 'online' | 'offline' | 'warning' | 'unknown';
  description: string | null;
  location: string | null;
  
  // SNMP Config
  snmp_version: 'v1' | 'v2c' | 'v3';
  snmp_community: string;
  snmp_port: number;
  snmpv3_username?: string;
  snmpv3_auth_proto?: string;
  snmpv3_auth_key?: string;
  snmpv3_priv_proto?: string;
  snmpv3_priv_key?: string;
  
  // Polling
  polling_interval: number; // segundos
  enabled: boolean;
  
  // Sys Info (from SNMP)
  sys_name?: string;
  sys_descr?: string;
  sys_oid?: string;
  sys_contact?: string;
  
  // Metadados
  tags: string[];
  last_seen_at: string | null;
  last_polled_at: string | null;
  created_at: string;
  updated_at: string;
}
```

#### Metric
```typescript
interface Metric {
  id: string;
  device_id: string;
  type: string;        // 'cpu', 'memory', 'temperature', 'interface_in', 'interface_out', etc.
  interface_index: number;
  oid: string;
  value: number;
  unit: string;        // '%', 'bytes', '°C', etc.
  collected_at: string;
}
```

#### AlertRule
```typescript
interface AlertRule {
  id: string;
  name: string;
  description: string | null;
  device_id: string | null; // null = todas
  metric_type: string;
  condition: 'gt' | 'gte' | 'lt' | 'lte' | 'eq' | 'neq' | 'down';
  threshold: number;
  severity: 'critical' | 'high' | 'medium' | 'low' | 'info';
  enabled: boolean;
  notifications: Array<{
    type: 'email' | 'telegram' | 'webhook';
    config: Record<string, string>;
  }>;
  created_at: string;
  updated_at: string;
}
```

#### AlertEvent
```typescript
interface AlertEvent {
  id: string;
  rule_id: string;
  device_id: string;
  status: 'active' | 'resolved' | 'acknowledged';
  severity: string;
  value: number | null;
  message: string;
  triggered_at: string;
  resolved_at: string | null;
  acked_at: string | null;
  acked_by_id: string | null;
  duration_seconds: number | null;
}
```

#### DowntimeEvent
```typescript
interface DowntimeEvent {
  id: string;
  device_id: string;
  started_at: string;
  ended_at: string | null;
  duration_seconds: number | null;
  cause: string | null;
  notes: string | null;
}
```

### API Endpoints

#### Autenticação
- `POST /api/auth/login` — Body: `{ email, password }` → Returns `{ token, refresh_token, user }`
- `POST /api/auth/refresh` — Body: `{ refresh_token }` → Returns `{ token }`
- `POST /api/auth/logout` — Invalida token

#### Dispositivos
- `GET /api/devices` — Lista dispositivos (com filtros)
- `GET /api/devices/:id` — Detalhes de um dispositivo
- `POST /api/devices` — Criar dispositivo
- `PUT /api/devices/:id` — Atualizar dispositivo
- `DELETE /api/devices/:id` — Remover dispositivo
- `POST /api/devices/discover` — Body: `{ cidr, snmp_community }` → Inicia discovery
- `POST /api/devices/discover/import` — Body: `{ ips: [] }` → Importa IPs descobertos

#### Métricas
- `GET /api/metrics/:deviceId` — Métricas atuais
- `GET /api/metrics/:deviceId/history?type=cpu&from=&to=` — Histórico
- `POST /api/metrics/:deviceId/poll` — Força coleta

#### Alertas
- `GET /api/alerts/rules` — Lista regras
- `POST /api/alerts/rules` — Criar regra
- `PUT /api/alerts/rules/:id` — Atualizar regra
- `DELETE /api/alerts/rules/:id` — Remover regra
- `GET /api/alerts/events` — Lista eventos
- `PUT /api/alerts/events/:id/ack` — Acknowledgement

#### Dashboard
- `GET /api/dashboard/stats` — Estatísticas gerais
- `GET /api/dashboard/top-devices` — Top dispositivos por status/métricas

## 4. Requisitos do Frontend

### 4.1 Autenticação
- [ ] Tela de Login (email + senha)
- [ ] Armazenar JWT no localStorage/electron-store
- [ ] Interceptor Axios para adicionar Authorization header
- [ ] Protected routes (redirecionar para login se não autenticado)
- [ ] Logout (limpa tokens, redireciona para login)
- [ ] Exibir nome do usuário e role no header

### 4.2 Layout Principal
- [ ] Sidebar navegável (collapse/expand)
- [ ] Header com: logo, user info, logout
- [ ] Área de conteúdo principal
- [ ] Temas: dark mode (padrão) + light mode

### 4.3 Dashboard
- [ ] Cards de estatísticas: total devices, online, offline, warnings
- [ ] Gráfico de status dos dispositivos (pizza/donut chart)
- [ ] Últimos alertas disparados
- [ ] Dispositivos recently seen / recently down
- [ ] Gráfico de uso de CPU/memória agregado (últimas 24h)

### 4.4 Gerenciamento de Dispositivos
- [ ] Lista de dispositivos com filtros (status, tipo, busca por nome/IP)
- [ ] Tabela com colunas: Nome, IP, Tipo, Status, Última Coleta, Ações
- [ ] Paginação
- [ ] Modal de criação/edição de dispositivo
- [ ] Teste de conectividade SNMP antes de salvar
- [ ] Exclusão com confirmação
- [ ] Botão para forçar poll de um dispositivo

### 4.5 Detalhes de Dispositivo
- [ ] Informações gerais (sysInfo, SNMP config)
- [ ] Gráficos em tempo real (CPU, Memória, Temperature)
- [ ] Lista de interfaces de rede com status
- [ ] Gráfico de tráfego de interface (in/out bytes)
- [ ] Histórico de alertas do dispositivo
- [ ] Timeline de eventos de downtime

### 4.6 Discovery de Rede
- [ ] Input para CIDR (ex: 192.168.1.0/24)
- [ ] Community string
- [ ] Progress bar durante varredura
- [ ] Lista de IPs encontrados com sysInfo preview
- [ ] Seleção múltipla de dispositivos para importar
- [ ] Botão "Importar Selecionados"

### 4.7 Sistema de Alertas
- [ ] Lista de regras de alerta
- [ ] CRUD de regras com UI intuitiva
- [ ] Condição visual: "CPU > 80%" etc
- [ ] Configuração de notificações (email, telegram webhook)
- [ ] Lista de eventos de alerta (ativos, resolvidos)
- [ ] Botão acknowledge em alertas
- [ ] Filtros: severidade, status, dispositivo

### 4.8 Mapas de Rede
- [ ] Canvas com React Flow
- [ ] Adicionar dispositivos ao mapa (drag from list)
- [ ] Layout automático (opcional)
- [ ] Status visual dos nós (verde/vermelho/amarelo)
- [ ] Zoom e pan
- [ ] Salvar/carregar layouts
- [ ] Múltiplos mapas

### 4.9 Configurações
- [ ] Gerenciamento de usuários (CRUD) — apenas admin
- [ ] Configurações gerais do sistema

## 5. Estrutura de Diretórios Sugerida

```
netwatch-frontend/
├── electron/
│   ├── main.ts              # Electron main process
│   ├── preload.ts           # Preload script
│   └── ipc-handlers.ts      # IPC communication
├── src/
│   ├── main.tsx             # React entry
│   ├── App.tsx              # Root component
│   ├── components/
│   │   ├── common/          # Button, Input, Card, Modal, Table, etc.
│   │   ├── layout/          # Sidebar, Header, MainLayout
│   │   ├── dashboard/       # Dashboard widgets
│   │   ├── devices/         # Device components
│   │   ├── alerts/          # Alert components
│   │   └── maps/            # Map components
│   ├── pages/
│   │   ├── Login.tsx
│   │   ├── Dashboard.tsx
│   │   ├── Devices.tsx
│   │   ├── DeviceDetail.tsx
│   │   ├── Discovery.tsx
│   │   ├── Alerts.tsx
│   │   ├── Maps.tsx
│   │   └── Settings.tsx
│   ├── hooks/
│   │   ├── useAuth.ts
│   │   ├── useDevices.ts
│   │   ├── useMetrics.ts
│   │   └── useAlerts.ts
│   ├── services/
│   │   ├── api.ts           # Axios instance
│   │   ├── auth.service.ts
│   │   ├── device.service.ts
│   │   ├── metric.service.ts
│   │   └── alert.service.ts
│   ├── stores/
│   │   ├── authStore.ts     # Zustand
│   │   ├── deviceStore.ts
│   │   └── uiStore.ts
│   ├── types/
│   │   └── index.ts         # TypeScript interfaces
│   └── styles/
│       └── globals.css      # Tailwind imports
├── package.json
├── electron-builder.yml
├── tailwind.config.js
├── vite.config.ts
└── tsconfig.json
```

## 6. Convenções de Código

### TypeScript
- Usar interfaces para tipos de dados (não types para objetos)
- Nomes: camelCase para variáveis, PascalCase para componentes
- Exports: named exports preferidos
- Nullability: usar `undefined` ao invés de `null` quando possível

### React
- Componentes funcionais com hooks
- Props com interfaces tipadas
- Extrair lógica para custom hooks quando reutilizável
- Componentes menores e mais focados

### CSS/Tailwind
- Mobile-first responsive design
- Usar classes do Tailwind, evitar CSS customizado
- Variáveis CSS para cores do tema (dark/light)

## 7. Configuração de Build (electron-builder)

```yaml
appId: com.netwatch.app
productName: NetWatch
directories:
  output: release
files:
  - build
  - package.json
win:
  target:
    - nsis
  icon: build/icon.ico
mac:
  target:
    - dmg
  icon: build/icon.icns
linux:
  target:
    - AppImage
    - deb
  icon: build/icon.png
```

## 8. Próximas Fases (para contexto)

- **Fase 4**: Mapas de Rede (React Flow, editor visual)
- **Fase 5**: Sistema de Alertas avançado (notificações, escalação, métricas derivadas)
- **Fase 6**: Polimento UX
- **Fase 7**: Release

## 9. Observações Importantes

1. **Backend rodando**: O backend deve estar configurado para aceitar conexões do frontend (CORS configurado)
2. **Tempo real**: Planejar WebSocket para atualizações em tempo real (futuro)
3. **Performance**: Lazy loading de páginas, memoização onde necessário
4. **Responsividade**: UI deve funcionar bem em telas de varios tamanhos

---

*Prompt gerado para implementação do Frontend NetWatch*
