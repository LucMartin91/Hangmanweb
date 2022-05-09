package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type HangManData struct {
	Word             string // Word composed of '_', ex: H_ll_
	ToFind           string // Final word chosen by the program at the beginning. It is the word to find
	Attempts         int    // Number of attempts left
	HangmanPositions string // It can be the array where the positions parsed in "hangman.txt" are stored
	index_tab_pendu  int    // ça va être là où on en est dans hangman positions.
}

func ToUpper(s string) string { // Ici on va créer une fonction toUpper qui va nous permettre de mettre les lettres minuscules dd'une string en majuscules.
	compteur := ""
	ch := []rune(s)
	for i := 0; i <= len(s)-1; i++ {
		if ch[i] >= 'a' && ch[i] <= 'z' {
			compteur = compteur + string(ch[i]-32)
		} else {
			compteur = compteur + string(ch[i])
		}
	}
	return compteur
}

func motatrouver() string { // la string qu'on va return sera un mot aléatoirement choisi dans words.txt
	var tabmots []string           // là où on va stocker les mots de words.txt
	f, err := os.Open("words.txt") // on va lire words.txt
	if err != nil {                // si erreur (words.txt vide ou autre)
		log.Fatal(err) // on va afficher l'erreur corespondante
	}

	scanner := bufio.NewScanner(f) // on va scanner le contenu de words.txt qui a été stocké dans f plus haut
	scanner.Split(bufio.ScanLines) // on va dire à notre scanner de scanner lignes par lignes (1 ligne = une place)
	for scanner.Scan() {           // tant que le scanner à du contenu à scanner
		tabmots = append(tabmots, (scanner.Text())) // on ajoute la ligne à  tabmots tant qu'il y'a des lignes dans words.txt
	}

	max := len(tabmots) // ici on définit les limites du tableau qui contient les mots pour en choisir un aléatoirement dedans
	min := 0
	rand.Seed(time.Now().UnixNano())
	rand := rand.Intn(max - min)  // On va stocker un int aléatoire choisi par la fonction rand.Intn situé entre les limites définies ici par max et min
	mot := ToUpper(tabmots[rand]) // Ici on va user de notre fonction ToUpper afin de mettre en majuscule le mot choisi dans tabmots par un int aléatoirement défini précédemment
	return mot                    // On retourne la variable mot dans laquelle est stocké notre mot choisi aléatoirement et mit ensuite en majuscules par ToUpper
}

func débutjeu(s string) string { // On va créer une fonction qui va servir a initialiser notre jeu.
	var motcaché []string // tableau contenant le mot dont certaines lettres vont être remplacées par des "_"
	for range s {
		motcaché = append(motcaché, "_") // ici on ajoute à motcaché un "_" pour chaque caractère de s
	}
	for i := 0; i < (len(s)/2)-1; i++ { // ici on va faire réveler de manière aléatoire un ou plusieurs des caractères du motcaché
		rand.Seed(time.Now().UnixNano())
		pos := rand.Intn(len(s))
		motcaché[pos] = string(s[pos])
	}
	res := "" // ici on créé juste une variable de type string à laquelle on va ajouter tous les caractères de notre motcaché fini
	for i := 0; i < len(motcaché); i++ {
		res = res + motcaché[i]
	}
	return res
}

func islettre(motcomplet string, lettrecherché string, motavecles_ string) bool {
	tab_motcomplet := []rune(motcomplet)  // on initialise un tableau de runes du mot complet quon cherche
	tabentrée := []rune(lettrecherché)    // On créé un tableau de runes contenant l'entrée utilisateur
	tabmotavecles_ := []rune(motavecles_) // Tableau de runes de notre mot caché
	if lettrecherché == motcomplet {
		return true
	}
	nb_de_fois_lalettre := 0                   // Compteur de fois où la lettre est dans notre mot
	for i := 0; i < len(tab_motcomplet); i++ { // On fait une boucle ici pour voir combien de fois la lettre donnée est présente dans notre mot complet
		if tabentrée[0] == tab_motcomplet[i] {
			nb_de_fois_lalettre++
		}
	}
	for i := 0; i < len(tab_motcomplet); i++ { // ici on va vérifier si on a déjà découvert la lettre qu'on cherche dans le mot
		if tab_motcomplet[i] == tabentrée[0] && len(tabentrée) == 1 {
			compteur := 0 // ce compteur sert à voir si on a autant de fois la lettre cherchée dans le mot caché que dans le mot complet
			for i := 0; i < len(tabmotavecles_); i++ {
				if tabentrée[0] == tabmotavecles_[i] {
					compteur++
				}
				if compteur == nb_de_fois_lalettre { // si on a déjà découvert tous les "a" par exemple, on retourne donc false
					return false
				}
			}
			return true // sinon, on retourne true, la lettre entrée par l'utilisateur n'a pas été découverte et elle fait partie du mot qu'on cherche.
		}
	}
	return false
}

