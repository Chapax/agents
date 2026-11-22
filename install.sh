#!/bin/sh
# agents — installeur.
#
#   curl -fsSL https://raw.githubusercontent.com/Chapax/agents/main/install.sh | sh
#
# Ce que ça fait : télécharge le binaire de la dernière version, le pose dans
# ~/.local/bin, vérifie qu'il s'exécute, et dit quoi faire ensuite.
# Ne demande pas sudo, n'écrit nulle part ailleurs, ne touche à aucune config.

set -eu

DEPOT="Chapax/agents"
DEST="${AGENTS_BIN_DIR:-$HOME/.local/bin}"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "architecture non prise en charge : $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
  linux|darwin) ;;
  *) echo "système non pris en charge : $OS" >&2; exit 1 ;;
esac

echo "agents — $OS/$ARCH"

# dernière version publiée
TAG=$(curl -fsSL "https://api.github.com/repos/$DEPOT/releases/latest" 2>/dev/null \
      | grep -m1 '"tag_name"' | sed 's/.*: *"//; s/".*//') || TAG=""

mkdir -p "$DEST"
TMP=$(mktemp)

if [ -n "$TAG" ]; then
  URL="https://github.com/$DEPOT/releases/download/$TAG/agents_${OS}_${ARCH}"
  echo "téléchargement de $TAG"
  if curl -fsSL -o "$TMP" "$URL"; then
    chmod +x "$TMP"
    mv "$TMP" "$DEST/agents"
    echo "installé dans $DEST/agents"
    "$DEST/agents" --check >/dev/null 2>&1 || true
    echo
    echo "Pour commencer :  cd <ton projet> && agents"
    echo "S'il n'est pas dans le PATH :  export PATH=\"\$HOME/.local/bin:\$PATH\""
    exit 0
  fi
  echo "binaire introuvable pour cette version, tentative depuis les sources" >&2
fi

# repli : compiler depuis les sources
if command -v go >/dev/null 2>&1; then
  echo "compilation depuis les sources"
  TMPDIR_SRC=$(mktemp -d)
  git clone --depth 1 "https://github.com/$DEPOT.git" "$TMPDIR_SRC/agents" 2>/dev/null
  ( cd "$TMPDIR_SRC/agents" && go build -o "$DEST/agents" . )
  rm -rf "$TMPDIR_SRC"
  echo "installé dans $DEST/agents"
  exit 0
fi

echo "échec : aucune version publiée et Go absent." >&2
echo "Installe Go (https://go.dev/dl) puis relance cette commande." >&2
rm -f "$TMP"
exit 1
