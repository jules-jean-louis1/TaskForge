.PHONY: env certs trust-certs auth_traefik init

env: # Crée le fichier .env à partir de l'exemple s'il n'existe pas déjà
	@if [ ! -f "docker/.env" ]; then \
			cp docker/.env.example docker/.env; \
			echo "Created docker/.env from example"; \
	else \
			echo "docker/.env already exists"; \
	fi

certs: # Génère les certificats SSL pour les domaines locaux en utilisant mkcert dans un conteneur Docker
	@docker run --rm -v "$(CURDIR)/traefik:/workspace" -w /workspace debian:bookworm-slim sh -c 'apt-get update >/dev/null && apt-get install -y --no-install-recommends mkcert ca-certificates >/dev/null && mkdir -p certs && mkcert -cert-file certs/local.crt -key-file certs/local.key localhost "*.localhost" app.localhost api.localhost db.localhost mail.localhost traefik.localhost monitor.localhost prometheus.localhost 127.0.0.1 ::1' && echo "Certificats générés"

trust-certs: # Installe l'autorité de certification locale de mkcert dans le trust store du système/navigateur
	@command -v mkcert >/dev/null 2>&1 || (echo "mkcert n'est pas installé localement. Installe-le puis relance 'make trust-certs'." && exit 1)
	@mkcert -install && echo "Autorité locale mkcert installée dans le trust store système/navigateur"

auth_traefik: # Génère un fichier .htpasswd pour l'authentification de base de Traefik en utilisant htpasswd dans un conteneur Docker
	@docker run --rm -v "$(CURDIR)/traefik":/data -w /data httpd:alpine htpasswd -b -B -c .htpasswd admin password

perms: # Donne les droits d'exécution au script docker helper
	@chmod +x ./docker/docker.sh

init: env perms certs trust-certs auth_traefik # Exécute toutes les étapes d'initialisation pour préparer l'environnement de développement

# --- Lancement des services Docker ---

# --- DEVELOPMENT ---	

launch-dev: # Lance les services Docker en mode développement avec docker-compose
	./docker/docker.sh up -d --build

remove: # Arrête les conteneurs et supprime les volumes associés (remise à zéro)
	./docker/docker.sh down -v

logs-all: # Affiche les logs de tous les conteneurs en continu
	./docker/docker.sh logs -f

logs: # Affiche les logs d'un conteneur spécifique ou général (ex: make logs s=api)
	./docker/docker.sh logs -f $(s)

ps: # Liste les conteneurs actifs et leur statut
	./docker/docker.sh ps

start: # Démarre les conteneurs existants (ex: make start ou make start s=api)
	./docker/docker.sh start $(s)

stop: # Arrête temporairement les conteneurs (ex: make stop ou make stop s=api)
	./docker/docker.sh stop $(s)

restart: # Redémarre un ou tous les conteneurs proprement (ex: make restart s=api)
	./docker/docker.sh restart $(s)

rebuild: # Force le rebuild et la relance complète d'un service (ex: make rebuild s=api)
	./docker/docker.sh stop $(s) || true
	./docker/docker.sh rm -f $(s) || true
	./docker/docker.sh up -d --build $(s)

remove-service: # Arrête un service, supprime son conteneur et SES VOLUMES associés (ex: make remove-service s=postgres-test)
	@if [ -z "$(s)" ]; then \
		echo "Erreur: Spécifiez un service (ex: make remove-service s=postgres-test)"; \
	else \
		./docker/docker.sh rm -f -s -v $(s); \
	fi

exec: # Ouvre un shell interactif dans un conteneur (ex: make exec s=api)
	@if [ -z "$(s)" ]; then \
		echo "Erreur: Vous devez spécifier un service (ex: make exec s=api)"; \
	else \
		./docker/docker.sh exec -it $(s) sh || ./docker/docker.sh exec -it $(s) /bin/bash; \
	fi

stats: # Affiche la consommation CPU/RAM des conteneurs en temps réel
	docker stats

clean: # Nettoie Docker en supprimant les images, conteneurs et réseaux inutilisés (Prune)
	docker system prune -f
	docker volume prune -f

# --- BACKEND ---

lint-back:
	@cd backend && golangci-lint run ./...

launch-prod: # Lance les services Docker en mode production avec docker-compose
	./docker/docker.sh -f docker/docker-compose.prod.yml up -d --build

format-back: # Formate le code source du backend avec gofumpt
	@cd backend && gofumpt -l -w .

test-all: # Exécute les tests unitaires du backend avec go test
	@cd backend && go test -v ./tests/...

# --- Aide ---

help: # Affiche la liste et la description de toutes les commandes disponibles
	@echo "Commandes disponibles :"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?# .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[1;32m%-20s\033[0m %s\n", $$1, $$2}'