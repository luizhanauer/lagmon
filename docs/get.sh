#!/bin/sh
set -e

# Configurações do repositório LAGMON
REPO="luizhanauer/lagmon"
BINARY_NAME="lagmon"

# Identificação da arquitetura
ARCH=$(uname -m)
case $ARCH in
    x86_64) ASSET_ARCH="amd64" ;;
    aarch64) ASSET_ARCH="arm64" ;;
    *) echo "❌ Arquitetura não suportada: $ARCH"; exit 1 ;;
esac

echo ">>> 📦 Instalador via Rede: LAGMON (Network Monitor)"
echo ">>> Fonte: https://github.com/$REPO"

# Verificação de dependências básicas
if ! command -v curl >/dev/null; then echo "❌ Erro: 'curl' necessário."; exit 1; fi
if ! command -v tar >/dev/null; then echo "❌ Erro: 'tar' necessário."; exit 1; fi

# Criação de diretório temporário
TMP_DIR=$(mktemp -d)
FILENAME="${BINARY_NAME}_linux_${ASSET_ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"

echo ">>> ⬇️  Baixando release mais recente..."
if ! curl -f -L "$URL" -o "$TMP_DIR/$FILENAME"; then
    echo "❌ Erro ao baixar release. Verifique se a tag de release existe no GitHub."
    rm -rf "$TMP_DIR"
    exit 1
fi

echo ">>> 📂 Extraindo arquivos..."
tar -xzf "$TMP_DIR/$FILENAME" -C "$TMP_DIR"

echo ">>> 🚀 Iniciando script de instalação..."
# Executa o install.sh que contém a configuração de permissões ICMP
cd "$TMP_DIR"
chmod +x install.sh
sh ./install.sh

# Limpeza e finalização
cd - > /dev/null
rm -rf "$TMP_DIR"
echo ">>> ✅ Setup finalizado. O LAGMON já pode ser encontrado no seu menu de aplicativos."