package main

import (
    "log"
    "fmt"
    "strings"
    "strconv"
    
    
    "github.com/jroimartin/gocui"
    
    "github.com/Vacheprime/gochat/network"
    "github.com/Vacheprime/gochat/widget"
    "github.com/Vacheprime/gochat/connclient"
)


// Définition du client de façon
// globale
var client connclient.Client = connclient.Client{"", "", nil, 0}

// Variable qui indique s'il y a eu une interface auparavant (si l'utilisateur a été déconnecté)
var precedent bool = false

// Structure qui permet de créer la zone de text
type EntrerMessage struct {
    LongueurMax int
}

// Fonction à executer lorsque l'utilisateur veut envoyer un message 
func envoyerMsgServeur(g *gocui.Gui, v *gocui.View) error {
    // Enlever le \n de la fin du texte
    contenu := strings.TrimRight(v.Buffer(), "\n")
    if len(contenu) != 0 {
        // Envoyer le message en le formattant de la façon suivante:
        // msg:NomUtilisateur|CouleurID|Contenu
        err := network.EnvoyerMessage("msg:" + client.NomUtilisateur + "|" + strconv.Itoa(int(client.CouleurID)) + "|" + contenu, client.Connection)
        if err != nil {
            return err
        } else {
            // Effacer le message entré et replacer le curseur au début
            v.Clear()
            v.SetCursor(0, 0)
            v.SetOrigin(0, 0)
        }
        
        return nil
    }
    return nil
}

// Fonction permettant de définir comment l'Entrée est modifié lorsqu'on appuit
// sur une touche.
// Code prit de: https://gist.github.com/jroimartin/3b2e943a3811d795e0718b4a95b89bec
func (e *EntrerMessage) Edit(v *gocui.View, cle gocui.Key, ch rune, mod gocui.Modifier) {
    switch {
        // Pour chaque touche appuyée, effectuer l'action requise
        case ch != 0 && mod == 0 && len(v.Buffer()) < e.LongueurMax:
            v.EditWrite(ch)
        case cle == gocui.KeySpace && len(v.Buffer()) < e.LongueurMax:
            v.EditWrite(' ')
        case cle == gocui.KeyBackspace || cle == gocui.KeyBackspace2:
            v.EditDelete(true)
        case cle == gocui.KeyArrowDown:
            v.MoveCursor(0, 1, false)
        case cle == gocui.KeyArrowUp:
            v.MoveCursor(0, -1, false)
        case cle == gocui.KeyArrowLeft:
            v.MoveCursor(-1, 0, false)
        case cle == gocui.KeyArrowRight:
            v.MoveCursor(1, 0, false)
    }
}

// Fonction qui créé les fenêtres Utilisateurs, Entrée et Affichage
func initialiserFenetres(g *gocui.Gui) error {
    // Définir les variables de taille d'écran
    maxX, maxY := g.Size()
    maxXf := float32(maxX)
    maxYf := float32(maxY)

    // Créer les fenêtres
    if VUtilisateurs, err := g.SetView("Utilisateurs", 0, 0, int(maxXf*0.3), int(maxYf*0.7-1)); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        VUtilisateurs.Clear()
        VUtilisateurs.Frame = true
        VUtilisateurs.Title = "Utilisateurs"
        VUtilisateurs.Wrap = true
    }
    
    if VAffichage, err := g.SetView("Affichage", int(maxXf*0.3+1), 0, maxX-1, int(maxYf*0.7-1)); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        
        VAffichage.Frame = true
        VAffichage.Title = "Messages"
        VAffichage.Wrap = true
    }
    if VEntree, err := g.SetView("Entree", 0, int(maxYf*0.7), maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        VEntree.Frame = true
        VEntree.Title = "Entree"
        VEntree.Editable = true
        VEntree.Wrap = true
        VEntree.Editor = &EntrerMessage{3000}
    }
    
    
    
    return nil
}

