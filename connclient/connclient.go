package connclient

import (
    "os"
    "log"
    "errors"
    "strings"
    "fmt"
    "net"
    "strconv"
    "math"
    "unicode/utf8"
    
    "github.com/Vacheprime/gochat/widget"
    "github.com/Vacheprime/gochat/network"
    "github.com/Vacheprime/gochat/couleurs"
    
    "github.com/jroimartin/gocui"
    
)

// Définition de la structure Client
type Client struct {
    NomUtilisateur string
    AdresseIP string
    Connection net.Conn
    CouleurID int64
}

func Connection(AdresseIP, NomUtili string, client *Client, g *gocui.Gui) error {
    
    // Essai de connection au serveur avec l'AdresseIP 
    connection, err := net.Dial("tcp", AdresseIP)
    if err != nil {
        // Définir le type d'erreur
        if erri, ok := err.(*net.OpError); ok {
            // Déterminer une erreur d'adresse invalide
            if _, ok := erri.Err.(*net.AddrError); ok {
                return errors.New("Adresse IP invalide!")
                
            // Déterminer une erreur de connection au serveur
            } else if _, ok := erri.Err.(*os.SyscallError); ok {
                return errors.New("Erreur de connection au serveur!")
            }
        }
        // Si l'erreur ne peut être définie, retourner une erreur disant qu'elle
        // est inconnue.
        
        return errors.New("Erreur inconnue est survenue!")
    }
    
    // Actualisation de la variable client avec les bonnes valeures
    client.NomUtilisateur = NomUtili
    client.AdresseIP = AdresseIP
    client.Connection = connection
    
    return nil
}

// Créer un avis et le montrer à partir du thread du gestion de la connection
func montrerAvis(avis string, g *gocui.Gui) error {
    g.Update(func (g *gocui.Gui) error {
        widget.CreerAvis(avis, g)
        return nil
    })
    return nil
}

// Imprimer un message à la fenetre Affichage
func imprimerFenetreAffichage(msg string, fenetre *gocui.View,  g *gocui.Gui) error {
    g.Update(func (g *gocui.Gui) error {
        
        // Imprimer le message à l'écran
        fmt.Fprintf(fenetre, msg + "\n")

        // Actualiser la position du curseur
        _, y := fenetre.Cursor()
        maxX, maxY := fenetre.Size()
        NbrLignesMsg := int(math.Ceil(float64(utf8.RuneCountInString(msg) - 6) / float64(maxX)))
        totalLignes := 0
        
        
        // Calculer le nombre de lignes imprimées à l'écran avant l'impression
        for _, ligne := range fenetre.ViewBufferLines() {
            longueurLigne := utf8.RuneCountInString(ligne)
            totalLignes += int(math.Ceil(float64(longueurLigne - 6) / float64(maxX)))
        }
        // Selon la position du curseur et du nombre de lignes du Message reçu, faire bouger le curseur et le texte vers le bas
        if totalLignes == y || totalLignes - 1 == y || totalLignes - 2 == y {
            // Si le message est plus grand que l'écran, faire défiler l'écran pour montrer le message au complet
            if NbrLignesMsg > maxY - y {
                _, oy := fenetre.Origin()
                fenetre.SetOrigin(0, oy + NbrLignesMsg)

            } else {
                fenetre.MoveCursor(0, NbrLignesMsg, true)
            }
            
        // Continuer de faire apparaître les messages si le curseur est rendu à la fin de l'écran
        } else if maxY == y || maxY - 1 == y {
            if NbrLignesMsg > 1 {
                _, oy := fenetre.Origin()
                fenetre.SetOrigin(0, oy + NbrLignesMsg)
            } else {
                fenetre.MoveCursor(0, NbrLignesMsg, true)
            }
            
        }
        

        
        return nil
    })
    return nil
}

