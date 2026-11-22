package main

import "github.com/charmbracelet/lipgloss"

// L'Akita de l'installateur.
//
// Dessiné en blocs pleins, trois teintes : le pelage fauve, le poitrail blanc,
// et l'œil. Il occupe 11 colonnes sur 9 lignes : assez grand pour qu'on
// reconnaisse la race, assez petit pour tenir à côté du texte sur 80 colonnes.
//
// Trois états, comme le dauphin : il attend, il travaille, il a fini.
// L'oreille pointue et la queue enroulée sont ce qui fait lire « akita » et pas
// « chien » — ce sont les deux traits à ne jamais simplifier.

var (
	poil   = lipgloss.NewStyle().Foreground(lipgloss.Color("180")) // fauve
	clair  = lipgloss.NewStyle().Foreground(lipgloss.Color("255")) // poitrail
	sombre = lipgloss.NewStyle().Foreground(lipgloss.Color("236")) // œil, truffe
	assise = lipgloss.NewStyle().Foreground(lipgloss.Color("94"))  // ombre au sol
)

type mascotte struct {
	nom    string
	lignes []string // chaque caractère : p = poil, b = blanc, n = noir, . = vide
}

var mascotteAssis = mascotte{
	nom: "assis",
	lignes: []string{
		`..p......p..`,
		`.ppp....ppp.`,
		`.pppppppppp.`,
		`pp.pppppp.pp`,
		`pppnnppnnppp`,
		`pppppbbppppp`,
		`.pppbbbbppp.`,
		`..ppbbbbpp..`,
		`...ppbbpp...`,
	},
}

var mascotteTravaille = mascotte{
	nom: "travaille",
	lignes: []string{
		`..p......p..`,
		`.ppp....ppp.`,
		`.pppppppppp.`,
		`pp.pppppp.pp`,
		`pppnnppnnppp`,
		`pppppbbppppp`,
		`.pppbbbbppp.`,
		`..ppbbbbpp..`,
		`..pp.bb.pp..`,
	},
}

var mascotteFini = mascotte{
	nom: "fini",
	lignes: []string{
		`..p......p..`,
		`.ppp....ppp.`,
		`.pppppppppp.`,
		`pp.pppppp.pp`,
		`pppnnppnnppp`,
		`pppppbbppppp`,
		`.pppbbbbppp.`,
		`..ppbbbbpp..`,
		`.ppppbbpppp.`,
	},
}

// dessine rend la mascotte coloriée, une ligne par entrée.
func (m mascotte) dessine() []string {
	out := make([]string, 0, len(m.lignes))
	for _, l := range m.lignes {
		var s string
		for _, c := range l {
			switch c {
			case 'p':
				s += poil.Render("█")
			case 'b':
				s += clair.Render("█")
			case 'n':
				s += sombre.Render("█")
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
		var g, d string
		if i < len(art) {
			g = art[i]
		} else {
			g = "           "
		}
		if i < len(texte) {
			d = texte[i]
		}
		lignes[i] = g + "   " + d
	}
	return lipgloss.JoinVertical(lipgloss.Left, lignes...)
}
