# Glossário de Termos

## Termos SNMP

| Termo | Definição |
|-------|-----------|
| **SNMP** | Simple Network Management Protocol - Protocolo para gerenciamento de dispositivos em redes IP |
| **OID** | Object Identifier - Identificador único para cada objeto gerenciável em SNMP |
| **MIB** | Management Information Base - Coleção de OIDs e suas definições |
| **Polling** | Método de monitoramento onde o servidor faz requisições periódicas aos dispositivos |
| **Trap** | Notificação enviada por um dispositivo quando ocorre um evento específico |
| **Community String** | Senha simples usada em SNMP v1 e v2c para autenticação |
| **PDU** | Protocol Data Unit - Unidade de dados do protocolo SNMP |

## Termos de Monitoramento

| Termo | Definição |
|-------|-----------|
| **Uptime** | Tempo que um dispositivo está operante sem reinicialização |
| **Latência** | Tempo de resposta de um dispositivo |
| **Packet Loss** | Percentual de pacotes perdidos na comunicação |
| **Throughput** | Quantidade de dados transferidos por segundo |
| **Bandwidth** | Largura de banda disponível |
| **Threshold** | Valor limite para disparar alertas |

## Termos do Projeto

| Termo | Definição |
|-------|-----------|
| **Discovery** | Processo automático de encontrar dispositivos na rede |
| **Node** | Dispositivo individual na rede monitorada |
| **Link** | Conexão entre dois nodes |
| **Map** | Representação gráfica da rede |
| **Widget** | Componente visual do dashboard |
| **Alert Rule** | Regra que define quando um alerta deve ser disparado |

## Termos Técnicos

| Termo | Definição |
|-------|-----------|
| **Agent** | Software que roda no dispositivo e responde queries SNMP |
| **Manager** | Sistema que faz queries aos agents (nosso backend) |
| **Goroutine** | Thread leve em Go para concorrência |
| **Goroutine Leak** | Goroutine que nunca termina, causando memory leak |
| **ORM** | Object-Relational Mapping - Mapeamento objeto-relacional |
| **WebSocket** | Protocolo de comunicação bidirecional em tempo real |
| **JWT** | JSON Web Token - Método de autenticação stateless |
| **Graceful Shutdown** | Desligamento orderly de um serviço |

## OIDs Comuns

| OID | Descrição | MIB |
|-----|-----------|-----|
| `.1.3.6.1.2.1.1.1.0` | SysDescr - Descrição do sistema | SNMPv2-MIB |
| `.1.3.6.1.2.1.1.3.0` | SysUpTime - Tempo de atividade | SNMPv2-MIB |
| `.1.3.6.1.2.1.1.5.0` | SysName - Nome do sistema | SNMPv2-MIB |
| `.1.3.6.1.2.1.2.2.1.*` | Interfaces de rede | IF-MIB |
| `.1.3.6.1.2.1.25.1.1.0` | Host resources - Uptime | HOST-RESOURCES-MIB |

## Abreviaturas

| Abreviação | Significado |
|------------|-------------|
| ISP | Internet Service Provider |
| NOC | Network Operations Center |
| VM | Virtual Machine |
| API | Application Programming Interface |
| REST | Representational State Transfer |
| CRUD | Create, Read, Update, Delete |
| DNS | Domain Name System |
| DHCP | Dynamic Host Configuration Protocol |
| LAN | Local Area Network |
| WAN | Wide Area Network |
| VLAN | Virtual Local Area Network |
| QoS | Quality of Service |
| MTU | Maximum Transmission Unit |

## Status de Dispositivo

| Status | Cor | Significado |
|--------|-----|-------------|
| `online` | 🟢 Verde | Dispositivo respondendo normalmente |
| `offline` | 🔴 Vermelho | Dispositivo não responde |
| `warning` | 🟡 Amarelo | Dispositivo com problemas menores |
| `unknown` | ⚪ Cinza | Status não determinado |
| `maintenance` | 🔵 Azul | Dispositivo em manutenção |

## Tipos de Métricas

| Tipo | Unidade | Descrição |
|------|---------|-----------|
| `gauge` | Variável | Valor atual (ex: temperatura) |
| `counter` | Incremento | Valor que só aumenta (ex: bytes) |
| `derive` | Delta | Taxa de variação |
| `absolute` | Absoluto | Valor que reset periodicamente |

---

*Última atualização: 2026-03-20*