func actions(g *gocui.Gui) error {
    // Actions générales
    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quitter); err != nil {
        return err
    }
    
    // Actions Flèche vers le haut à l'accueil
    if err := g.SetKeybinding("EntreeAdresseIP", gocui.KeyArrowUp, gocui.ModNone, flecheHautAccueil); err != nil {
        return err
    }
    if err := g.SetKeybinding("AccueilDemande", gocui.KeyArrowUp, gocui.ModNone, flecheHautAccueil); err != nil {
        return err
    }
    
    // Actions Flèche vers le bas à l'accueil
    if err := g.SetKeybinding("EntreeNomUtili", gocui.KeyArrowDown, gocui.ModNone, flecheBasAccueil); err != nil {
        return err
    }
    if err := g.SetKeybinding("EntreeAdresseIP", gocui.KeyArrowDown, gocui.ModNone, flecheBasAccueil); err != nil {
        return err
    }
    
    // Actions Entrer à l'accueil
    if err := g.SetKeybinding("AccueilDemande", gocui.KeyEnter, gocui.ModNone, entrerAccueil); err != nil {
        return err
    }
    
    // Actions TAB 
    if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, changerFenetre); err != nil {
        return err
    }
    
    // Actions Entrer à l'Entrée
    if err := g.SetKeybinding("Entree", gocui.KeyEnter, gocui.ModNone, envoyerMsgServeur); err != nil {
        return err
    }
    
    // Actions flèche vers le haut/bas à l'Affichage
    if err := g.SetKeybinding("Affichage", gocui.KeyArrowUp, gocui.ModNone, flecheHautAffichage); err != nil {
        return err
    }
    
    if err := g.SetKeybinding("Affichage", gocui.KeyArrowDown, gocui.ModNone, flecheBasAffichage); err != nil {
        return err
    }
    
    return nil
}

// Fonction à exécuter lorsque l'utilisateur veut changer de fenêtre (en appuyant sur TAB)
func changerFenetre(g *gocui.Gui, v *gocui.View) error {
    if v != nil {
        switch v.Name() {
            case "Entree":
                if _, err := g.SetCurrentView("Affichage"); err != nil {
                    return err
                }
            case "Affichage":
                if _, err := g.SetCurrentView("Utilisateurs"); err != nil {
                    return err
                }
            case "Utilisateurs":
                if _, err := g.SetCurrentView("Entree"); err != nil {
                    return err
                }
        }
    }
    return nil
}

// Sélectionner l'entrée en-haut de celle actuelle
func flecheHautAccueil(g *gocui.Gui, v *gocui.View) error {
    nom_fenetre := v.Name()
    switch nom_fenetre {
        case "EntreeAdresseIP":
            if _, err := g.SetCurrentView("EntreeNomUtili"); err != nil {
                return err
            }
        case "AccueilDemande":
            vDemande, _ := g.View("AccueilDemande")
            vDemande.Highlight = false
            if _, err := g.SetCurrentView("EntreeAdresseIP"); err != nil {
                return err
            }
    }
    return nil
}

// Sélectionner l'entrée en-dessous de celle actuelle
func flecheBasAccueil(g *gocui.Gui, v *gocui.View) error {
    switch v.Name() {
        case "EntreeNomUtili":
            if _, err := g.SetCurrentView("EntreeAdresseIP"); err != nil {
                return err
            }
        case "EntreeAdresseIP":
            vDemande, _ := g.View("AccueilDemande")
            vDemande.Highlight = true

            if _, err := g.SetCurrentView("AccueilDemande"); err != nil {
                return err
            }
            vDemande.SetCursor(1, 6)
    }
    
    return nil
}