func replace(entreeutilisateur string, motcomplet string, gruyère string) string { // fonction qui va nous servir à modifier le motcaché tout au long du jeu
	tabgruyère := []rune(gruyère)          // tableau de runes de notre mot caché
	tabentrée := []rune(entreeutilisateur) // tableau de runes de notre entrée utilisateur
	for i := 0; i < len(motcomplet); i++ { // on boucle sur la len de motcomplet
		if entreeutilisateur[0] == motcomplet[i] && len(entreeutilisateur) == 1 { // si l'entrée utilisateur est égale à une ou plusieurs lettres du mot complet et que le scanner détecte bien une seule lettre et pas un mot dans l'input après le choose :
			tabgruyère[i] = tabentrée[0] // alors on modifie le ou les "_" de tabgruyères par la lettre en question selon sa ou ses positions dans le mot complet.
		} else if entreeutilisateur == motcomplet { // Sinon, si l'entréeU correspond à notre motcomplet
			tabgruyère = tabentrée // alors on replace chaque lettre par les lettres du motcomplet.
		}
	}
	return string(tabgruyère) // Enfin on retourne la string correspondante à notre résultat.
}

func main() {
	var hangman HangManData                 // on créé une variable de type HangManData(voir en haut)
	hangman.Attempts = 10                   // On initialise le nomrbe d'erreurs permises à 10
	hangman.ToFind = motatrouver()          // sinon on initialise le jeu à 0 avec la fonction motatrouver expliquée plus haut (elle trv un mot o pif dans words.txt qui va etre celui à trouver)
	hangman.Word = débutjeu(hangman.ToFind) // on va créer notre mot caché par rapport au résultat de motatrouver() qu'on a fait au dessus
	hangman.index_tab_pendu = 0
	f, err := os.Open("hangman.txt") // on ouvre hangman.txt et on stock son contenu dans f
	if err != nil {
		log.Fatal(err)
	}
	//var tabjoser []string          // ici on va créer un tab de string qui va stocker les positions de josé
	scanner := bufio.NewScanner(f) // On initialise un scanner sur f qui contient ttes les pos de notre pendu
	scanner.Split(bufio.ScanLines) // ici on dit au scanner de scanner ligne par ligne
	paragraphe := ""               // On créé une variable paragraphe vide qui va stocker ligne par ligne notre position de hangman
	for scanner.Scan() {           // on boucle tant que le scanner scan
		paragraphe = paragraphe + scanner.Text() + "\n" // adds the value of scanner (that contains the characters from StylizedFile) to source
	}

	hangman_positions := strings.Split(paragraphe, "=========") // On va split la string paragraphe à chaque fois que string.Split détecte "========="

	for index, v := range hangman_positions { // On va ensuite boucler pour ajouter chaque string petit à petit à chaque position dans notre tableau
		hangman_positions[index] = v + "========="
	}
	hangman_positions = hangman_positions[:len(hangman_positions)-1] // Ici on supprime la dernière position de hangmanposition car c une position vide
	// là on met la variable qui va nous servir à aller chercher les positions du pendu dans notre tableau du pendu en fonction des attemps restantes
	hangman.HangmanPositions = ""
	tmpl, _ := template.ParseGlob("./www/*.gohtml")
	fmt.Println("\n", "LOGS SERVEUR :", "\n\n", "Position du pendu actuelle : ", hangman.HangmanPositions, "\n", "position dans tableau pendu : ", hangman.index_tab_pendu, "\n", "Chances restantes : ", hangman.Attempts, "\n", "Mot caché : ", hangman.Word, "\n", "Mot à trouver actuellement : ", hangman.ToFind) // Ici on affiche nos logs dans la console pour le début de jeu affiché sur la page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {                                                                                                                                                                                                                                                // Ici on va définir le chemin où on va jouer ( ici "/" qui est la racine de notre serveur)
		if r.Method != http.MethodPost { // Ici on va régler nos erreurs de quand le scanner est vide pour lancer quand même sinon ça casse tout
			tmpl.ExecuteTemplate(w, "index", hangman)
			return
		}
		scanner_input := r.FormValue("scanner")                 // on stock le résultat du scanner nommé "scanner" dans mon html dans cette variable
		if scanner_input[0] >= 'a' && scanner_input[0] <= 'z' { // si la lettre ecrite est en minuscule, on la met en maj avec toupper
			scanner_input = ToUpper(scanner_input)
		}
		if islettre(hangman.ToFind, scanner_input, hangman.Word) { // si la fonction islettre renvoie true
			hangman.Word = (replace(scanner_input, hangman.ToFind, hangman.Word)) // Là on fait évoluer notre mot caché avec replace
		} else {
			hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
			if len(scanner_input) == 1 { // Si l'entréeU correspond à une lettre et pas à un mot
				hangman.Attempts--         // Sinon, on enlève une chance au joueur
				if hangman.Attempts == 0 { // Si le joueur n'a plus de chance
					hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
					// on print la dernière position du pendu, puis GAME OVER
					// et on sort de notre fonction.
				} else {
					hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
					// fmt.Print("Not present in the word, ", hangman.Attempts, " attempts remaining", "\n", hangman.HangmanPositions, "\n") // sinon on print le message correspondant à une mauvaise réponse
					hangman.index_tab_pendu++ // On augmente ensuite l'index du tableau du pendu pour print la prochaine position à la prochaine erreur.
				}
			} else if len(scanner_input) > 1 {
				if hangman.Attempts-2 < 0 { // si normalement on doit faire attempts -2 mais qu'il reste qu'une attempt par exemple, alors pour éviter de merder :
					hangman.Attempts--         // on enlève une chance au joueur au lieu de 2 car si il reste 1 attempts et si on fait -2 on se retrouverait à Attempts = -1, le programme ne passerait donc pas par la condition du dessous qui se déclanche si on n'a plus d'attempts, c'est à dire attempts = 0.
					if hangman.Attempts == 0 { // Si le joueur n'a plus de chance
						hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
					}
				} else {
					hangman.index_tab_pendu++               // On incrémente une première fois l'index de notre pendu (il faudra le faire deux fois car on doit avoir 2 chances en moins dans le cas où la len de l'entréeU est supérieure à 1 (donc le prgramme détecte un mot)
					hangman.Attempts = hangman.Attempts - 2 // On décrémente 2 fois les attempts car on a deux chances en moins dans ce cas
					if hangman.Attempts > 0 {               // si les attempts sont supérieurs à 0 (si il nous reste au moins 1 chance)
						hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
						//	fmt.Print("Not present in the word, ", hangman.Attempts, " attempts remaining", "\n", hangman.HangmanPositions, "\n") // on print le message correspondant à une mauvaise réponse
						hangman.index_tab_pendu++ // On incrémente une deuxième fois l'index du tableau après le print pour préparer la prochaine output
					} else { // Sinon on va print notre game over et quitter le jeu du pendu.
						hangman.HangmanPositions = hangman_positions[hangman.index_tab_pendu]
						//	fmt.Print(hangman.HangmanPositions, "\nGAME OVER")
						//	return
					}
				}
			}
		}
		fmt.Println("\n\n", "Position du pendu actuelle : ", hangman.HangmanPositions, "\n", "position dans tableau pendu : ", hangman.index_tab_pendu, "\n", "Chances restantes : ", hangman.Attempts, "\n", "Mot caché : ", hangman.Word, "\n", "Mot à trouver actuellement : ", hangman.ToFind, "\n", "Keylog :", scanner_input) // Ici on affiche nos logs dans la console à chaque refresh de la page
		tmpl.ExecuteTemplate(w, "index", hangman)                                                                                                                                                                                                                                                                                   //On va exécuter notre template en lui implantant notre structure hangman pour pouvoir l'afficher dans notre html
	})
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) { // ici on va créer une page reload sur notre serveur qui si on s'y rend, va réinitialiser notre jeu pour ensuite nous renvoyer dans une nouvelle partie.

		hangman.HangmanPositions = ""
		hangman.Attempts = 10                         // On initialise le nomrbe d'erreurs permises à 10
		hangman.ToFind = motatrouver()                // sinon on initialise le jeu à 0 avec la fonction motatrouver expliquée plus haut (elle trv un mot o pif dans words.txt qui va etre celui à trouver)
		hangman.Word = débutjeu(hangman.ToFind)       // on va créer notre mot caché par rapport au résultat de motatrouver() qu'on a fait au dessus
		hangman.index_tab_pendu = 0                   // On reset l'index du tableau de notre pendu
		http.Redirect(w, r, "/", http.StatusSeeOther) // Enfin, on redirige l'utilisateur une fois ces modifications appotrtées à nos variables dans la racine serveur où on affiche notre jeu.
	})
	http.ListenAndServe("localhost:555", nil) // ici on héberge notre serveur en local sur le port 555.

}
