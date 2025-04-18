# **`Golang + Gin` REST API template**
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Firebase](https://img.shields.io/badge/firebase-a08021?style=for-the-badge&logo=firebase&logoColor=ffcd34)

### Small **REST API** template written in `Golang` using `Gin`, `PostgreSQL`, `JWT` and `Redis`. Project: user profile *CRUD* with authorization, authentication and token validation, `Swagger` generator and metrics gainer 
---
### Endpoints:
1. `POST:/api/v1/users` - registrates new user
2. `POST:/api/v1/auth` - authorizes user, generates JWT token and writing it in cache
3. `POST:/api/v1/auth/changePassword` - changes password with deleting valid token in `Redis` valid tokens storage *(requires token in Bearer header)*
4. `GET:/api/v1/users/:username` - get user profile by username *(requires token in Bearer header)*
5. `GET:/api/v1/users/me` - returns current user profile *(requires token in Bearer header)*
6. `PATCH:/api/v1/users/me` - updates current user profile *(requires token in Bearer header)*

+ Visit http://localhost:8000/swagger/index.html to see `Swagger` specification
+ Use `swag init -g ./cmd/main.go -o docs` to auto-update `openapi.yaml` 
---
### Launch guide:
+ To launch webapp run:
```Shell
docker compose up --build
```
