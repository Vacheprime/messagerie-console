package main 

import (
    "log"
    "net"
    "strings"
    "fmt"
    "io"
    "strconv"
    
    "github.com/Vacheprime/gochat/network"
    "github.com/Vacheprime/gochat/couleurs"
)

// Définition de constantes qui serviront à initialiser le serveur
const (
    IP string = "127.0.0.1"
    PORT string = "5555"
    TYPE string = "tcp"
)

// Définition de la liste de clients et du nombre de clients connectés
var ListeClients [10]Client
var NbrClients uint

// Définition de la structure client
type Client struct {
    NomUtilisateur string
    Connection net.Conn
    ID int
    CouleurAVP string
    CouleurARP string
    CouleurID int
}

// Enlever un client de la liste de clients grâce à son identificateur
func enleverClient(listeClients *[10]Client, ID int) {
    for index, client := range listeClients {
        if client.ID == ID {
            listeClients[index] = Client{"", nil, 0, "", "", 0}
            break
        }
    }
}

// Ajouter un client à la liste de clients connectés
func ajouterClient(listeClients *[10]Client, clientAjouter Client) int{
    for index, client := range listeClients {
        if client.Connection == nil {
            // Actualiser les valeurs de couleur du client 
            clientAjouter.CouleurAVP = couleurs.CouleursAVP[index + 1]
            clientAjouter.CouleurARP = couleurs.CouleursARP[index + 1]
            clientAjouter.CouleurID = index + 1
            listeClients[index] = clientAjouter
            return index
        }
    }
    return 0
}

// Diffuser un message à tous les clients connectés
func diffuser(msg string) error {
    for _, utilisateur := range ListeClients {
        if utilisateur.Connection != nil {
            if err := network.EnvoyerMessage(msg, utilisateur.Connection); err != nil {
                return err
            }
        }
    }
    return nil
}

// Fonction qui gère les échanges avec le client en continu
func GererConnection(connection net.Conn, client *Client) {
    
    // Recevoir le nom d'utilisateur du client
    NomUtili, err := network.RecevoirMsg(connection)
    log.Printf("Nom Utilisateur reçu: %v", NomUtili)
    
    if err != nil {
        log.Printf("Erreur de reception du nom d'utilisateur")
        connection.Close()
    }
    client.NomUtilisateur = NomUtili
    
    // Envoyer l'index de couleur au client
    if err := network.EnvoyerMessage(strconv.Itoa(client.CouleurID), connection); err != nil {
        log.Println("Erreur en envoyant l'index de couleur de l'utilisateur.")
    }
    
    // Envoyer la liste d'utilisateurs connectés au client connecté
    liste := ""
    for _, utilisateur := range ListeClients {
        if utilisateur.Connection != nil{
            liste += utilisateur.NomUtilisateur + "|" + strconv.Itoa(utilisateur.CouleurID) + "/"
        }
    }
    
    // Diffuser le nom d'utilisateur de l'utilisateur qui vient de se connecter
    if err := diffuser("utilisateurs:" + liste); err != nil {
        log.Println("Erreur en diffusant la liste d'utilisateurs.")
    }
    
    // Diffuser un message d'accueil à tout le monde
    if err := diffuser(fmt.Sprintf("msg:[SERVEUR] %v%v%v vient de se connecter!", client.CouleurAVP, client.NomUtilisateur, couleurs.CouleursAVP[11])); err != nil {
        log.Printf("Erreur d'envoi du message de notification de connection")
    }
    
    // Boucle infinie qui reçoit les messages du client
    for {
        
        // Recevoir le message du client
        msg, err := network.RecevoirMsg(connection)
        if err != nil {
            if err == io.EOF {
                log.Println("Client déconnecté")
                break
            }
            log.Println(err)
            break
        }
        // Vérifier le type de message et agir en conséquence
        if strings.HasPrefix(msg, "msg:") {
            if err := diffuser(msg); err != nil {
                log.Printf("Erreur en diffusant le message: %v", msg)
            }
        }
    }
    
    // Créer la nouvelle liste d'utilisateurs
    liste = ""
    for _, utilisateur := range ListeClients {
        if utilisateur.Connection != nil{
            liste += utilisateur.NomUtilisateur + "|" + strconv.Itoa(utilisateur.CouleurID) + "/"
        }
    }
    
    // Diffuser la nouvelle liste d'utilisateurs connectés
    if err := diffuser("utilisateurs:" + liste); err != nil {
        log.Println(err)
    }
    
    // Diffuser un message indiquant la déconnection du client
    diffuser(fmt.Sprintf("msg:[SERVEUR] %v%v%v vient de se déconnecter!", client.CouleurAVP, client.NomUtilisateur, couleurs.CouleursAVP[11]))
    
    // Enlever le client de la liste de clients, fermer sa connection et diminuer le compteur de clients
    NbrClients--
    enleverClient(&ListeClients, client.ID)
    connection.Close()
}

func main() {
    // Initialiser le serveur
    
    AdresseServeur := IP + ":" + PORT
    serveur, err := net.Listen(TYPE, AdresseServeur)
    if err != nil {
        log.Panicln(err)
    }
    defer serveur.Close()
    log.Printf("Serveur initialisé sur %v, en attente de clients...\n", AdresseServeur)
    
    // Initialiser le compteur d'identificateurs
    compteur := 0
    // Boucle infinie qui accepte les connections
    for {
        if NbrClients < 10 {
            // Attendre pour une connexion et l'accepter
            
            connection, err := serveur.Accept()
            if err != nil {
                log.Panicln(err)
            }
            log.Println("Client connecté!")
            index := ajouterClient(&ListeClients, Client{"", connection, compteur, "", "", 0})
            compteur++
            NbrClients++
            // Partir un nouveau fil d'exécution/thread gérant les échanges avec 
            // la connection  
            go GererConnection(connection, &ListeClients[index])
        }
    }
    return 
}
