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
ICON_NAME="network-transmit-receive"

echo "========================================================"
echo "🚀 Iniciando Instalação do LAGMON para o usuário: $REAL_USER"
echo "========================================================"

# [PASSO 1/5] Diretórios
echo "📂 [1/5] Criando estrutura de diretórios..."
mkdir -p "$BIN_DIR" "$AUTOSTART_DIR" "$APPS_DIR" "$ICONS_DIR" "$CONFIG_DIR" "$DATA_DIR"

# [PASSO 2/5] Grupo de Segurança
echo "🔐 [2/5] Configurando permissões de rede (ICMP/Ping)..."
getent group "$GROUP_NAME" > /dev/null || sudo groupadd -r "$GROUP_NAME"
sudo usermod -aG "$GROUP_NAME" "$REAL_USER"

GID_ALVO=$(getent group "$GROUP_NAME" | cut -d: -f3)
echo "net.ipv4.ping_group_range = $GID_ALVO $GID_ALVO" | sudo tee /etc/sysctl.d/99-lagmon.conf > /dev/null
sudo sysctl -p /etc/sysctl.d/99-lagmon.conf > /dev/null || true

# [PASSO 3/5] Instalação do Binário
echo "📦 [3/5] Instalando binário em $BIN_DIR..."
[ ! -f "bin/lagmon" ] && echo "❌ Erro: Binário bin/lagmon não encontrado." && exit 1

pkill lagmon || true
cp bin/lagmon "$FINAL_BIN"
chmod 755 "$FINAL_BIN"
sudo chown "$REAL_USER":"$REAL_USER" "$FINAL_BIN"
sudo setcap cap_net_raw=+ep "$FINAL_BIN" 2>/dev/null || true

# [PASSO 4/5] Interface e Ícones
echo "🎨 [4/5] Configurando identidade visual..."
[ -f "appicon.png" ] && cp appicon.png "$ICONS_DIR/lagmon.png" && ICON_NAME="$ICONS_DIR/lagmon.png" || true

# [PASSO 5/5] Atalhos de Sistema
echo "🖥️  [5/5] Criando entradas de menu e autostart..."
DESKTOP_FILE="/tmp/lagmon.desktop"

cat <<EOF > "$DESKTOP_FILE"
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

cp "$DESKTOP_FILE" "$APPS_DIR/lagmon.desktop"
mv "$DESKTOP_FILE" "$AUTOSTART_DIR/lagmon.desktop"
update-desktop-database "$APPS_DIR" 2>/dev/null || true

# Conclusão exata solicitada
echo "--------------------------------------------------------"
echo "✅ Instalação concluída com sucesso!"
echo "📂 Configurações: $CONFIG_DIR"
echo "📊 Base de dados: $DATA_DIR"
echo "🌐 Site: https://luizhanauer.github.io/lagmon/"
echo "--------------------------------------------------------"
echo "⚠️  IMPORTANTE: Você foi adicionado ao grupo '$GROUP_NAME'."
echo "👉 Para que as permissões de rede funcionem, faça REINICIALIZAÇÃO DO SISTEMA."
echo "--------------------------------------------------------"