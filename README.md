[![Go](https://github.com/Floflow-V/go-rest-api-chi-example/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Floflow-V/go-rest-api-chi-example/actions/workflows/ci.yml)

# go-rest-api-chi-example

## Présentation générale

Ce projet est une **API REST** en Go utilisant le framework [Chi](https://github.com/go-chi/chi). Il adopte une **architecture monolithique modulaire** où chaque domaine métier est isolé pour faciliter la maintenance et l'évolutivité.

---

## Table des matières

- [go-rest-api-chi-example](#go-rest-api-chi-example)
  - [Présentation générale](#présentation-générale)
  - [Table des matières](#table-des-matières)
  - [Architecture du projet](#architecture-du-projet)
  - [Technologies utilisées](#technologies-utilisées)
  - [Gestion de la base de données avec SQLC](#gestion-de-la-base-de-données-avec-sqlc)
  - [Gestion de la configuration](#gestion-de-la-configuration)
  - [Fichiers d'environnement](#fichiers-denvironnement)
  - [Utilisation de Docker](#utilisation-de-docker)
  - [Automatisation avec Task](#automatisation-avec-task)
  - [Tests et qualité](#tests-et-qualité)
  - [Intégration Continue (CI)](#intégration-continue-ci)
  - [Documentation API (Swagger \& Scalar)](#documentation-api-swagger--scalar)
  - [Collections Bruno](#collections-bruno)
  - [Démarrage rapide](#démarrage-rapide)

---

## Architecture du projet

L'architecture est **monolithique modulaire** :

- Chaque domaine (ex : `author`, `book`) possède son propre dossier dans `internal/`, avec ses handlers, services, repositories, DTO, etc.
- Les dépendances sont injectées explicitement, facilitant les tests et la maintenance.
- La configuration, la base de données, le logger, la validation et les mocks sont séparés dans des modules dédiés.
- Les routes sont centralisées dans `internal/app`.
- Le point d'entrée se trouve dans `cmd/go-rest-api-chi-example/main.go`.

```
internal/
	author/
	book/
	app/
	config/
	database/
	logger/
	mocks/
	model/
	response/
	validator/
```

---

## Technologies utilisées

- **Go** (>=1.25)
- **Chi** : router HTTP léger et performant
- **SQLC** : génération automatique de code SQL
- **Zerolog** : logging structuré et performant
- **Swaggo** : génération automatique de documentation Swagger
- **Scalar** : UI moderne pour la doc Swagger
- **Bruno** : gestionnaire de collections de requêtes API
- **Task** : automatisation des tâches de développement
- **Docker & Docker Compose** : conteneurisation
- **Testify, GoMock** : tests unitaires et mocks

---

## Gestion de la base de données avec SQLC

[SQLC](https://sqlc.dev/) génère automatiquement du **code Go type-safe à partir de requêtes SQL**. Vous écrivez du SQL pur, SQLC produit les fonctions Go correspondantes.

- Les requêtes SQL sont écrites dans `internal/database/queries/` avec des annotations (`:exec`, `:one`, `:many`)
- Le code Go est généré dans `internal/database/sqlc/` (structs, fonctions, interface `Querier`)
- L'interface `Querier` est injectée dans les services et utilisée pour les mocks de tests
- Les fichiers générés ne doivent **jamais** être modifiés manuellement

**Régénérer après modification des fichiers SQL :**

```sh
sqlc generate
```

---

## Gestion de la configuration

La configuration est centralisée dans [`internal/config`](internal/config/) et chargée à partir de fichiers d'environnement (`.env`) et de variables système. Les paramètres de connexion à la base de données, au serveur, au logger, etc., sont tous configurables.

- **.env.example** : modèle à copier pour créer votre propre `.env`
- **.env** : contient les variables locales (ignoré par git)
- Les valeurs sont lues automatiquement au démarrage

**Variables principales :**

```
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=secret
DB_NAME=go_boilerplate_rest_api_chi
API_PORT=8080
LOG_LEVEL=debug
```

---

## Fichiers d'environnement

- **.env.example** : référence de toutes les variables nécessaires
- **.env** : à personnaliser selon votre environnement local ou de production
- **Sécurité** : `.env` est ignoré par git pour éviter toute fuite de secrets

---

## Utilisation de Docker

Le projet est prêt à l'emploi avec Docker :

- **docker-compose.yml** orchestre l'API et la base de données
- **Dockerfile** construit une image Go optimisée
- Les volumes assurent la persistance des données
- Les variables d'environnement sont injectées automatiquement

**Commande utile :**

```sh
docker compose up --build
```

---

## Automatisation avec Task

Le projet utilise [Task](https://taskfile.dev) pour automatiser les tâches courantes :

- **Formatage** : `task format` (go fmt)
- **Lint** : `task lint` (golangci-lint)
- **Tests** : `task test` (unitaires), `task test-cover` (avec couverture), `task test-cover-details` (rapport HTML)
- **Génération de documentation** : `task doc` (Swagger)
- **Génération des mocks** : `task generate`
- **Build** : `task build` (exécutable statique)
- **Démarrage complet (API + DB)** : `task dev`

**Exemple de workflow développeur :**

```sh
cp .env.example .env
task dev
```

---

## Tests et qualité

- **Tests unitaires** : chaque module possède ses tests (`*_test.go`)
- **Mocks** : générés automatiquement (voir `internal/mocks/`)
- **Couverture** : mesurée et exportée

**Commandes utiles :**

```sh
task test
task test-cover
task test-cover-details
```

---

## Intégration Continue (CI)

Le projet utilise **GitHub Actions** (`.github/workflows/ci.yml`) pour :

- **Lint** : vérification du style et erreurs statiques
- **Tests** : exécution automatique à chaque push/PR
- **Build** : compilation de l'application
- **Govulncheck** : vérification des vulnérabilités
- **Docker Build** : vérification de la construction de l'image

---

## Documentation API (Swagger & Scalar)

- **Swagger** : la documentation OpenAPI est générée automatiquement à partir des annotations dans le code (voir `docs/`).
- **Scalar** : une UI moderne pour explorer et tester l'API, accessible sur `/api/docs` en local.
- **Mise à jour** : `task doc` régénère la documentation après modification des routes ou des schémas.

---

## Collections Bruno

Le dossier [`bruno-collection/`](bruno-collection/) contient des collections de requêtes prêtes à l'emploi pour [Bruno](https://www.usebruno.com/), un outil open-source pour tester et documenter les APIs :

- **CRUD complet** sur les entités (auteur, livre, etc.)
- **Healthcheck**
- **Organisation par dossier**
- **Environnements** (local, prod, etc.)

Importez la collection dans Bruno pour tester rapidement tous les endpoints de l'API.

---

## Démarrage rapide

1. Copier le fichier d'environnement :

   ```sh
   cp .env.example .env
   ```

2. Lancer l'environnement de dev (API + DB) :

   ```sh
   task dev
   ```

3. Accéder aux services :
   - **API** : [http://localhost:8080](http://localhost:8080)
   - **Documentation** : [http://localhost:8080/api/docs](http://localhost:8080/api/docs)
