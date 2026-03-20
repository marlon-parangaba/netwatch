# Roadmap de Desenvolvimento

## Visão Geral

O desenvolvimento será dividido em **fases iterativas**, com release early e often.

## Fases do Projeto

### ✅ Fase 0: Fundação (Semanas 1-2) — CONCLUÍDA
> Setup inicial e arquitetura base

**Objetivos:**
- [x] Setup repositório Git
- [x] Configurar estrutura de diretórios
- [x] Criar documentação inicial (esta documentação)
- [x] Setup ambiente de desenvolvimento
- [x] Criar Dockerfile/docker-compose para backend

**Entregáveis:**
- ✅ Repositório com estrutura base
- ✅ Documentação completa
- ✅ Ambiente Docker para desenvolvimento

**Implementado:**
- `docker-compose.yml` — PostgreSQL 16 + Backend + Adminer
- `netwatch-backend/Dockerfile` — Build multi-stage (Go → Alpine)
- `.env.example` — Variáveis documentadas

**Status:** 🟢 Concluída

---

### ✅ Fase 1: Backend Core (Semanas 3-5) — CONCLUÍDA
> API básica e conexão com banco de dados

**Objetivos:**
- [x] Configurar PostgreSQL com Docker
- [x] Setup GORM e criar modelos base
- [x] Implementar migrações
- [x] Criar API REST de dispositivos (CRUD)
- [x] Implementar autenticação JWT
- [x] Setup logging e error handling

**Entregáveis:**
- ✅ API REST funcional
- ✅ Autenticação funcionando (JWT + refresh)
- ✅ Migrations executando
- [ ] Swagger/OpenAPI documentado

**Implementado:**
- `internal/models/` — User, Device, Metric, Alert
- `internal/repository/` — User, Device, Metric (interfaces limpas)
- `internal/services/` — AuthService (JWT), DeviceService (CRUD)
- `internal/handlers/` — Auth, Device, Metric + middleware JWT/RBAC
- `cmd/server/main.go` — Graceful shutdown, CORS, rate limiting, health check
- `config.yaml` + `internal/config/` — Suporte YAML + env vars
- `migrations/001_initial_schema.up.sql` — Schema completo com:
  - Particionamento de métricas por mês (retenção 1 ano)
  - `downtime_events` vitalício (MTBF/MTTR)
  - `alert_events` com rolling window de 6 meses
  - Triggers para `updated_at` automático

**Status:** 🟢 Concluída (pendendo Swagger)

---

### ✅ Fase 2: SNMP Engine (Semanas 6-8) — CONCLUÍDA
> Comunicação SNMP com dispositivos

**Objetivos:**
- [x] Implementar cliente SNMP com gosnmp
- [x] Criar serviço de discovery de rede
- [x] Implementar coleta de métricas (polling ativo)
- [x] **Implementar receptor de SNMP Traps**
- [x] Criar polling schedule
- [x] Suportar SNMP v1, v2c e v3

**Entregáveis:**
- ✅ Módulo SNMP funcional
- ✅ Discovery automático (CIDR, worker pool, detecção de vendor)
- ✅ Coleta de métricas
- ✅ Receptor de Traps funcional

**Implementado:**
- `internal/snmp/oids.go` — OIDs MIB-II, IF-MIB, HOST-RESOURCES, UCD-SNMP, Mikrotik
- `internal/snmp/client.go` — Wrapper gosnmp (v1/v2c/v3), GetSysInfo, GetInterfaces, Walk
- `internal/snmp/discovery.go` — Varredura CIDR concorrente, detecção de vendor por sysOID
- `internal/snmp/poller.go` — Poller (device) + PollerScheduler (todos os devices)
- `internal/snmp/trap_receiver.go` — UDP listener, classifica coldStart, warmStart, linkDown, linkUp, authFailure
- `internal/services/poller_service.go` — Orquestra Scheduler + TrapReceiver
- `internal/handlers/discovery_handler.go` — POST /api/devices/discover, /api/devices/discover/import

**Status:** 🟢 Concluída

---

### 🔲 Fase 3: Frontend Base (Semanas 9-12)
> Interface desktop básica

**Objetivos:**
- [ ] Setup Electron + React + TypeScript
- [ ] Implementar sistema de autenticação no frontend
- [ ] Criar dashboard básico
- [ ] Listar e gerenciar dispositivos
- [ ] Visualizar métricas de dispositivo
- [ ] Configurar electron-builder para builds

