package widget

import (
    "github.com/jroimartin/gocui"
)

// Déclaration d'un type Entrée qui contiendra les données nécessaires à créer
// une fenêtre qui prendra l'entrée de l'utilisateur.
// Code prit de: https://gist.github.com/jroimartin/3b2e943a3811d795e0718b4a95b89bec
type Entree struct {
    nom string
    x, y int
    w int
    longueur_max int
}

// Fonction qui créé une variable de type Entrée.
// Code prit de: https://gist.github.com/jroimartin/3b2e943a3811d795e0718b4a95b89bec
func CreerEntree(nom string, x, y, w, longueur_max int) *Entree{
    return &Entree{nom, x, y, w, longueur_max}
}

// Fonction qui initialise l'Entrée.
// Code prit de: https://gist.github.com/jroimartin/3b2e943a3811d795e0718b4a95b89bec
func (e *Entree) Disposition(g *gocui.Gui) error {
    v, err := g.SetView(e.nom, e.x, e.y, e.x+e.w, e.y+2)
    if err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Frame = false
        v.Editor = e
        v.Editable = true
    }
    return nil
}

// Fonction permettant de définir comment l'Entrée est modifié lorsqu'on appuit
// sur une touche.
// Code prit de: https://gist.github.com/jroimartin/3b2e943a3811d795e0718b4a95b89bec
func (e *Entree) Edit(v *gocui.View, cle gocui.Key, ch rune, mod gocui.Modifier) {
    cx, _ := v.Cursor()
    ox, _ := v.Origin()
    limite := ox+cx+1 > e.longueur_max
    switch {
        case ch != 0 && mod == 0 && !limite:
            v.EditWrite(ch)
        case cle == gocui.KeySpace && !limite:
            v.EditWrite(' ')
        case cle == gocui.KeyBackspace || cle == gocui.KeyBackspace2:
            v.EditDelete(true)
    }
}
