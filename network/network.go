package network

import (
    "net"
    "encoding/binary"
)

// Formatter en envoyer un message
func EnvoyerMessage(message string, connection net.Conn) error {
    // Créer un entête d'exactement 8 octets qui indique la longueur du message en octets
    entete := make([]byte, 8)
    binary.BigEndian.PutUint64(entete, uint64(len(message)))
    
    // Ajouter l'entête au message
    msg := append(entete, []byte(message)...)
    
    // Envoyer le message
    if _, err := connection.Write(msg); err != nil {
        return err
    }
    return nil
}

// Recevoir un message
func RecevoirMsg(connection net.Conn) (string, error) {
    // Créer l'entête de 8 octets et le décoder du message à recevoir
    entete := make([]byte, 8)
    if _, err := connection.Read(entete); err != nil {
        return "", err
    }
    longueur := binary.BigEndian.Uint64(entete)
    msg := make([]byte, longueur)
    
    if _, err := connection.Read(msg); err != nil {
        return "", err
    }
    // Retourner le message reçu sans l'entête
    return string(msg), nil
}

