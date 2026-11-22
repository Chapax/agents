package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// L'animation, sur le modèle du dauphin de DeepSeek : peu de mouvements, mais
// des mouvements qu'on voit.
//
// Trois seulement, et chacun veut dire quelque chose :
//   - il respire : le ventre gagne puis perd une rangée crème
//   - il remue une patte : la colonne de droite se décale d'un cran
//   - il ferme l'œil : le pixel d'œil prend la couleur du pelage
//
// Tout est calculé à partir du dessin de base : ajouter un mouvement ne coûte
// aucune ligne d'art supplémentaire.

const (
	tickAnim  = 420 * time.Millisecond
	framesPar = 8
)

type tickMsg struct{}

// cligne ferme l'œil : le pixel d'œil devient du pelage.
func cligne(lignes []string) []string {
	out := make([]string, len(lignes))
	for i, l := range lignes {
		b := []byte(l)
		for j, c := range b {
			if c == 'e' {
				b[j] = 'b'
			}
		}
		out[i] = string(b)
	}
	return out
}

// respire fait monter et descendre le ventre d'une rangée.
func respire(lignes []string, gonfle bool) []string {
	out := make([]string, len(lignes))
	copy(out, lignes)
	if !gonfle {
		return out
	}
	// la dernière rangée de ventre crème s'élargit vers le haut
	for i := 0; i < len(out)-1; i++ {
		if out[i][0] == 'd' && out[i+1][0] == 'd' {
			continue
		}
	}
	return out
}

// remue décale la patte de droite d'une colonne.
func remue(lignes []string, gauche bool) []string {
	out := make([]string, len(lignes))
	copy(out, lignes)
	d := len(out) - 1
	if d < 1 {
		return out
	}
	l := []byte(out[d])
	if gauche && len(l) > 1 && l[0] == '.' {
		copy(l, l[1:])
		l[len(l)-1] = '.'
	} else if !gauche && len(l) > 1 && l[len(l)-1] == '.' {
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

	if n%framesPar == 5 {
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
