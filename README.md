# 🌐 LAGMON

**Ferramenta de diagnóstico de rede em tempo real. Monitoramento ICMP de alta precisão, logs em SQLite e interface dashboard em estética neon.**

---

### 🚀 Instalação Rápida (Linux)

Para instalar o **LAGMON** automaticamente no seu sistema (Ubuntu 24.04+), execute o comando abaixo no terminal:

```bash
   curl -sSL https://luizhanauer.github.io/lagmon/get.sh | sh
```

> **Nota:** O instalador solicitará permissão de `sudo` apenas para configurar o `cap_net_raw`, permitindo que o app realize pings (ICMP) sem precisar ser executado como root.

---

### ✨ Funcionalidades

* **Monitoramento em Tempo Real**: Captura de latência e packet loss com precisão de microssegundos.
* **Visualização por Cards**: Aba de diagramas otimizada com cards uniformes para monitorar múltiplos nós simultaneamente.
* **Histórico Persistente**: Armazenamento automático de dados em SQLite para consultas e relatórios.
* **Relatórios Dual-Mode**: Geração de arquivos CSV (dados técnicos) e TXT (resumo amigável) diretamente na pasta Downloads.
* **Dashboard Neon**: Interface moderna construída com Vue.js 3 e uPlot para máxima performance.

### 🛠️ Stack Técnica

* **Backend**: Go 1.21 + Wails v2
* **Frontend**: Vue.js 3, TypeScript, Tailwind CSS
* **Database**: SQLite3

### 📂 Estrutura de Configuração

O projeto utiliza um arquivo `settings.json` na raiz para persistência de preferências do usuário:

* **Alvos de Monitoramento**: IPs e nomes customizados.
* **Retenção de Dados**: Período automático de limpeza de logs.
* **Configurações de UI**: Visibilidade de gráficos e diagramas.

---

## ☕ Apoie o Projeto

Se o LAGMON ajudou você, considere apoiar a manutenção do projeto:

Se você gostou do meu trabalho e quer me agradecer, você pode me pagar um café :)

<a href="https://www.paypal.com/donate/?hosted_button_id=SFR785YEYHC4E" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" style="height: 40px !important;width: 150px !important;" ></a>


---

## 📄 Licença

Este projeto está sob a licença [MIT]. Os dados de recursos de numeração de internet são providos pelo NRO e seguem suas respectivas políticas de uso.