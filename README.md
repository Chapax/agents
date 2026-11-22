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
agents              # interface : naviguer, activer, désactiver
agents --json       # état complet, pour les agents et les scripts
agents --check      # code 1 si une entrée activée est absente de la machine
```

Dans l'interface :

- `↑` `↓` ou `j` `k` — naviguer
- `espace` ou `entrée` — activer / désactiver
- **clic de souris** sur une ligne — activer / désactiver directement
- molette — défiler
- `s` — enregistrer
- `q` — quitter

## Le manifeste

Tout tient dans `agents.yaml` :

```yaml
version: 1
entries:
  - id: superpowers
    kind: skill          # skill | mcp | plugin | config
    name: Superpowers
    source: obra/superpowers
    detect: ~/.hermes/skills/superpowers   # le chemin qui prouve l'installation
    enabled: false
```

Le champ `detect` est le cœur du projet : il ne dit pas « je veux », il dit
**« voici comment vérifier »**. Une entrée activée dont le chemin n'existe pas
s'affiche en orange — c'est le décalage entre l'intention et la machine, celui
qui casse silencieusement une configuration.

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
