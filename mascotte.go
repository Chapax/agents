package main

import "github.com/charmbracelet/lipgloss"

// L'Akita de l'installateur, d'après les photos de Chapax.
//
// Ce qui le fait reconnaître, dans l'ordre d'importance :
//   1. il est BRINGÉ — des rayures sombres sur un brun chaud, pas un beige uni
//   2. il est SUR LE DOS — ventre crème offert, quatre pattes en l'air
//   3. museau crème, truffe noire, poil long sur le poitrail
//
// Légende du dessin :
//   b = brun de base        d = rayure bringée (sombre)
//   c = crème (ventre, pattes, museau)   n = truffe   e = œil   . = vide
//
// La pose sur le dos n'est pas décorative : un chien qui s'offre le ventre ne
// monte pas la garde. C'est le bon signal pour un installateur qui travaille
// pendant qu'on le regarde.

var (
	brun   = lipgloss.NewStyle().Foreground(lipgloss.Color("137")) // brun chaud
	rayure = lipgloss.NewStyle().Foreground(lipgloss.Color("238")) // bringé sombre
	creme  = lipgloss.NewStyle().Foreground(lipgloss.Color("223")) // ventre, museau
	noir   = lipgloss.NewStyle().Foreground(lipgloss.Color("232")) // truffe
	oeil   = lipgloss.NewStyle().Foreground(lipgloss.Color("94"))  // œil ambre
)

type mascotte struct {
	nom    string
	lignes []string
}

// Sur le dos, tête à droite, quatre pattes en l'air, ventre crème au centre.
var mascotteSurLeDos = mascotte{
	nom: "sur le dos",
	lignes: []string{
		`.....cc...........cc....`,
		`..cc.cc.....cc....cc.cc.`,
		`..cc.cc.....cc....cc.cc.`,
		`.ddccddccddddccddddccdd.`,
		`ddbbbbbbbbbbbbbbbbbbbbdd`,
		`dbbddddddddddddddddddbbd`,
		`dbbccccccccccccccccccbne`,
		`dbbccccccccccccccccccbnn`,
		`dbbddddddddddddddddddbbd`,
		`.ddbbbbbbbbbbbbbbbbbbdd.`,
		`...dd..dd.......dd..dd..`,
	},
}

var mascotteFini = mascotte{nom: "fini", lignes: mascotteSurLeDos.lignes}

// dessine rend la mascotte coloriée, une ligne par entrée.
func (m mascotte) dessine() []string {
	out := make([]string, 0, len(m.lignes))
	for _, l := range m.lignes {
		var s string
		for _, ch := range l {
			switch ch {
			case 'b':
				s += brun.Render("█")
			case 'd':
				s += rayure.Render("█")
			case 'c':
				s += creme.Render("█")
			case 'n':
				s += noir.Render("█")
			case 'e':
				s += oeil.Render("█")
			default:
				s += " "
			}
		}
		out = append(out, s)
	}
	return out
}

// cote rend la mascotte à gauche d'un bloc de texte, ligne par ligne.
// Le texte est aligné en haut ; les lignes de mascotte en trop sont ignorées.
func cote(m mascotte, texte []string) string {
	art := m.dessine()
	n := len(art)
	if len(texte) > n {
		n = len(texte)
	}
	lignes := make([]string, n)
	for i := 0; i < n; i++ {
		g := "                        "
		if i < len(art) {
			g = art[i]
		}
		var d string
		if i < len(texte) {
			d = texte[i]
		}
		lignes[i] = g + "   " + d
	}
	return lipgloss.JoinVertical(lipgloss.Left, lignes...)
}
