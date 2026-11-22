// agents — un seul endroit pour les skills, serveurs MCP, plugins et réglages.
//
// Lit un manifeste (agents.yaml), affiche l'état réel de la machine, et permet
// d'activer ou de désactiver chaque entrée au clavier ou à la souris.
//
//	go run .            interface
//	go run . --json     état en JSON, pour les agents
//	go run . --check    code de sortie 1 si une entrée activée n'est pas installée
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gopkg.in/yaml.v3"
)

// ---------------------------------------------------------------- manifeste

type Entry struct {
	ID      string `yaml:"id" json:"id"`
	Kind    string `yaml:"kind" json:"kind"` // skill | mcp | plugin | config
	Name    string `yaml:"name" json:"name"`
	Source  string `yaml:"source,omitempty" json:"source,omitempty"`
	Detect  string `yaml:"detect,omitempty" json:"detect,omitempty"` // chemin, ~ accepté
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Note    string `yaml:"note,omitempty" json:"note,omitempty"`
}

type Manifest struct {
	Version int     `yaml:"version" json:"version"`
	Entries []Entry `yaml:"entries" json:"entries"`
	path    string
}

func expand(p string) string {
	if strings.HasPrefix(p, "~/") {
		if h, err := os.UserHomeDir(); err == nil {
			return filepath.Join(h, p[2:])
		}
	}
	return p
}

func (e Entry) Installed() bool {
	if e.Detect == "" {
		return false
	}
	_, err := os.Stat(expand(e.Detect))
	return err == nil
}

func loadManifest(path string) (*Manifest, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Manifest
	if err := yaml.Unmarshal(raw, &m); err != nil {
		return nil, fmt.Errorf("%s : %w", path, err)
	}
	m.path = path
	return &m, nil
}

func (m *Manifest) save() error {
	out, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	tmp := m.path + ".tmp"
	if err := os.WriteFile(tmp, out, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, m.path)
}

// ---------------------------------------------------------------- modèle TUI

type model struct {
	man      *Manifest
	cursor   int
	width    int
	height   int
	status   string
	quitting bool
}

func newModel(m *Manifest) model { return model{man: m, status: "espace : activer/désactiver · j/k ou ↑↓ : naviguer · s : enregistrer · q : quitter"} }

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

	case tea.MouseMsg:
		// clic : la liste commence à la ligne 4 (titre + en-têtes)
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			row := msg.Y - 4
			if row >= 0 && row < len(m.man.Entries) {
				m.cursor = row
				m.man.Entries[row].Enabled = !m.man.Entries[row].Enabled
				m.status = "modifié — « s » pour enregistrer"
			}
		}
		// molette
		if msg.Button == tea.MouseButtonWheelUp && m.cursor > 0 {
			m.cursor--
		}
		if msg.Button == tea.MouseButtonWheelDown && m.cursor < len(m.man.Entries)-1 {
			m.cursor++
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.man.Entries)-1 {
				m.cursor++
			}
		case " ", "enter":
			if len(m.man.Entries) > 0 {
				m.man.Entries[m.cursor].Enabled = !m.man.Entries[m.cursor].Enabled
				m.status = "modifié — « s » pour enregistrer"
			}
		case "s":
			if err := m.man.save(); err != nil {
				m.status = "échec de l'enregistrement : " + err.Error()
			} else {
				m.status = "enregistré dans " + m.man.path
			}
		}
	}
	return m, nil
}

var (
	titre   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	actif   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	inactif = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	sel     = lipgloss.NewStyle().Background(lipgloss.Color("237")).Bold(true)
	alerte  = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	pied    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m model) View() string {
	if m.quitting {
		return ""
	}
	var b strings.Builder
	b.WriteString(titre.Render("agents") + inactif.Render("  —  skills, MCP, plugins, réglages") + "\n")
	b.WriteString(inactif.Render(fmt.Sprintf("%d entrées · %s", len(m.man.Entries), m.man.path)) + "\n\n")
	b.WriteString(inactif.Render(fmt.Sprintf("  %-10s %-28s %-10s %s", "TYPE", "NOM", "ÉTAT", "INSTALLÉ")) + "\n")

	for i, e := range m.man.Entries {
		inst := e.Installed()
		etat := "désactivé"
		if e.Enabled {
			etat = "activé"
		}
		instTxt := "non"
		if inst {
			instTxt = "oui"
		}
		// activé mais absent de la machine : c'est le cas qui mérite un signal
		ligne := fmt.Sprintf("  %-10s %-28s %-10s %s", e.Kind, tronque(e.Name, 28), etat, instTxt)
		switch {
		case i == m.cursor:
			b.WriteString(sel.Render(ligne))
		case e.Enabled && !inst:
			b.WriteString(alerte.Render(ligne))
		case e.Enabled:
			b.WriteString(actif.Render(ligne))
		default:
			b.WriteString(inactif.Render(ligne))
		}
		b.WriteString("\n")
	}

	if m.cursor < len(m.man.Entries) {
		e := m.man.Entries[m.cursor]
		b.WriteString("\n")
		if e.Source != "" {
			b.WriteString(inactif.Render("  source : ") + e.Source + "\n")
		}
		if e.Note != "" {
			b.WriteString(inactif.Render("  note   : ") + e.Note + "\n")
		}
		if e.Enabled && !e.Installed() {
			b.WriteString(alerte.Render("  activé mais absent de la machine") + "\n")
		}
	}

	b.WriteString("\n" + pied.Render(m.status))
	return b.String()
}

func tronque(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n < 2 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

// ---------------------------------------------------------------- main

func main() {
	manifestPath := "agents.yaml"
	if v := os.Getenv("AGENTS_MANIFEST"); v != "" {
		manifestPath = v
	}
	for i, a := range os.Args[1:] {
		if a == "--manifest" && i+2 <= len(os.Args)-1 {
			manifestPath = os.Args[i+2]
		}
	}

	m, err := loadManifest(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "agents : %v\n", err)
		os.Exit(2)
	}

	sort.SliceStable(m.Entries, func(i, j int) bool { return m.Entries[i].Kind < m.Entries[j].Kind })

	for _, a := range os.Args[1:] {
		switch a {
		case "--json":
			out, _ := json.MarshalIndent(map[string]any{
				"manifest": m.path,
				"entries":  m.Entries,
			}, "", "  ")
			fmt.Println(string(out))
			return
		case "--check":
			manquant := 0
			for _, e := range m.Entries {
				if e.Enabled && !e.Installed() {
					fmt.Printf("MANQUANT  %-10s %s (%s)\n", e.Kind, e.Name, e.Detect)
					manquant++
				}
			}
			if manquant > 0 {
				os.Exit(1)
			}
			fmt.Println("tout ce qui est activé est présent")
			return
		}
	}

	p := tea.NewProgram(newModel(m), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "agents : %v\n", err)
		os.Exit(1)
	}
}
