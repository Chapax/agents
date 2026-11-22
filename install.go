package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Écran d'installation.
//
// L'idée : on ne demande pas « veux-tu installer ? » à un humain qui a déjà
// tapé la commande. On lui montre ce qui va se passer, entrée par entrée, avec
// son Akita à côté qui attend. Rien ne se fait sans qu'il ait vu la liste.

var (
	titreInstall = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	etapeFaite   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	etapeAPlus   = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	etapeAttente = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	encadre      = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)
)

type etape struct {
	libelle string
	etat    string // fait | aplus | attente
	detail  string
}

func marque(etat string) string {
	switch etat {
	case "fait":
		return etapeFaite.Render("✓")
	case "aplus":
		return etapeAPlus.Render("→")
	default:
		return etapeAttente.Render("·")
	}
}

// ecranInstallation compose l'écran. pur : aucun effet de bord, ce qui permet
// de le rendre à l'identique dans un test ou une capture.
func ecranInstallation(entrees []Entry, etapes []etape, courant int) string {
	return ecranInstallationAnime(entrees, etapes, courant, 0)
}

// ecranInstallationAnime est la même chose, avec l'image n de l'Akita.
func ecranInstallationAnime(entrees []Entry, etapes []etape, courant, frame int) string {
	var texte []string
	texte = append(texte, titreInstall.Render("Bonjour.")+"  "+etapeAttente.Render("je suis ton Akita d'installation."))
	texte = append(texte, "")
	texte = append(texte, etapeAttente.Render(fmt.Sprintf("%d entrées dans le manifeste · %d étapes", len(entrees), len(etapes))))
	texte = append(texte, "")

	m := mascotteAssis
	if courant >= len(etapes) {
		m = mascotteFini
	} else if courant > 0 {
		m = mascotteTravaille
	}
	m = m.image(frame)

	var droite []string
	for i, e := range etapes {
		ligne := marque(e.etat) + " " + e.libelle
		if i == courant {
			ligne = lipgloss.NewStyle().Bold(true).Render(ligne)
		}
		droite = append(droite, ligne)
		if e.detail != "" && i <= courant {
			droite = append(droite, "  "+etapeAttente.Render(e.detail))
		}
	}

	corps := strings.Join(texte, "\n") + "\n\n" + cote(m, droite)
	bas := etapeAttente.Render("entrée : continuer · échap : tout arrêter, rien n'est modifié")
	return encadre.Render(corps) + "\n\n" + bas
}

// etapesParDefaut décrit ce que l'installateur fait, dans l'ordre.
func etapesParDefaut() []etape {
	return []etape{
		{"lire le manifeste", "fait", "~/.agents/agents.yaml"},
		{"regarder ce qui est déjà là", "fait", "aucun fichier existant n'est écrasé"},
		{"proposer les entrées absentes", "aplus", "tu choisis, une par une"},
		{"écrire les fichiers", "attente", ""},
		{"vérifier que ça répond", "attente", ""},
	}
}

// rendreDemo imprime l'écran sans terminal interactif — pour les captures.
func rendreDemo() string {
	m, err := loadManifest(trouverManifeste(""))
	var entrees []Entry
	if err == nil {
		entrees = m.Entries
	}
	return ecranInstallation(entrees, etapesParDefaut(), 2) + "\n"
}

func init() {
	_ = strings.TrimSpace
}

// isTTY dit si la sortie est un vrai terminal. Sans ça, on n'ouvre pas
// d'interface plein écran : on imprime l'écran tel quel.
func isTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
