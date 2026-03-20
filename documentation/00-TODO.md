# Tarefas Pendentes

## Tarefas Imediatas

### Setup Inicial

- [ ] Criar repositório Git
- [ ] Configurar Git hooks (pre-commit, commit-msg)
- [ ] Criar `.gitignore` completo
- [ ] Setup workspace no VSCode (settings, extensions recomendadas)
- [ ] Criar arquivo `CONTRIBUTING.md`

### Documentação

- [ ] Definir convenções de commit (Conventional Commits)
- [ ] Criar template de PR
- [ ] Criar template de issue
- [ ] Documentar TODOs do código com tags específicas

### Infraestrutura

- [ ] Criar `docker-compose.yml` para PostgreSQL
- [ ] Criar `docker-compose.yml` para backend completo
- [ ] Criar Dockerfile para backend
- [ ] Setup ambiente de variáveis (`.env.example`)

---

## Tarefas de Design

### Arquitetura

- [ ] Definir estrutura de OIDs SNMP mais comuns
- [ ] Criar enum de status de dispositivos
- [ ] Definir schema do banco de dados completo
- [ ] Documentar API com OpenAPI/Swagger

### UI/UX

- [ ] Criar design system base
- [ ] Definir paleta de cores (dark/light)
- [ ] Criar componentes base (Button, Input, Card, etc.)
- [ ] Definir tipografia
- [ ] Criar biblioteca de ícones

---

## Tarefas Técnicas

### Backend

- [ ] Implementar logging estruturado
- [ ] Configurar health check endpoint
- [ ] Implementar graceful shutdown
- [ ] Setup metrics/monitoring (Prometheus?)
- [ ] Implementar rate limiting
- [ ] Adicionar caching (Redis?)

### Frontend

- [ ] Implementar cache local (IndexedDB?)
- [ ] Configurar PWA (offline support)
- [ ] Implementar auto-update
- [ ] Adicionar crash reporting (Sentry?)

---

## Decisões Tomadas

- [x] **Multi-usuário**: Sistema com múltiplos usuários e roles
- [x] **Polling + Traps**: Receber tanto polling ativo quanto SNMP Traps
- [x] **Retenção de métricas de tráfego**: 1 ano
- [x] **Retenção de alertas**: 6 meses com sobreposição (rolling window)
- [x] **Métricas derivadas de alertas**: Ex: calcular tempo de duração de baterias baseado em eventos de "ping down"

## Decisões a Tomar

- [ ] Qual formato de configuração? (YAML, TOML, JSON, ENV)
- [ ] Fazer backup automático? Em que frequência?
- [ ] Suportar plugins/extensões no futuro?
- [ ] API pública para integrações?

---

## Perguntas em Aberto

1. ~~**Autenticação**: Será multi-usuário ou single-user?~~ → **Multi-usuário**
2. ~~**SNMP Traps**: Precisará receber traps ou apenas polling?~~ → **Ambos (polling + traps)**
3. **Mapas**: Será limitado por número de dispositivos?
4. ~~**Armazenamento de métricas**: Por quanto tempo reter dados?~~ → **1 ano para tráfego, 6 meses para alertas**
5. **Updates**: Canal stable vs beta vs nightly?
6. **Roles**: Quais tipos de usuários? (admin, operator, viewer)

---

## Prioridades

### Alta Prioridade
- Repositório Git configurado
- docker-compose para PostgreSQL
- Backend básico funcional
- SNMP discovery básico

### Média Prioridade
- Frontend desktop
- Mapas de rede
- Sistema de alertas

### Baixa Prioridade
- PWA support
- Multi-tenant
- API pública

---

*Última atualização: 2026-03-20*
