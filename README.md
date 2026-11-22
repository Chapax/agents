# agents

Un seul endroit pour tes **skills**, **serveurs MCP**, **plugins** et **réglages** d'agents.

Aujourd'hui, chaque outil range ses skills dans son coin, les serveurs MCP vivent
dans trois fichiers de configuration différents, et personne ne sait dire si ce qui
est déclaré est réellement installé sur la machine. `agents` lit un manifeste
unique et te montre la vérité.

## Installation

```sh
curl -fsSL https://raw.githubusercontent.com/Chapax/agents/main/install.sh | sh
```

Pas de `sudo`, aucune écriture ailleurs que dans `~/.local/bin`.

## Usage

```sh
agents --init       # crée ~/.agents/agents.yaml à partir du modèle embarqué
agents              # interface : naviguer, activer, désactiver
agents --json       # état complet, pour les agents et les scripts
agents --check      # code 1 si une entrée activée manque, agent par agent
```

Le manifeste est cherché en remontant depuis le dossier courant, puis dans
`~/.agents/agents.yaml`. Pas besoin d'être « dans un projet » pour lancer l'outil.

Dans l'interface :

- `↑` `↓` ou `j` `k` — naviguer
- `espace` ou `entrée` — activer / désactiver
- **clic de souris** sur une ligne — activer / désactiver directement
- molette — défiler
- `s` — enregistrer
- `q` — quitter

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

### Emplacements connus

| Agent | Config | Skills | MCP |
|---|---|---|---|
| Hermes | `~/.hermes/config.yaml` | `~/.hermes/skills` | `~/.hermes/config.yaml` |
| Claude Code | `~/.claude/settings.json` | `~/.claude/skills` | `~/.claude.json` |
| Codex | `~/.codex/config.toml` | `~/.codex/skills` | `~/.codex/config.toml` |
| Gemini / Antigravity | `~/.gemini/settings.json` | — | `~/.gemini/config/mcp_config.json` |
| opencode | `~/.config/opencode/opencode.json` | — | — |
| Cursor | `~/.cursor/hooks.json` | — | — |

Hermes garde ses profils séparés (`~/.hermes/profiles/<nom>/`), donc deux agents
Hermes peuvent avoir deux configurations distinctes depuis le même manifeste.

## Ce que ce projet ne fait pas

- il **n'installe pas** les skills à ta place : il montre, il écrit le manifeste,
  et l'installation reste une décision explicite
- il ne modifie aucune configuration d'un autre outil
- il ne télécharge rien sans que tu l'aies demandé

## État

v0.1 — lecture du manifeste, affichage de l'état réel, activation, souris,
sortie JSON, contrôle `--check`.

## Licence

MIT.
