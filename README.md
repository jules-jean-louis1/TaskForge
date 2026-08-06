# TaskForge

TaskForge est un MVP de gestion de tickets d’incidents (helpdesk interne) développé avec une architecture web conteneurisée.

## Contenu du projet

- `frontend/` : application Angular
- `backend/` : API REST en Go avec Gin
- `db/` : schéma PostgreSQL et initialisation
- `docker/` : configuration Docker Compose et helper
- `docs/` : architecture, ADR et gestion de projet

## Objectif

Offrir un produit fonctionnel avec :
- création, consultation, mise à jour et fermeture de tickets
- assignation et réassignation des tickets
- filtres par statut, priorité et technicien
- dashboard métier avec statistiques
- authentification et rôles : Admin, Technicien, Standard
- health checks et métriques applicatives
- lancement complet via Docker Compose

## Lancement rapide

1. Copier l’exemple de configuration :

```bash
cp docker/.env.example docker/.env
```

2. Lancer la stack de développement :

```bash
make launch-dev
```

3. Arrêter et supprimer la stack :

```bash
make remove
```

4. Lancer la stack de production :

```bash
make launch-prod
```

> `docker/docker.sh` charge `docker/.env` et exécute Docker Compose avec le nom de projet défini.

## Accès

- Frontend : `http://localhost:4200`
- Backend API : `http://localhost:8080/api/v1`
- Health API : `http://localhost:8080/api/v1/health`
- Health frontend : `http://localhost:4200/healthz`
- Metrics : `http://localhost:8080/api/v1/metrics`

## Variables d’environnement

Le fichier d’exemple se trouve dans `docker/.env.example`.
Copiez-le en `docker/.env` et ajustez les valeurs si nécessaire.

## Commandes utiles

- `make help` : liste des commandes Make disponibles
- `make test-all` : exécute les tests unitaires backend
- `make format-back` : formate le backend avec gofumpt
- `make lint-back` : lance golangci-lint sur le backend
- `make logs-all` : visualise les logs de tous les conteneurs

## Architecture et justification

- Architecture backend : Go + Gin, séparation contrôleurs/services/repository
- Architecture frontend : Angular, composants, services et routage
- Déploiement : Docker Compose avec frontend, backend, base de données et Traefik

Consulter également :
- `docs/infra.md` : diagrammes de composants et déploiement
- `docs/ADR_001.md` : décision technologique et justification
- `docs/project-management.md` : board Kanban, sprint backlog, daily logs et burn-down

## Notes supplémentaires

- Le backend expose un endpoint `/api/v1/health` et `/api/v1/metrics`.
- La base PostgreSQL est initialisée via `db/init.sql`.
- `docker/docker.sh` est le helper utilisé pour exécuter `docker compose` depuis le dossier racine.
