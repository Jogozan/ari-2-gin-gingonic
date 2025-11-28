package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// albums slice to seed record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	router.Run("localhost:8080")
}

// getAlbums responds with the list of all albums as JSON.
func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON received in the request body.
func postAlbums(c *gin.Context) {
	var newAlbum album

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// Add the new album to the slice.
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

// getAlbumByID locates the album whose ID value matches the id
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of albums, looking for
	// an album whose ID value matches the parameter.
	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

/* Fil rouge : Pokédex “avancé” en 4 h

On part de :

    API JSON CRUD complète Pokemon.

    Vues HTML basiques : liste + formulaire de création.
    Les étudiant·es vont enrichir ce socle avec logique métier, middlewares, groupes, validations avancées et petites features “fun”.

1. Rappels + lecture du code existant (20–30 min)

Objectifs : comprendre le socle fourni et identifier les points d’extension.

Contenu :

    Lecture rapide du routeur : où sont définies les routes API, les routes HTML, les templates chargés.

    Explication de la struct Pokemon et des handlers actuels (liste, détail, création).

    Point d’entrée sur la notion de groupe (/api/v1) et de contexte (*gin.Context) pour préparer la suite.

2. Logique métier “Pokémon” (évolution, types, stats) (45 min)

Objectif : montrer que les handlers Gin peuvent embarquer une vraie logique métier, pas juste du CRUD.

Idées d’extensions (tu peux en choisir 2–3) :

    Route d’évolution :

        POST /api/v1/pokemons/:id/level-up qui augmente Level et HP selon une règle (ex. +1 niveau = +10 HP, ou dépendant du type).

        Gestion d’erreurs : 404 si ID inconnu, 400 si niveau max atteint.

    Calcul de “power” :

        Champ calculé (non stocké) envoyé en JSON et affiché dans le template HTML (par ex. Power = Level * HP ou en fonction du type).

    Filtrage serveur :

        Sur l’API JSON : GET /api/v1/pokemons?type=fire&minLevel=10.

        Sur la vue HTML : un petit formulaire en GET qui utilise ces query params pour filtrer la liste.

On insiste sur :

    Lecture des query params via le contexte.

    Structs de réponse dédiées si besoin (ajouter Power sans casser la struct de stockage).

3. Middlewares utiles + sous-groupe “admin” (45 min)

Objectif : exploiter les middlewares Gin et les groupes pour illustrer des cas concrets.

Ajouts proposés :

    Middleware de rate-limiting simple ou de “fatigue du serveur” :

        Par exemple, limiter les appels à la route level-up ou ajouter un délai artificiel pour montrer l’impact.

    Groupe /api/v1/admin avec :

        POST /api/v1/admin/pokemons/:id/level-up protégé par un middleware d’“auth” très simple : vérification d’un header ou d’un cookie.

    Middleware de logging enrichi :

        Ajout dans les logs de la “dresseur·euse” (valeur lue dans un header ou paramètre) et du Pokémon ciblé.

But pédagogique :

    Bien distinguer middlewares globaux, de groupe et de route.

    Montrer comment attacher/propager des infos dans le contexte pour les récupérer plus tard dans le handler.

4. Validation avancée + retours d’erreurs propres (40–45 min)

Objectif : dépasser le “BindJSON basique” pour montrer la puissance de la validation et des réponses structurées.

Travail proposé :

    Structs d’input avec tags de validation plus riches :

        Name non vide, longueur max.

        Type dans une liste autorisée (Fire, Water, Grass, etc.) – validation custom.

        Level min/max.

    Fonction helper pour les réponses :

        Standardiser les réponses JSON en { data: ..., error: ... }.

        Renvoyer des listes d’erreurs de validation lisibles pour le front.

    Intégration dans les templates HTML :

        En cas d’erreur de validation en POST HTML, réafficher le formulaire avec les messages d’erreurs à côté des champs.

Ce bloc met en avant la séparation input / modèle / vue et l’importance d’une couche de validation propre.
5. Templates HTML enrichis + petite “gamification” (40–45 min)

Objectif : rendre la partie graphique plus vivante, tout en restant côté serveur avec Gin.

Idées d’amélioration :

    Page Pokédex “tableau de bord” :

        Affichage des Pokémon triés par Power ou Level.

        Badges de type (couleur différente selon le type).

    Actions depuis le HTML :

        Bouton “Level up” sur chaque ligne qui appelle la route correspondante (méthode POST via un formulaire caché).

        Peut-être un bouton “Relâcher” (DELETE) avec confirmation simple côté front.

    Une page “stats dresseur” qui :

        Compte le nombre de Pokémon par type et les affiche dans un petit tableau ou liste.

        Utilise une route GET /pokemons/stats qui calcule ces infos côté serveur.

L’intérêt est de montrer que les mêmes données et règles métiers servent à la fois en JSON et dans les vues HTML.
6. Optionnel si temps : petite intégration externe ou tests (20–30 min)

Deux pistes, selon ton public :

    Intégration externe légère : appeler une API publique Pokémon pour récupérer une info (type officiel, sprite) et l’afficher dans le template (en gardant la logique côté serveur).

    Tests d’handlers Gin : écrire 1–2 tests simples sur un handler clé (liste ou level-up) pour montrer comment instancier un routeur, faire une requête de test et vérifier la réponse.

*/