// Fonction à exécuter lorsque l'utilisateur essaye de se connecter.
func entrerAccueil(g *gocui.Gui, v *gocui.View) error {
    // Obtenir les entrées et leur textes
    VNomUtili, _ := g.View("EntreeNomUtili")
    VAdresseIP, _ := g.View("EntreeAdresseIP")
    
    // Enlever les \n de la fin des textes entrés
    NomUtili := strings.TrimRight(VNomUtili.Buffer(), "\n")
    AdresseIP := strings.TrimRight(VAdresseIP.Buffer(), "\n")
    
    // Exécuter seulement si la longueur est plus que 0
    if len(NomUtili) != 0 && len(AdresseIP) != 0 {
        // Les / ne sont pas autorisés dans le nom d'utilisateur
        if strings.Contains(NomUtili, "/") {
            VNomUtili.Clear()
            VNomUtili.SetCursor(0, 0)
            widget.CreerAvis("Nom d'Utilisateur ne peut contenir de '/'!", g)
        
        // Tenter de se connecter    
        } else {
            err := connclient.Connection(AdresseIP, NomUtili, &client, g)
            if err != nil {
                // Créer un avis d'erreur
                VAdresseIP.Clear()
                VAdresseIP.SetCursor(0, 0)
                widget.CreerAvis(err.Error(), g)
 
            } else {
                // Initialiser les fenêtres
                g.SetManagerFunc(initialiserFenetres)
                if err := actions(g); err != nil {
                    return err
                }
                
                err := initialiserFenetres(g)
                if err != nil {
                    return err
                }
                
                if _, err := g.SetCurrentView("Entree"); err != nil {
                    return err
                }
                // Partir le thread de gestion de la connection
                go connclient.GererConnectionServeur(&client, g)
                
            }
        }
    }
    return nil
}

// Monter le curseur de 1 caractère. Sert à voir les messages reçus auparavant
func flecheBasAffichage(g *gocui.Gui, v *gocui.View) error {
    v.MoveCursor(0, 1, false)
    return nil
}

// Descendre le curseur de 1 caractère. Sert à voir les messages reçus actuellement.
func flecheHautAffichage(g *gocui.Gui, v *gocui.View) error {
    v.MoveCursor(0, -1, false)
    return nil
}

// Renvoyer une erreur lorsqu'il faut quitter l'application
func quitter(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}

func main() {
    for {
        // Initialiser l'interface
        g, err := gocui.NewGui(gocui.Output256)
        if err != nil {
            log.Panicln(err)
        }
        defer g.Close()
        
        // Définir quelques propriétés
        g.Cursor = true
        g.SetManagerFunc(DispositionInitiale)
        g.SelFgColor = gocui.ColorMagenta
        g.Highlight = true
        
        if err := actions(g); err != nil {
            log.Panicln(err)
        }
        // Commencer l'application et quitter lors d'une erreur.
        if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
            if err.Error() == "Erreur serveur" {
                precedent = true
                g.Close()
                continue
            } else {
                log.Panicln(err)
            }
        }
        return 
    }
    
}

// Fonction à exécuter pour créer la disposition initiale
func DispositionInitiale(g *gocui.Gui) error {
    maxX, _ := g.Size()
    
    if v, err := g.SetView("AcceuilTitre", maxX/2-10, -1, maxX/2+9, 1); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        v.Frame = false
        fmt.Fprintf(v, "Messagerie Console")
    }
    // Nom d'utilisateur:
    // Adresse IP du serveur:
    if v, err := g.SetView("AccueilDemande", -1, 2, 22, 6); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        v.Frame = false
        v.SelBgColor = gocui.ColorGreen
        v.SelFgColor = gocui.ColorBlack
        
        fmt.Fprintf(v, "Nom d'utilisateur:\nAdresse IP du serveur:\n[x]Se Connecter")
        err := v.SetCursor(1, 2)
        if err != nil {
            log.Panicln(err)
        }
        EntreeNomUtili := widget.CreerEntree("EntreeNomUtili", 18, 2, 22, 20)
        EntreeNomUtili.Disposition(g) 
        
        // 255.255.255.255:xxxx
        EntreeAdresseIP := widget.CreerEntree("EntreeAdresseIP", 22, 3, 22, 20)
        EntreeAdresseIP.Disposition(g)
        
        if _, err := g.SetCurrentView("EntreeNomUtili"); err != nil {
                return err
        }
        
        // S'il y a eu une erreur avec le serveur auparavant, le mentionner
        if precedent == true {
            widget.CreerAvis("Erreur avec la communication au serveur, vous vous êtes déconnectés!", g)
        }
    }
    return nil
}

