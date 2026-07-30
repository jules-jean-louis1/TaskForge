# Diagramme de Composants (Logique)
Ce diagramme montre comment les briques logicielles communiquent entre elles, indépendamment de l'infrastructure.

```mermaid
graph TB
    subgraph "Frontend Layer (Angular)"
        UI[User Interface]
        Auth[Auth Service]
        ApiClient[API Client]
    end

    subgraph "Gateway Layer (Traefik)"
        Router[Traefik Router]
        SSL[TLS Termination]
        MW[Middlewares: Compression/Limit]
    end

    subgraph "Backend Layer (Go)"
        API[Go API / Gin / Air]
        DB_Logic[Database Handler]
        Middleware[Auth/Logging Middleware]
    end

    subgraph "Data Layer"
        Postgres[(PostgreSQL 17)]
    end

    %% Interactions
    UI --> Auth
    UI --> ApiClient
    ApiClient -- HTTP/JSON --> Router
    Router --> SSL --> MW
    MW --> API
    API --> DB_Logic
    DB_Logic -- SQL --> Postgres
```

# Diagramme de Déploiement (Infrastructure Docker)
Ce diagramme montre comment les conteneurs sont isolés dans leurs réseaux respectifs et comment le trafic circule à travers les ports.

```mermaid
graph TB
    subgraph Internet
        User((Utilisateur))
    end

    subgraph "Docker Host"
        subgraph "Frontend Network (bridge)"
            Traefik[[Conteneur Traefik]]
            Frontend[[Conteneur Frontend]]
        end

        subgraph "Backend Network (bridge)"
            Backend[[Conteneur Backend]]
            Database[(Conteneur Postgres)]
        end

        %% Ports
        User -- "80/443" --> Traefik
        
        %% Routing Traefik
        Traefik -- "app.localhost" --> Frontend
        Traefik -- "api.localhost" --> Backend

        %% Inter-service communication
        Backend -- "Port 5432" --> Database
    end

    %% Volumes
    subgraph "Volumes"
        DB_Data[(postgres_data)]
        Certs[(Certs/Config)]
    end

    Database -.-> DB_Data
    Traefik -.-> Certs
```

### Isolation Réseau :

La Database est totalement isolée du monde extérieur. Seul le Backend peut lui parler via le backend_network.
Traefik sert de point d'entrée unique (Reverse Proxy) pour le Frontend et le Backend.

### Flux de Données :

Le trafic entrant arrive sur Traefik (Port 80/443).
Traefik analyse le Host (ex: api.localhost) et redirige vers le bon conteneur.
Note : Ton Backend appartient aux deux réseaux car il doit être joignable par Traefik (pour l'API) et doit pouvoir joindre la Database.

### Persistance :

Les données de la base sont stockées dans un volume nommé pour survivre au redémarrage des conteneurs.