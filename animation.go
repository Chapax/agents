package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// L'animation.
//
// Le dauphin de DeepSeek clignote et remue : c'est ce qui fait qu'on le regarde
// au lieu de lire. On reproduit la même mécanique, avec un Akita.
//
// Deux mouvements seulement, parce que deux suffisent et que trois font du bruit :
//   - il cligne des yeux, une fois toutes les huit images
//   - sa queue balance, une image sur deux
//
// Les deux sont calculés à partir du dessin de base, pas dessinés à la main :
// une image de plus ne coûte alors aucune ligne d'art supplémentaire.

const (
	tickAnim  = 420 * time.Millisecond
	framesPar = 8
)

type tickMsg struct{}

// cligne renvoie les lignes avec les yeux fermés : les pixels d'œil prennent la
// couleur du pelage. Fermer un œil, c'est le faire disparaître, pas le dessiner.
func cligne(lignes []string) []string {
	out := make([]string, len(lignes))
	for i, l := range lignes {
		b := []byte(l)
		for j, c := range b {
			if c == 'n' {
				b[j] = 'p'
			}
		}
		out[i] = string(b)
	}
	return out
}

// remue décale le bloc de queue d'une colonne, alternativement à gauche et à
// droite. On agit sur la dernière ligne, celle qui porte la queue.
func remue(lignes []string, versLaGauche bool) []string {
	out := make([]string, len(lignes))
	copy(out, lignes)
	d := len(out) - 1
	if d < 1 {
		return out
	}
	l := []byte(out[d])
	if versLaGauche && len(l) > 1 && l[0] == '.' {
		copy(l, l[1:])
		l[len(l)-1] = '.'
	} else if !versLaGauche && len(l) > 1 && l[len(l)-1] == '.' {
		copy(l[1:], l[:len(l)-1])
		l[0] = '.'
	}
	out[d] = string(l)
	return out
}

// image construit l'image n du dessin.
func (m mascotte) image(n int) mascotte {
	lignes := make([]string, len(m.lignes))
	copy(lignes, m.lignes)

	if n%framesPar == 3 {
		lignes = cligne(lignes)
	}
	lignes = remue(lignes, n%2 == 0)

	return mascotte{nom: m.nom, lignes: lignes}
}

// ---------------------------------------------------------------- modèle animé

type installModel struct {
	man     *Manifest
	etapes  []etape
	courant int
	frame   int
}

func newInstallModel(m *Manifest) installModel {
	return installModel{man: m, etapes: etapesParDefaut(), courant: 2}
}

func (m installModel) Init() tea.Cmd {
	return tea.Tick(tickAnim, func(time.Time) tea.Msg { return tickMsg{} })
}

func (m installModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "enter", " ":
			if m.courant < len(m.etapes) {
				m.etapes[m.courant].etat = "fait"
				m.courant++
			} else {
				return m, tea.Quit
			}
		}
	case tickMsg:
		m.frame++
		return m, tea.Tick(tickAnim, func(time.Time) tea.Msg { return tickMsg{} })
	}
	return m, nil
}

func (m installModel) View() string {
	var entrees []Entry
	if m.man != nil {
		entrees = m.man.Entries
	}
	return ecranInstallationAnime(entrees, m.etapes, m.courant, m.frame)
}
