# TP — Middlewares utiles & sous-groupe `admin` (≈ 40–45min)

Objectif pédagogique
--------------------
- Découvrir la notion de middleware dans Gin — global vs groupe vs route
- Montrer comment propager des informations via le contexte Gin (c.Set / c.Get)
- Exemples concrets : logging enrichi, authentification simple, rate-limiting et "fatigue" serveur

Structure ajoutée
-----------------
Fichiers modifiés/ajoutés (pokedex):

- `main.go`
  - Enregistrement des middlewares globaux : `EnrichedLogger` et `FatigueMiddleware`.
- `pokemon/middleware.go` (nouveau)
  - `EnrichedLogger()` : récupère `X-Trainer` et, si présent, stocke `trainer` dans le contexte; attache aussi le `target_pokemon` (si `:id` présent) au contexte.
  - `SimpleAuth(adminSecret string)` : middleware minimal d'authentification (header `X-Admin-Token` ou cookie `admin_token`).
  - `RateLimitMiddleware(maxRequests int, window time.Duration)` : demo d'un rate-limiter global par route (très simple, pour la pédagogie).
  - `FatigueMiddleware(delay time.Duration)` : ajoute un délai artificiel si header `X-Server-Fatigue=true`.
- `pokemon/model.go`
  - Ajout du champ `Level` (valeur par défaut 1) — utilisé pour le `level-up`.
- `pokemon/repository.go`
  - Initialisation de `Level` au chargement et à la création.
  - Ajout de `LevelUp(id, levels)` pour modifier l'état du pokémon en mémoire.
- `pokemon/handler_api.go`
  - Ajout du sous-groupe `admin` : `POST /api/v1/admin/pokemons/:id/level-up`.
  - Le route-handler `levelUpPokemon` montre l'utilisation du contexte (lecture de `trainer`) et modifie l'état métier.
- `pokemon/repository_test.go` (nouveau)
  - Test minimal pour `LevelUp`.

Comment les middlewares se comportent (points didactiques)
------------------------------------------------------

- Middlewares globaux (attachés à `router.Use(...)`) s'exécutent pour TOUTES les routes — c'est utile pour un logger, une gestion commune.
  - Ici `EnrichedLogger` fixe `trainer` dans le contexte et ajoute du contexte (target_pokemon) quand `:id` est présent.
  - `FatigueMiddleware` illustre l'impact d'un délai applicatif; essayer `X-Server-Fatigue=true` pour simuler un serveur lent.

- Middlewares de groupe (ex : admin group) s'appliquent uniquement aux routes du groupe — très pratique pour des politiques d'accès.
  - `admin.Use(SimpleAuth("admin-secret"))` : authentification simple, démonstration d'une couche de sécurité par groupe.

- Middlewares de route (ex : rate limiting) s'appliquent à une route précise.
  - Dans l'exemple `RateLimitMiddleware(5, 10*time.Minute)` empêche trop d'appels à la route `level-up`.

Propagation d'informations dans le contexte
-----------------------------------------

Le middleware `EnrichedLogger` montre comment :

- lire un header (`X-Trainer`) et le stocker en contexte via `c.Set("trainer", value)`
- tenter de résoudre `:id` et placer une structure `target_pokemon` dans le contexte si trouvée

Le handler `levelUpPokemon` récupère ces valeurs via `c.Get("trainer")` et `c.Get("target_pokemon")` ; c'est la façon recommandée d'envoyer des informations d'une middleware vers un handler en Gin.

Exemples / commandes à essayer (terminal)
----------------------------------------

1) Lancer le serveur

```powershell
cd pokedex
go run .
```

2) Tester l'endpoint ping (vérifier que Global middleware tourne)

```powershell
curl -i http://localhost:8080/ping
```

3) Appeler la route `level-up` sans authentification (devrait échouer)

```powershell
curl -i -X POST http://localhost:8080/api/v1/admin/pokemons/47/level-up
```

4) Appeler la route `level-up` en tant qu'admin (valeur d'exemple `admin-secret`)

```powershell
curl -i -X POST http://localhost:8080/api/v1/admin/pokemons/47/level-up -H "X-Admin-Token: admin-secret"
```

5) Passer le param `levels` pour monter plusieurs niveaux

```powershell
curl -i -X POST "http://localhost:8080/api/v1/admin/pokemons/47/level-up?levels=3" -H "X-Admin-Token: admin-secret"
```

6) Simuler la fatigue serveur (middleware global) et mesurer l'impact

```powershell
curl -i -X POST "http://localhost:8080/api/v1/admin/pokemons/47/level-up" -H "X-Admin-Token: admin-secret" -H "X-Server-Fatigue: true"
```

7) Tester la propagation du `trainer` dans le log (middleware global) :

```powershell
curl -i -X POST "http://localhost:8080/api/v1/admin/pokemons/47/level-up" -H "X-Admin-Token: admin-secret" -H "X-Trainer: ash"
```

8) Déclencher le rate-limit sur `level-up` (exécuté 5 fois, vous devriez recevoir 429 ensuite):

```powershell
for ($i=0; $i -lt 6; $i++) { curl -i -X POST "http://localhost:8080/api/v1/admin/pokemons/47/level-up" -H "X-Admin-Token: admin-secret" }
```

Proposition de déroulé pédagogique (≈ 40 min)
------------------------------------------

Total approximé : 40–45 minutes, exercices et démo inclus.

1) (5 min) Rappel rapide : qu'est-ce qu'un middleware ? exemples concrets (auth, logging, rate-limit).

2) (10 min) Global middleware : examiner `main.go` + `EnrichedLogger`
   - Expliquer pourquoi logger global est utile
   - Démonstration : envoyer `X-Trainer` et voir les logs / handlers lire `trainer` depuis le contexte

3) (10 min) Groupe `admin` et middleware d'auth : lecture du code `RegisterAPIRoutes` et `SimpleAuth`
   - Discussion : comparaison cookie vs header, pourquoi pas garder des secrets dans le code
   - Exo : tester la route sans token puis avec token

4) (8 min) Middleware de route : `RateLimitMiddleware` pour `level-up`
   - Simuler les appels et montrer le 429

5) (5 min) Fatigue serveur — `FatigueMiddleware`
   - Montrer l'effet de `X-Server-Fatigue` et discussion performance

6) (restant) Questions / variations

Variantes / exercices supplémentaires (si vous avez du temps)
- Implémenter un rate-limiter per-client (IP ou header-based) au lieu d'un compteur global.
- Remplacer `SimpleAuth` par un JWT + middleware de validation.
- Propager et stocker logs avancés (correlation id / request id) dans le contexte.
