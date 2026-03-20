# Visão Geral - NetWatch

## O que é o NetWatch?

O NetWatch é uma ferramenta de monitoramento SNMP desenvolvida especificamente para **Provedores de Internet (ISPs)**. O objetivo é combinar as melhores funcionalidades do **Zabbix** e do **The Dude** da Mikrotik em uma solução unificada.

## Problema Identificado

- **Zabbix**: Poderoso, mas complexo demais para monitorar ambientes de ISP pequeno/médio
- **The Dude**: Excelente para discovery e visualização de mapas, mas limitado em funcionalidades de alerta e análise histórica

## Solução Proposta

O NetWatch busca unir:
- ✅ Discovery automático de dispositivos (como The Dude)
- ✅ Visualização em mapa de rede (como The Dude)
- ✅ Monitoramento SNMP eficiente (como Zabbix)
- ✅ Sistema de alertas flexível (como Zabbix)
- ✅ Interface moderna e intuitiva
- ✅ Aplicação desktop instalável (Windows, Linux, macOS)
- ✅ Backend leve executado em VM Linux

## Funcionalidades Principais

### 1. Discovery Automático
- Varredura de redes via SNMP
- Identificação automática de dispositivos
- Detecção de tipo de equipamento (Mikrotik, Cisco, Huawei, etc.)

### 2. Monitoramento SNMP
- Coleta de métricas via SNMP (CPU, Memória, Interface, Temperatura, etc.)
- Gráficos históricos
- Suporte a SNMP v1, v2c e v3

### 3. Mapas de Rede
- Visualização拓扑ica dos dispositivos
- Status visual (online/offline/warning)
- Atualização em tempo real
- Editor de mapas interativo

### 4. Sistema de Alertas
- Múltiplos canais (Email, Telegram, SMS, Webhook)
- Regras de alerta configuráveis
- Escalação de alertas
- Histórico de notificações

### 5. Dashboard
- Visão geral do ambiente
- Widgets customizáveis
- Indicadores de saúde da rede

## Público-Alvo

- Provedores de Internet (ISPs) pequenos e médios
- NOCs de telecomunicações
- Administradores de rede que precisam de uma ferramenta simples e eficaz

## Arquitetura Básica

```
┌─────────────────────────────────────────────────────────────┐
│                    Desktop App (Electron)                   │
│                  Windows / Linux / macOS                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ REST API / WebSocket
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Backend - Debian 13 (VM Proxmox)            │
│                      Go + Fiber + PostgreSQL                │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ SNMP
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Dispositivos de Rede                     │
│          Roteadores, Switches, APs, OLTs, etc.              │
└─────────────────────────────────────────────────────────────┘
```

---

*Última atualização: 2026-03-20*
