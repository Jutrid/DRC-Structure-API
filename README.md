# DRC Structure API 🇨🇩

API REST en **Go (Gin)** + **SQLite (GORM)** exposant les données administratives de la République Démocratique du Congo : **Provinces**, **Villes** et **Territoires**.

Les données sont seedées automatiquement depuis `Data_in_json/` au premier démarrage.

## Stack

- Go 1.21+
- [Gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io) + `gorm.io/driver/sqlite`
- SQLite

## Structure

```
.
├── main.go
├── go.mod
├── database/
│   ├── database.go   # init, migrate, seed
│   └── db.sqlite
├── models/
│   ├── province.go
│   ├── city.go
│   └── territory.go
├── handlers/
│   ├── common.go
│   ├── province.go
│   ├── city.go
│   ├── territory.go
│   └── stats.go
├── routes/
│   └── router.go
└── Data_in_json/
    ├── provinces.json   (26 provinces)
    ├── cities.json      (102 villes)
    └── territories.json (164 territoires)
```

### Modèles

**Province** `provinces`
- `id`, `name`, `code`, `principal_town`, `surface`, `population`, `latitude`, `longitude`

**City** `cities`
- `id`, `province_id` (FK), `name`, `code` + relation `Province`

**Territory** `territories`
- `id`, `province_id` (FK), `name`, `code` + relation `Province`

## Installation

```bash
git clone <repo>
cd DRC_Structure_API

go mod tidy

# Variables d'environnement optionnelles
# DB_PATH=database/db.sqlite PORT=8080

go run main.go
```

> Au premier lancement, la DB est migrée et seedée depuis `Data_in_json/*.json`. Les lancements suivants sautent le seed si la table `provinces` n'est pas vide.

## Endpoints

Base URL: `http://localhost:8080`

| Méthode | Route | Description |
|---------|-------|-------------|
| GET | `/` | Welcome + liste des endpoints |
| GET | `/health` | Healthcheck |
| GET | `/api/v1/stats` | Compteurs provinces/villes/territoires |
| GET | `/api/v1/provinces?search=&page=1&limit=50` | Liste provinces |
| GET | `/api/v1/provinces/:id` | Détail province (+ counts) |
| GET | `/api/v1/provinces/:id/cities?search=&page=&limit=` | Villes d'une province |
| GET | `/api/v1/provinces/:id/territories?search=&page=&limit=` | Territoires d'une province |
| GET | `/api/v1/cities?province_id=&search=&page=&limit=` | Liste villes (filtre province_id) |
| GET | `/api/v1/cities/:id` | Détail ville (+ province) |
| GET | `/api/v1/territories?province_id=&search=&page=&limit=` | Liste territoires |
| GET | `/api/v1/territories/:id` | Détail territoire (+ province) |

### Pagination

Tous les endpoints de liste retournent :

```json
{
  "data": [...],
  "meta": {
    "page": 1,
    "limit": 50,
    "total": 26,
    "total_pages": 1
  }
}
```

- `page` défaut 1, `limit` défaut 50 (max 100)
- `search` filtre `LIKE %search%` sur `name` (et `principal_town` pour provinces)

### Exemples curl

```bash
# Toutes les provinces
curl http://localhost:8080/api/v1/provinces

# Recherche
curl "http://localhost:8080/api/v1/provinces?search=Kivu"

# Détail province 19 (Nord-Kivu)
curl http://localhost:8080/api/v1/provinces/19

# Villes du Haut-Katanga (id 3)
curl http://localhost:8080/api/v1/provinces/3/cities

# Toutes les villes, filtrées par province
curl "http://localhost:8080/api/v1/cities?province_id=19"

# Territoires avec pagination
curl "http://localhost:8080/api/v1/territories?page=2&limit=10"

# Stats
curl http://localhost:8080/api/v1/stats
```

## Build

```bash
go build -o drc_api .
./drc_api
```

## Variables d'environnement

| Var | Défaut | Description |
|-----|--------|-------------|
| `PORT` | `8080` | Port HTTP |
| `DB_PATH` | `database/db.sqlite` | Chemin SQLite |

## TODO (extensions possibles)

- POST/PUT/DELETE avec auth
- Filtre par code (`CD...`)
- Export GeoJSON
