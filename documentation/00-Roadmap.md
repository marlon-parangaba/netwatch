# Roadmap de Desenvolvimento

## Visão Geral

O desenvolvimento será dividido em **fases iterativas**, com release early e often.

## Fases do Projeto

### 🔲 Fase 0: Fundação (Semanas 1-2)
> Setup inicial e arquitetura base

**Objetivos:**
- [ ] Setup repositório Git
- [ ] Configurar estrutura de diretórios
- [ ] Criar documentação inicial (esta documentação)
- [ ] Setup ambiente de desenvolvimento
- [ ] Criar Dockerfile/docker-compose para backend

**Entregáveis:**
- Repositório com estrutura base
- Documentação completa
- Ambiente Docker para desenvolvimento

**Responsável:** A definir
**Status:** 🟡 A fazer

---

### 🔲 Fase 1: Backend Core (Semanas 3-5)
> API básica e conexão com banco de dados

**Objetivos:**
- [ ] Configurar PostgreSQL com Docker
- [ ] Setup GORM e criar modelos base
- [ ] Implementar migrações
- [ ] Criar API REST de dispositivos (CRUD)
- [ ] Implementar autenticação JWT
- [ ] Setup logging e error handling

**Entregáveis:**
- API REST funcional
- Autenticação funcionando
- Migrations executando
- Swagger/OpenAPI documentado

**Responsável:** A definir
**Status:** 🔴 Pendente

---

### 🔲 Fase 2: SNMP Engine (Semanas 6-8)
> Comunicação SNMP com dispositivos

**Objetivos:**
- [ ] Implementar cliente SNMP com gosnmp
- [ ] Criar serviço de discovery de rede
- [ ] Implementar coleta de métricas (polling ativo)
- [ ] **Implementar receptor de SNMP Traps**
- [ ] Criar polling schedule
- [ ] Suportar SNMP v1, v2c e v3

**Entregáveis:**
- Módulo SNMP funcional
- Discovery automático
- Coleta de métricas
- **Receptor de Traps funcional**

**Responsável:** A definir
**Status:** 🔴 Pendente

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

     Set          Out          Nov
     ▼            ▼            ▼
┌─────────┐ ┌─────────┐ ┌─────────┐
│ Fase 6  │ │ Fase 7  │ │ v1.0.0! │
│Polimento│ │ Release │
│  3 sem  │ │  2 sem  │
└─────────┘ └─────────┘ └─────────┘
```

##里程碑 (Milestones)

| Versão | Descrição | Data Alvo |
|--------|-----------|-----------|
| v0.1.0 | Backend + API base | 2026-04-15 |
| v0.2.0 | SNMP funcional | 2026-05-10 |
| v0.3.0 | Frontend básico | 2026-06-15 |
| v0.4.0 | Mapas de rede | 2026-07-15 |
| v0.5.0 | Sistema de alertas | 2026-08-15 |
| v1.0.0 | Release production-ready | 2026-11-01 |

---

*Última atualização: 2026-03-20*
