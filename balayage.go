package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Le balayage, repris de deepseek-harness-cli
// (packages/ui/tui/tests/snapshots/banner-gradient.expected.txt).
//
// Leur bannière n'anime pas un dessin : elle colorie le TITRE caractère par
// caractère sur huit arrêts de couleur, du bleu profond au bleu clair, puis
// fait glisser ce dégradé. C'est ce qui rend l'écran vivant sans rien dessiner.
//
// On reprend exactement l'idée, avec les couleurs de l'Akita : du brun sombre
// au crème, en passant par le bringé.

var arrets = [][3]int{
	{0x6b, 0x4a, 0x2f}, // brun profond
	{0x8a, 0x63, 0x3d},
	{0xa8, 0x7c, 0x4c},
	{0xc0, 0x96, 0x62},
	{0xd2, 0xae, 0x7d},
	{0xdf, 0xc2, 0x99},
	{0xe8, 0xd3, 0xb3},
	{0xd9, 0xcd, 0xba}, // crème
}

// couleur rend la n-ième teinte du dégradé, en interpolant entre deux arrêts.
func couleur(n, total int) lipgloss.Color {
	if total <= 1 {
		return lipgloss.Color("#d9cdba")
	}
	pos := float64(n) / float64(total-1) * float64(len(arrets)-1)
	i := int(math.Floor(pos))
	if i >= len(arrets)-1 {
		i = len(arrets) - 2
	}
	f := pos - float64(i)
	a, b := arrets[i], arrets[i+1]
	r := int(float64(a[0]) + (float64(b[0])-float64(a[0]))*f)
	g := int(float64(a[1]) + (float64(b[1])-float64(a[1]))*f)
	bl := int(float64(a[2]) + (float64(b[2])-float64(a[2]))*f)
	return lipgloss.Color("#" + hex2(r) + hex2(g) + hex2(bl))
}

func hex2(v int) string {
	if v < 0 {
		v = 0
	}
	if v > 255 {
		v = 255
	}
	s := strconv.FormatInt(int64(v), 16)
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

// balaie colorie le texte caractère par caractère, en décalant le dégradé de
// `phase` crans. La phase avance à chaque image : le dégradé glisse.
func balaie(texte string, phase int) string {
	r := []rune(texte)
	n := len(r)
	if n == 0 {
		return ""
	}
	var b strings.Builder
	for i, c := range r {
		// le décalage tourne en boucle sur la longueur du texte
		idx := (i + phase) % (n + len(arrets))
		if idx > n-1 {
			idx = n - 1
		}
		st := lipgloss.NewStyle().Foreground(couleur(idx, n)).Bold(true)
		b.WriteString(st.Render(string(c)))
	}
	return b.String()
}

// slogans : une phrase différente à chaque lancement, comme leur
// tui-startup-slogans. Courtes, et jamais deux fois la même de suite.
var slogans = []string{
	"un seul manifeste, et la vérité de la machine.",
	"ce qui est déclaré n'est pas toujours installé.",
	"on vérifie avant de promettre.",
	"tes agents, leur configuration, un seul endroit.",
	"pas de magie : des chemins, des fichiers, des faits.",
}

var dernierSlogan = -1

func slogan() string {
	if len(slogans) == 0 {
		return ""
	}
	if len(slogans) == 1 {
		return slogans[0]
	}
	i := rand.Intn(len(slogans))
	for i == dernierSlogan {
		i = rand.Intn(len(slogans))
	}
	dernierSlogan = i
	return slogans[i]
}

// entete compose la ligne de titre balayée, prête pour View().
func entete(phase int) string {
	return balaie("agents", phase) + "  " + lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).Render("— skills · MCP · plugins · réglages")
}

var _ = fmt.Sprintf
