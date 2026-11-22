package main

import (
	"strings"
	"testing"
)

// L'animation doit réellement changer le dessin, sinon on aurait une mascotte
// figée qui prétend bouger.
func TestImagesDifferentes(t *testing.T) {
	vues := map[string]int{}
	for n := 0; n < framesPar; n++ {
		vues[strings.Join(mascotteAssis.image(n).lignes, "\n")]++
	}
	if len(vues) < 3 {
		t.Fatalf("l'animation ne produit que %d image(s) distincte(s), il en faut au moins 3", len(vues))
	}
}

// Il cligne une fois par cycle : les yeux disparaissent à l'image 3.
func TestClignement(t *testing.T) {
	ouvert := strings.Join(mascotteAssis.image(0).lignes, "")
	ferme := strings.Join(mascotteAssis.image(3).lignes, "")
	if strings.Count(ouvert, "n") == 0 {
		t.Fatal("l'image de départ n'a pas d'yeux")
	}
	if strings.Count(ferme, "n") != 0 {
		t.Fatalf("à l'image 3 les yeux devraient être fermés, il en reste %d", strings.Count(ferme, "n"))
	}
}

// La queue change de côté entre deux images.
func TestQueueRemue(t *testing.T) {
	a := mascotteAssis.image(0).lignes
	b := mascotteAssis.image(1).lignes
	if a[len(a)-1] == b[len(b)-1] {
		t.Fatal("la queue ne bouge pas entre l'image 0 et l'image 1")
	}
}

// L'écran se compose sans paniquer, même sans manifeste.
func TestEcranSansManifeste(t *testing.T) {
	out := ecranInstallationAnime(nil, etapesParDefaut(), 2, 1)
	if !strings.Contains(out, "Akita") {
		t.Fatal("l'écran ne mentionne pas la mascotte")
	}
}
