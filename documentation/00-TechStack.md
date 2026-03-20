# Stack Tecnológica

## Decisões de Tecnologia

### Backend

| Tecnologia | Escolha | Alternativas Consideradas |
|------------|---------|---------------------------|
| **Linguagem** | Go | Rust, Python |
| **Framework** | Fiber | Gin, Echo, stdlib |
| **Banco de Dados** | PostgreSQL | MySQL, SQLite, TimescaleDB |
| **ORM** | GORM | Raw SQL, sqlx, ent |
| **Migrações** | golang-migrate | Goose, Flyway |
| **Validação** | go-playground/validator | ozzo-validation |
| **Config** | Viper | godotenv, standard JSON |

### Frontend Desktop

| Tecnologia | Escolha | Alternativas Consideradas |
|------------|---------|---------------------------|
| **Framework** | Electron | Tauri |
| **UI Library** | React | Vue, Svelte |
| **State Management** | Zustand | Redux, Jotai |
| **Routing** | React Router | wouter |
| **HTTP Client** | Axios | ky, fetch |
| **Styling** | Tailwind CSS | Chakra UI, Mantine |
| **Gráficos** | Recharts | Chart.js, ApexCharts |
| **Mapas** | React Flow | D3.js, vis.js |
| **Build** | electron-builder | electron-forge |

### Ferramentas de Desenvolvimento

| Categoria | Ferramenta |
|----------|------------|
| **Containerização** | Docker, Docker Compose |
| **CI/CD** | GitHub Actions |
| **Linting (Go)** | golangci-lint |
| **Linting (TS)** | ESLint |
| **Formatting (Go)** | gofmt, goimports |
| **Formatting (TS)** | Prettier |
| **Testes (Go)** | testing package, testify |
| **Testes (React)** | Vitest, React Testing Library |

### Bibliotecas SNMP (Go)

| Biblioteca | Uso |
|-----------|-----|
| `gosnmp` | Cliente SNMP principal |
| `goiana/snmp` | Helpers e utilitários |

## Justificativas

### Por que Go + Fiber?

1. **Performance**: Go é extremamente rápido na execução
2. **Concorrência**: Goroutines são perfeitas para coletar SNMP de múltiplos dispositivos
3. **Binário único**: Deploy simplificado no servidor
4. **Ecosistema SNMP maduro**: Bibliotecas estáveis disponíveis
5. **Fiber**: Framework Express-like, fácil migração se necessário

### Por que Electron + React?

1. **Maturidade**: Electron é a opção mais testada para desktop cross-platform
2. **Ecossistema React**: Grande quantidade de bibliotecas e componentes prontos
3. **TypeScript**: Melhor experiência de desenvolvimento
4. **WebView**: Pega o melhor do desenvolvimento web

### Por que PostgreSQL?

1. **Confiabilidade**: Robustz e maduro
2. **JSON**: Suporte nativo a JSON para dados flexíveis
3. **TimescaleDB**: Possibilidade de extensão futura para timeseries
4. **Geospatial**: PostGIS se precisar de geo-localização

## Dependências Principais (Backend)

```go
// go.mod - Dependências principais
github.com/gofiber/fiber/v2          // Web framework
github.com/gofiber/websocket/v2       // WebSocket
gorm.io/gorm                          // ORM
gorm.io/driver/postgres               // Driver PostgreSQL
github.com/gosnmp/gosnmp              // Cliente SNMP
github.com/spf13/viper                // Configuração
github.com/golang-jwt/jwt/v5          // JWT
github.com/golang-migrate/migrate/v4  // Migrações
golang.org/x/crypto                   // bcrypt
```

## Dependências Principais (Frontend)

```json
{
  "dependencies": {
    "react": "^18.x",
    "react-dom": "^18.x",
    "react-router-dom": "^6.x",
    "zustand": "^4.x",
    "axios": "^1.x",
    "recharts": "^2.x",
    "@xyflow/react": "^12.x",
    "@tanstack/react-query": "^5.x",
    "date-fns": "^3.x",
    "clsx": "^2.x"
  }
}
```

## Infraestrutura

| Componente | Especificação |
|-----------|---------------|
| **VM Host** | Proxmox |
| **OS Backend** | Debian 13 |
| **RAM Backend** | 4GB (inicial) |
| **CPU Backend** | 2 vCPU |
| **Storage** | 50GB SSD |
| **Docker** | Para desenvolvimento e futura containerização |

---

*Última atualização: 2026-03-20*
