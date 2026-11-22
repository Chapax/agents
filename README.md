<div align="center">

<img src="docs/banner.svg" alt="agents — l'interface, lancée dans un terminal" width="100%">

<br>

**Un seul endroit pour tes skills, tes serveurs MCP, tes plugins et tes réglages d'agents.**

Aujourd'hui, chaque agent range ses skills dans son coin, les serveurs MCP vivent
dans quatre fichiers différents, et personne ne sait dire si ce qui est déclaré
existe vraiment sur la machine.

`agents` lit **un manifeste unique** et te montre **la vérité**.

<br>

[![Release](https://img.shields.io/github/v/release/Chapax/agents?style=flat-square&color=ff79c6)](https://github.com/Chapax/agents/releases)
[![Licence](https://img.shields.io/badge/licence-MIT-50fa7b?style=flat-square)](LICENSE)
[![Plateformes](https://img.shields.io/badge/linux%20·%20macOS-x86__64%20·%20arm64-8be9fd?style=flat-square)](#installation)
[![Go](https://img.shields.io/badge/go-1.24-00add8?style=flat-square)](https://go.dev)

</div>

---

## Le problème

Un skill installé pour Hermes et absent pour Codex ne se voit nulle part.
Un serveur MCP déclaré dans quatre fichiers dont deux sont orphelins ne se voit
pas non plus. On découvre le décalage au moment où ça casse — c'est-à-dire au
pire moment.

```
MANQUANT  skill  Superpowers  agent=codex  ~/.codex/skills/superpowers
```

Cette ligne, c'est le cœur du projet : **l'écart entre ce que tu veux et ce que
la machine a réellement**, nommé, agent par agent.

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/Chapax/agents/main/install.sh | sh
```

Un binaire de 3,5 Mo dans `~/.local/bin`. Pas de `sudo`, pas de dépendance,
aucune écriture ailleurs.

<details>
<summary>Autres méthodes</summary>

**Depuis les sources** — Go 1.24 ou plus récent :

```sh
git clone https://github.com/Chapax/agents.git
cd agents && go build -o ~/.local/bin/agents .
```

**Binaires directs** — [dernière version](https://github.com/Chapax/agents/releases/latest),
pour `linux/amd64` et `darwin/arm64`.

</details>

## Démarrage

```sh
agents --init     # crée ~/.agents/agents.yaml à partir du modèle embarqué
agents            # l'interface
```

| Touche | Action |
|---|---|
| `↑` `↓` · `j` `k` | naviguer |
| `espace` · `entrée` | activer / désactiver |
| **clic de souris** | activer / désactiver la ligne visée |
| molette | défiler |
| `s` | enregistrer |
| `q` | quitter |

Le manifeste est cherché en remontant depuis le dossier courant, puis dans
`~/.agents/agents.yaml`. Pas besoin d'être « dans un projet ».

## Pour les agents, pas seulement pour toi

```sh
agents --json     # l'état complet, structuré
agents --check    # code de sortie 1 si une entrée activée manque
```

Un agent peut lire son propre état avant de commencer à travailler, et refuser
de partir sur une configuration incomplète. C'est fait pour être branché dans un
script de démarrage, un hook, ou une CI.

## Le manifeste, multi-agents

Une entrée décrit **une chose**, et dit **où elle doit vivre pour chaque agent** :

```yaml
- id: superpowers
  kind: skill              # skill | mcp | plugin | config
  name: Superpowers
  source: obra/superpowers
  enabled: true
  agents: [hermes, claude, codex]
  targets:
    hermes: ~/.hermes/skills/superpowers
    claude: ~/.claude/skills/superpowers
    codex:  ~/.codex/skills/superpowers
```

`--check` ne répond pas « installé / pas installé », mais **pour quel agent** :

```
MANQUANT  skill  Superpowers  agent=codex  ~/.codex/skills/superpowers
```

C'est le décalage entre l'intention et la machine, celui qui casse
silencieusement une configuration.

### Agents et emplacements reconnus

| Agent | Configuration | Skills | MCP |
|---|---|---|---|
| **Hermes** | `~/.hermes/config.yaml` | `~/.hermes/skills` | `~/.hermes/config.yaml` |
| **Claude Code** | `~/.claude/settings.json` | `~/.claude/skills` | `~/.claude.json` |
| **Codex** | `~/.codex/config.toml` | `~/.codex/skills` | `~/.codex/config.toml` |
| **Gemini / Antigravity** | `~/.gemini/settings.json` | — | `~/.gemini/config/mcp_config.json` |
| **opencode** | `~/.config/opencode/opencode.json` | — | — |
| **Cursor** | `~/.cursor/hooks.json` | — | — |

Hermes garde ses profils séparés (`~/.hermes/profiles/<nom>/`) : deux agents
Hermes peuvent donc avoir deux configurations distinctes depuis le même
manifeste.

## Ce que ce projet ne fait pas

- il **n'installe pas** les skills à ta place : il montre, il écrit le manifeste,
  et l'installation reste une décision explicite
- il ne modifie la configuration d'aucun autre outil
- il ne télécharge rien sans que tu l'aies demandé

Un installeur qui télécharge et exécute du code doit demander confirmation, pas
agir en silence. C'est pour ça que la v0.2 montre au lieu d'installer.

## Feuille de route

- [x] **v0.1** — manifeste, état réel, activation, souris, sortie JSON
- [x] **v0.2** — une entrée par agent, `--init`, `--check` par agent
- [ ] **v0.3** — installation guidée, avec confirmation par entrée
- [ ] **v0.4** — profils par projet (`agents --profile planquex`)
- [ ] **v0.5** — vérification que chaque serveur MCP répond vraiment
- [ ] **Traduction** — anglais, puis interface bilingue

## Contribuer

Les retours les plus utiles sont les **agents que tu utilises et qui ne sont pas
dans le tableau**. Ouvre une issue avec le chemin réel de leur configuration :
c'est la donnée qui manque le plus.

## Licence

[MIT](LICENSE) — fais-en ce que tu veux.

<div align="center">
<br>
<sub>Construit avec <a href="https://github.com/charmbracelet/bubbletea">Bubble Tea</a>. Un seul binaire, aucune dépendance.</sub>
</div>
