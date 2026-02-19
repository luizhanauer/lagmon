#!/bin/sh
set -e

cd "$(dirname "$0")"

# Caminhos Padrão de Sistema (XDG)
BIN_DIR="$HOME/.local/bin"
AUTOSTART_DIR="$HOME/.config/autostart"
APPS_DIR="$HOME/.local/share/applications"
ICONS_DIR="$HOME/.local/share/icons"

# Pastas de Dados e Configuração do LAGMON
CONFIG_DIR="$HOME/.config/lagmon"
DATA_DIR="$HOME/.local/share/lagmon"

FINAL_BIN="$BIN_DIR/lagmon"

echo ">>> 🚀 Instalando LAGMON..."

# 1. Cria diretórios necessários
mkdir -p "$BIN_DIR" "$AUTOSTART_DIR" "$APPS_DIR" "$ICONS_DIR" "$CONFIG_DIR" "$DATA_DIR"

# 2. Instala o binário e aplica permissões ICMP
if [ -f "bin/lagmon" ]; then
    cp bin/lagmon "$FINAL_BIN"
    chmod +x "$FINAL_BIN"
    
    echo ">>> 🔓 Solicitando permissão para monitoramento ICMP (setcap)..."
    sudo setcap cap_net_raw=+ep "$FINAL_BIN"
else
    echo "❌ Erro: Binário não encontrado em bin/lagmon."
    exit 1
fi

# 3. Gerenciamento do ícone
ICON_NAME="network-transmit-receive"
if [ -f "appicon.png" ]; then
    cp appicon.png "$ICONS_DIR/lagmon.png"
    ICON_NAME="$ICONS_DIR/lagmon.png"
fi

# 4. Criação do atalho (.desktop)
cat <<EOF > /tmp/lagmon.desktop
[Desktop Entry]
Type=Application
Name=LAGMON
Comment=Monitor de Latência em Tempo Real
Exec=$FINAL_BIN
Icon=$ICON_NAME
Terminal=false
Categories=Utility;Network;
Keywords=network;latency;ping;
X-GNOME-Autostart-enabled=true
EOF

cp /tmp/lagmon.desktop "$APPS_DIR/lagmon.desktop"
mv /tmp/lagmon.desktop "$AUTOSTART_DIR/lagmon.desktop"

# Isso garante que o ícone apareça no menu "Mostrar Aplicativos" na hora
update-desktop-database "$APPS_DIR" 2>/dev/null || true

echo "--------------------------------------------------------"
echo "✅ Instalação concluída com sucesso!"
echo "📂 Configurações: $CONFIG_DIR"
echo "📊 Base de dados: $DATA_DIR"
echo "🌐 Site: https://luizhanauer.github.io/lagmon/"
echo "--------------------------------------------------------"