**Entregáveis:**
- App desktop executável
- Login/logout funcionando
- Dashboard com widgets básicos

**Responsável:** A definir
**Status:** 🔴 Pendente

---

### 🔲 Fase 4: Mapas de Rede (Semanas 13-16)
> Editor de mapas como The Dude

**Objetivos:**
- [ ] Implementar canvas com React Flow
- [ ] Criar editor de mapa (drag & drop)
- [ ] Mostrar status dos dispositivos visualmente
- [ ] Permitir customização de ícones
- [ ] Salvar/carregar layouts de mapa
- [ ] Zoom e pan automáticos

**Entregáveis:**
- Editor de mapas funcional
- Mapas saváveis
- Atualização em tempo real

**Responsável:** A definir
**Status:** 🔴 Pendente

---

### 🔲 Fase 5: Sistema de Alertas (Semanas 17-19)
> Notificações, triggers e métricas derivadas

**Objetivos:**
- [ ] Criar modelo de regras de alerta
- [ ] Implementar motor de avaliação de alertas
- [ ] **Implementar retenção de 6 meses com sobreposição (rolling window)**
- [ ] Suporte a múltiplos canais de notificação
- [ ] Implementar escalação de alertas
- [ ] Criar UI de configuração de alertas
- [ ] **Criar métricas derivadas de alertas** (ex: duração de baterias baseada em eventos "ping down")

**Entregáveis:**
- Motor de alertas funcionando
- Retenção e sobreposição de alertas
- UI de configuração
- Notificações funcionando
- **Métricas derivadas calculadas**

**Responsável:** A definir
**Status:** 🔴 Pendente

---

### 🔲 Fase 6: Polimento e UX (Semanas 20-22)
> Melhorias de interface e experiência

**Objetivos:**
- [ ] Implementar tema dark/light
- [ ] Melhorar responsividade
- [ ] Adicionar animações e transições
- [ ] Implementar atalhos de teclado
- [ ] Adicionar onboarding para novos usuários
- [ ] Otimizar performance

**Entregáveis:**
- UI polida
- Performance otimizada
- Documentação de usuário

**Responsável:** A definir
**Status:** 🔴 Pendente

---

### 🔲 Fase 7: Release v1.0 (Semanas 23-24)
> Preparação para lançamento

**Objetivos:**
- [ ] Testes de regressão
- [ ] Builds para Windows, Linux, macOS
- [ ] Setup de update automático
- [ ] Criar guia de instalação
- [ ] Versãoing semântico (v1.0.0)
- [ ] Release notes

**Entregáveis:**
- Builds finais
- Instaladores para todas plataformas
- Guia de instalação

**Responsável:** A definir
**Status:** 🔴 Pendente

---

## Timeline Visual

```
2026
     Mar          Abr          Mai          Jun          Jul          Ago
     ▼            ▼            ▼            ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│ Fase 0  │ │ Fase 1  │ │ Fase 2  │ │ Fase 3  │ │ Fase 4  │ │ Fase 5  │
│Fun-dação│ │ Backend │ │  SNMP   │ │Front-end│ │ Mapas   │ │ Alertas │
│  2 sem  │ │  3 sem  │ │  3 sem  │ │  4 sem  │ │  4 sem  │ │  3 sem  │
└─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
  ✅ ✅ ✅                                             

     Set          Out          Nov
     ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Fase 6  │ │ Fase 7  │ │ v1.0.0! │
│Polimento│ │ Release │
│  3 sem  │ │  2 sem  │
└─────────┘ └─────────┘ └─────────┘
```

## Milestones

| Versão | Descrição | Status | Data |
|--------|-----------|--------|------|
| v0.1.0 | Backend + API base | 🟢 Concluída | 2026-03-20 |
| v0.2.0 | SNMP funcional | 🟢 Concluída | 2026-03-20 |
| v0.3.0 | Frontend básico | 🔄 Em progresso | — |
| v0.4.0 | Mapas de rede | 🔴 Pendente | — |
| v0.5.0 | Sistema de alertas | 🔴 Pendente | — |
| v1.0.0 | Release production-ready | 🔴 Pendente | — |

---

*Última atualização: 2026-03-20*
