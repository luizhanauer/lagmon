#!/bin/sh
set -e

# Detecta o usuário real mesmo se o script for rodado com sudo
REAL_USER="${SUDO_USER:-$USER}"

# Configurações de Caminhos baseadas no usuário real
BIN_DIR="/home/$REAL_USER/.local/bin"
AUTOSTART_DIR="/home/$REAL_USER/.config/autostart"
APPS_DIR="/home/$REAL_USER/.local/share/applications"
ICONS_DIR="/home/$REAL_USER/.local/share/icons"
CONFIG_DIR="/home/$REAL_USER/.config/lagmon"
DATA_DIR="/home/$REAL_USER/.local/share/lagmon"

FINAL_BIN="$BIN_DIR/lagmon"
GROUP_NAME="lagmon-users"
SYSCTL_CONF="/etc/sysctl.d/99-lagmon.conf"

echo "========================================================"
echo "🧹 Iniciando Desinstalação do LAGMON para o usuário: $REAL_USER"
echo "========================================================"

# [PASSO 1/4] Encerrar Processos
echo "🛑 [1/4] Encerrando instâncias em execução..."
pkill lagmon || true

# [PASSO 2/4] Remover Arquivos e Atalhos
echo "🗑️  [2/4] Removendo binário, atalhos e ícones..."
rm -f "$FINAL_BIN"
rm -f "$AUTOSTART_DIR/lagmon.desktop"
rm -f "$APPS_DIR/lagmon.desktop"
rm -f "$ICONS_DIR/lagmon.png"
update-desktop-database "$APPS_DIR" 2>/dev/null || true

# [PASSO 3/4] Limpeza de Permissões de Rede e Grupos (Sem Else)
echo "🔐 [3/4] Revertendo configurações de segurança do kernel..."
[ -f "$SYSCTL_CONF" ] && sudo rm -f "$SYSCTL_CONF" && sudo sysctl --system > /dev/null 2>&1 || true

# Remove o usuário do grupo
id -Gn "$REAL_USER" | grep -q "$GROUP_NAME" && sudo gpasswd -d "$REAL_USER" "$GROUP_NAME" > /dev/null 2>&1 || true

# Deleta o grupo do sistema
getent group "$GROUP_NAME" > /dev/null && sudo groupdel "$GROUP_NAME" > /dev/null 2>&1 || true

# [PASSO 4/4] Limpeza de Dados
echo "📂 [4/4] Removendo pastas de configuração e banco de dados..."
rm -rf "$CONFIG_DIR"
rm -rf "$DATA_DIR"

echo "--------------------------------------------------------"
echo "✅ Desinstalação concluída com sucesso!"
echo "O LAGMON e todas as suas permissões foram removidos."
echo "--------------------------------------------------------"