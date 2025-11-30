# TP — Templates HTML enrichis + petite « gamification » (40–45 min)

Objectif : enrichir l'interface HTML côté serveur avec Gin pour rendre la partie graphique plus vivante et montrer que les mêmes règles métiers servent en JSON et en HTML.

Ce que j'ai implémenté :

- Route HTML : tri par "level" ou "power" sur la page index
- Badges colorés par type (palette simple côté serveur)
- Boutons/actions côté HTML : "Level up" (POST) et "Relâcher" (POST + confirmation)
- Page de statistiques : `/pokemons/stats` qui agrège le nombre de pokémon par type
- Documentation détaillée et suggestions d'exercices pour un TP de ~40–45 min

---

Durée prévue : 40–45 min (expliquée en étapes ci-dessous)

## Fichiers modifiés / ajoutés

- `pokemon/handler_html.go` — ajout des routes HTML et handlers :
  - `POST /pokemons/:id/level-up` (level up depuis HTML)
  - `POST /pokemons/:id/release` (relâcher / supprimer depuis HTML)
  - `GET /pokemons/stats` (statistiques dresseur)
  - support d'un param `sort=level|power` pour la liste
  - envoi d'une petite map `type_colors` vers les templates pour colorer les badges

- `templates/pokemons_index.tmpl` — mise à jour :
  - nouvelles actions (Level up, Relâcher)
  - affichage du Level
  - contrôles de tri (Original, Par level, Par power)
  - affichage d'un message (retours d'actions)
  - badges colorés par type

- `templates/pokemons_detail.tmpl` — mise à jour :
  - badges colorés par type
  - boutons Level up et Relâcher depuis la page détail

- `templates/pokemons_stats.tmpl` — NOUVEAU : affiche le nombre de Pokémon par type

- `pokedex/TP5-templates-gamification.md` — ce document

---

## Déroulé pédagogique (étapes / minuteur)

1) (5 min) Introduction et exploration
   - Montrer l'ancienne page HTML simple : `templates/pokemons_index.tmpl` et `pokemon/handler_html.go`.
   - Expliquer que l'objectif est d'ajouter des actions côté HTML qui utilisent la même logique métier que l'API JSON (`LevelUp`, `Delete`).

2) (10–15 min) Ajout d'actions côté HTML
   - Implémenter `POST /pokemons/:id/level-up` : appelle la fonction `LevelUp(id, n)` côté repository.
   - Implémenter `POST /pokemons/:id/release` : appelle `Delete(id)`.
   - Après exécution, rediriger vers `/pokemons` avec un message court (query param `msg` utilisé par le template).
   - Exercices bonus : limiter les level-up à certaines conditions (max par jour), ou ajouter un petit middleware "admin" (démo). 

3) (8–10 min) Améliorer l'affichage (templates)
   - Ajouter colonne Level dans la table du Pokédex.
   - Ajouter un bouton "Level up" par ligne qui soumet un formulaire POST.
   - Ajouter badge coloré par type : passer une map `type_colors` depuis le handler vers le template et utiliser `index` pour récupérer la couleur.
   - Ajouter simple confirmation HTML/JS pour la suppression / relâcher d'un Pokémon.

4) (6–8 min) Page statistiques du dresseur
   - Créer `GET /pokemons/stats` qui parcourt la liste des Pokémon et calcule un mapping `type -> count`.
   - Afficher le résultat dans `templates/pokemons_stats.tmpl` (table simple).

5) (3–5 min) Vérification et démonstration finale
   - Lancer le serveur et parcourir les pages.
   - Montrer la redirection et les messages après level-up / relâcher.
   - Montrer la page `pokemons/stats` et expliquer comment la logique métier est réutilisée.

---

## Comment tester et lancer

Depuis le dossier `pokedex` :

```bash
# démarrer
go run .
# ouvrir http://localhost:8080/pokemons
```

Actions disponibles côté HTML :
- Depuis la liste : cliquer `Level up` pour augmenter d'un niveau.
- Depuis la liste : cliquer `Relâcher` (avec confirmation) supprime le Pokémon.
- Tri : `/pokemons?sort=level` (par level), `/pokemons?sort=power` (par une métrique simple combinant stats + base XP)
- Statistiques : `/pokemons/stats` affiche le nombre de Pokémon par type

Notes pratiques :
- Les actions HTML utilisent des POST simples (forms) pour rester compatibles avec les formulaires HTML.
- La logique métier (LevelUp, Delete) reste centralisée dans `pokemon/repository.go` et est réutilisée par les handlers JSON et HTML.

---

## Suggestions pour exercices / extensions

- Implémenter un petit middleware d'auth (ou ré-utiliser le groupe admin existant pour l'API) et restreindre le level-up depuis l'interface HTML (ex : bouton visible seulement si X-Trainer présent).
- Garder un historique des actions (journaling) pour faire une mini-visualisation.
- Ajouter un badge "rare" si la somme des stats > threshold.
- Remplacer la table par une grille plus visuelle ou ajouter pie charts sur la page stats (via une librairie JS incluse côté HTML).

---

Si vous voulez, je peux maintenant :
- ajouter quelques tests unitaires pour les nouveaux endpoints HTML (si vous voulez garder un coverage basique),
- améliorer l'apparence / CSS pour rendre les badges et boutons plus visibles,
- ou préparer une version du TP en deux variantes (démo + exercice guidé pour les étudiants) — dites-moi ce que vous préférez.
