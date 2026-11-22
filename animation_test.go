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
		vues[strings.Join(mascotteSurLeDos.image(n).lignes, "\n")]++
	}
	if len(vues) < 2 {
		t.Fatalf("l'animation ne produit que %d image(s) distincte(s)", len(vues))
	}
}

// Il cligne une fois par cycle : le pixel d'œil disparaît à l'image 5.
func TestClignement(t *testing.T) {
	ouvert := strings.Join(mascotteSurLeDos.image(0).lignes, "")
	ferme := strings.Join(mascotteSurLeDos.image(5).lignes, "")
	if strings.Count(ouvert, "e") == 0 {
		t.Fatal("l'image de départ n'a pas d'œil")
	}
	if strings.Count(ferme, "e") != 0 {
		t.Fatalf("à l'image 5 l'œil devrait être fermé, il en reste %d", strings.Count(ferme, "e"))
	}
}

// La patte bouge entre deux images.
func TestPatteRemue(t *testing.T) {
	a := mascotteSurLeDos.image(0).lignes
	b := mascotteSurLeDos.image(1).lignes
	if a[len(a)-1] == b[len(b)-1] {
		t.Fatal("la patte ne bouge pas entre l'image 0 et l'image 1")
	}
}

// La palette est bringée : au moins trois teintes distinctes dans le dessin.
func TestPaletteBringee(t *testing.T) {
	art := strings.Join(mascotteSurLeDos.lignes, "")
	for _, c := range []string{"b", "d", "c"} {
		if !strings.Contains(art, c) {
			t.Fatalf("teinte %q absente du dessin", c)
		}
	}
}

// L'écran se compose sans paniquer, même sans manifeste.
func TestEcranSansManifeste(t *testing.T) {
	out := ecranInstallationAnime(nil, etapesParDefaut(), 2, 1)
	if !strings.Contains(out, "Akita") {
		t.Fatal("l'écran ne mentionne pas la mascotte")
	}
}
