package couleurs


// Couleurs premier plan 
var CouleursAVP = map[int]string{
    1 : "\x1b[38;5;124m", // rouge
    2 : "\x1b[38;5;040m", // vert
    3 : "\x1b[38;5;020m", // bleu
    4 : "\x1b[38;5;202m", // orange
    5 : "\x1b[38;5;220m", // jaune
    6 : "\x1b[38;5;200m", // rose
    7 : "\x1b[38;5;055m", // mauve
    8 : "\x1b[38;5;073m", // turquoise
    9 : "\x1b[38;5;123m", // bleu clair
    10 : "\x1b[38;5;225m", // beige
    11 : "\x1b[0m", // normal
}

// Couleurs arrière plan
var CouleursARP = map[int]string{
    1 : "\x1b[48;5;124m", // rouge
    2 : "\x1b[48;5;40m", // vert
    3 : "\x1b[48;5;20m", // bleu
    4 : "\x1b[48;5;202m", // orange
    5 : "\x1b[48;5;220m", // jaune
    6 : "\x1b[48;5;200m", // rose
    7 : "\x1b[48;5;55m", // mauve
    8 : "\x1b[48;5;73m", // turquoise
    9 : "\x1b[48;5;123m", // bleu clair
    10 : "\x1b[48;5;225m", // beige
    11 : "\x1b[0m", // normal
}