// Imprimer la liste d'utilisateurs en mettant les couleurs appropriées
func imprimerListeUtilisateurs(listeUtilisateurs []string, CouleurID int64, fenetre *gocui.View, g *gocui.Gui) {
    g.Update(func (g *gocui.Gui) error {
        
        // Passer par tous les utilisateurs connectés
        for _, utilisateur := range listeUtilisateurs {
            // Définir l'information de chaque utilisateur reçu
            // Format:
            // NomUtilisateur|CouleurID
            infoUtilisateur := strings.Split(utilisateur, "|")
            nomUtili := infoUtilisateur[0]
            couleurUtili, _ := strconv.ParseInt(infoUtilisateur[1], 10, 64)
            
            // Vérifier si l'utilisateur est nous
            if couleurUtili == CouleurID {
                
                // Imprimer avec une couleur d'arrière plan
                fmt.Fprintf(fenetre, couleurs.CouleursARP[int(couleurUtili)] + nomUtili + couleurs.CouleursARP[11] + "\n")
            } else {
                
                // Imprimer avec une couleur d'avant plan si ce n'est pas nous
                fmt.Fprintf(fenetre, couleurs.CouleursAVP[int(couleurUtili)] + nomUtili + couleurs.CouleursAVP[11] + "\n")
            }
        }
        return nil
    })
}

func GererConnectionServeur(client *Client, g *gocui.Gui) {
    // Envoyer le nom d'utilisateur
    if err := network.EnvoyerMessage(client.NomUtilisateur, client.Connection); err != nil {
        log.Panicln(err)
    }
    // Définir quelques variables
    VUtilisateurs, _ := g.View("Utilisateurs")
    VAffichage, _ := g.View("Affichage")
    
    // Recevoir l'identificateur d'utilisateur
    if IDstr, err := network.RecevoirMsg(client.Connection); err != nil {
        log.Panicln(err)
    } else {
        client.CouleurID, _ = strconv.ParseInt(IDstr, 10, 64)
    }
    // Boucle infinie qui gère les messages envoyés par le serveur
    for {
        
        // Attendre pour les messages du serveur
        msg, err := network.RecevoirMsg(client.Connection)
        if err != nil {
            montrerAvis("Erreur de reception d'un message!", g)
            break
        }
        
        // Vérifier le type de message reçu
        if strings.HasPrefix(msg, "msg:") {
            
            // Format du message: 
            // msg:NomUtilisateur|CouleurID|Contenu
            msg = strings.TrimPrefix(msg, "msg:")
            infoMessage := strings.Split(msg, "|")
            
            // S'assurer que le message n'est pas un message du serveur du format suivant:
            // msg:Contenu
            
            if len(infoMessage) > 1 {
                couleurUtili, _ := strconv.ParseInt(infoMessage[1], 10, 64)
                
                if couleurUtili == client.CouleurID {
                    msg = couleurs.CouleursARP[int(client.CouleurID)] + infoMessage[0] + couleurs.CouleursARP[11] + ": " + infoMessage[2]
                } else {
                    msg = couleurs.CouleursAVP[int(couleurUtili)] + infoMessage[0] + couleurs.CouleursAVP[11] + ": " + infoMessage[2]
                }
            }
            
            // Imprimer le message à la fenêtre de messages
            imprimerFenetreAffichage(msg , VAffichage, g)
        } else if strings.HasPrefix(msg, "utilisateurs:") {
 
            // Effacer le contenu de la fenêtre
            VUtilisateurs.Clear()
            
            // Format du message:
            // utilisateurs:NomUtilisateur|CouleurID/NomUtilisateur2|CouleurID2/{....}
            utilisateurs := strings.TrimPrefix(msg, "utilisateurs:")
            listeUtilisateurs := strings.Split(utilisateurs, "/")
            
            // Imprimer la liste d'utilisateurs à l'écran
            imprimerListeUtilisateurs(listeUtilisateurs[:len(listeUtilisateurs) - 1], client.CouleurID, VUtilisateurs, g)
        }
    }
    
    // Retourner une erreur de déconnection
    g.Update(func (g *gocui.Gui) error {
        return errors.New("Erreur serveur")
    })
}
