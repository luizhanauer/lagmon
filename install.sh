#!/bin/sh
set -e

cd "$(dirname "$0")"

BIN_DIR="$HOME/.local/bin"
AUTOSTART_DIR="$HOME/.config/autostart"
APPS_DIR="$HOME/.local/share/applications"
ICONS_DIR="$HOME/.local/share/icons"

FINAL_BIN="$BIN_DIR/lagmon"

echo ">>> 🚀 Instalando LAGMON..."

mkdir -p "$BIN_DIR" "$AUTOSTART_DIR" "$APPS_DIR" "$ICONS_DIR"

# Copia o binário
if [ -f "bin/lagmon" ]; then
    cp bin/lagmon "$FINAL_BIN"
    chmod +x "$FINAL_BIN"
    
    # Adicionando permissões para Raw Sockets (ICMP) sem precisar de sudo constante
    echo ">>> 🔓 Solicitando permissões para monitoramento de rede (ICMP)..."
    sudo setcap cap_net_raw=+ep "$FINAL_BIN"
else
    echo "❌ Erro: Binário não encontrado."
    exit 1
fi

ICON_NAME="network-transmit-receive" # Ícone de rede padrão do sistema
if [ -f "appicon.png" ]; then
    cp appicon.png "$ICONS_DIR/lagmon.png"
    ICON_NAME="$ICONS_DIR/lagmon.png"
fi

# Criação do Arquivo .desktop atualizado para LAGMON
cat <<EOF > /tmp/lagmon.desktop
[Desktop Entry]
Type=Application
Name=LAGMON
Comment=Monitor de Latência de Rede em Tempo Real
Exec=$FINAL_BIN
Icon=$ICON_NAME
Terminal=false
Categories=Utility;Network;
Keywords=network;latency;ping;monitor;
X-GNOME-Autostart-enabled=true
EOF

cp /tmp/lagmon.desktop "$APPS_DIR/lagmon.desktop"
mv /tmp/lagmon.desktop "$AUTOSTART_DIR/lagmon.desktop"

update-desktop-database "$APPS_DIR" 2>/dev/null || true

echo "--------------------------------------------------------"
echo "✅ LAGMON instalado com sucesso!"
echo "📂 Menu: Disponível em 'Mostrar Aplicativos'"
echo "⚙️  Config: Iniciará automaticamente com o sistema"
echo "--------------------------------------------------------"