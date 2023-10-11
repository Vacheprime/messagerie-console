package widget

import (
    "fmt"
    "log"
    "strings"
    
    "github.com/jroimartin/gocui"
)

// Variable globale qui indique quelle fenêtre était selectionnée avant l'apparition des fenêtres d'avis
var FenetreActuelle string

// Vérifier si une Fenêtre Avis existe autre que celle avec le nom "nom"
func verifierFenetreAvis(nom string, liste []*gocui.View) bool {
    // Itérer à chaque fenêtre présente dans la liste
    for _, fenetre := range liste {
        // Si le nom de la fenêtre commence par Avis et elle est différente du 
        // nom, retourner true
        if strings.HasPrefix(fenetre.Name(), "Avis") && fenetre.Name() != nom {
            return true
        } 
    }
    // Retourner false s'il n'y en a pas
    return false
}

// Retourner la prochaine fenêtre avis autre que celle avec le nom "nom"
func retournerProchainAvis(nom string, liste []*gocui.View) *gocui.View {
    // Itérer à chaque fenêtre présente dans la liste
    for _, fenetre := range liste {
        // Si le nom de la fenêtre commence par Avis et elle est différente du 
        // nom, retourner cette fenêtre
        if strings.HasPrefix(fenetre.Name(), "Avis") && fenetre.Name() != nom {
            return fenetre
        } 
    }
    // Retourner nil si rien n'est trouvé
    return nil
}

// Créer un avis qui apparaît au millieu de l'écran et qui disparait quand 
// "Entrer" est appuyéé
func CreerAvis(message string, g *gocui.Gui) error {
    // Modifier la variable globale FenetreActuelle si aucune fenêtre avis n'existe
    if verifierFenetreAvis("Avis" + message, g.Views()) == false {
        fenetre := g.CurrentView()
        FenetreActuelle = fenetre.Name()
    }
    
    // Déterminer la grandeur de l'écran
    maxX, maxY := g.Size()
    
    // Créer la fenêtre avis avec comme nom: Avis + le message à afficher. Les dimensions sont déterminés à l'aide de la taille de l'écran.
    if v, err := g.SetView("Avis" + message, maxX/2-(len(message)/2)-1, maxY/2-2, maxX/2+(len(message)/2)+1, maxY/2); err != nil {
        if err != gocui.ErrUnknownView {
            log.Panicln(err)
        }
        
        // Imprimer le message à l'écran
        v.Frame = true
        fmt.Fprintf(v, message)
        
        // Fonction à exectuer lorsque la touche "Entrer" est appuyée dans une 
        // fenêtre d'avis
        fermer := func(g *gocui.Gui, VAvis *gocui.View) error {
            
            // Vérifier si une fenêtre d'Avis autre que VAvis existe.
            if verifierFenetreAvis(VAvis.Name(), g.Views()) {
                
                // Vérifier si la fenêtre VAvis n'a pas été supprimée
                if _, err := g.View(VAvis.Name()); err == nil {
                
                    // Supprimer la fenêtre VAvis
                    if err = g.DeleteView(VAvis.Name()); err != nil {
                        return err
                    }
                }
                // Obtenir la prochaine fenêtre d'Avis, la selectionner et la 
                // mettre en premier plan
                NAvis := retournerProchainAvis(VAvis.Name(), g.Views())
                if _, err = g.SetCurrentView(NAvis.Name()); err != nil {
                    return err
                }
                if _, err = g.SetViewOnTop(NAvis.Name()); err != nil {
                    return err
                }
            } else {
                // Vérifier si la fenêtre VAvis existe
                if _, err := g.View(VAvis.Name()); err == nil {
                    
                    // Supprimer la fenêtre VAvis
                    if err := g.DeleteView(VAvis.Name()); err != nil {
                        return err
                    }
                }
                // Selectionner la fenêtre précédente, soit celle qui était 
                // selectionnée avant la création de l'Avis
                if _, err = g.SetCurrentView(FenetreActuelle); err != nil {
                    return err
                }
            }
            return nil
        }
        
        // Définir le raccourci clavier de la fenêtre Avis
        if err := g.SetKeybinding(v.Name(), gocui.KeyEnter, gocui.ModNone, fermer); err != nil {
            return err
        }
        
        // Vérifier si cette fenêtre d'Avis est la première ou non.
        // Mettre l'Avis en arrière plan si oui, selectionner la fenêtre si non.
        if verifierFenetreAvis(v.Name(), g.Views()) {
            if _, err = g.SetViewOnBottom(v.Name()); err != nil {
                return err
            }
        } else {
            if _, err = g.SetCurrentView(v.Name()); err != nil {
                return err
            }
        }
        
    }
    return nil
}
