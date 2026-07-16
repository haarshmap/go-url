# go-url
a slightly complex url shortener to learn backend with `go` using `echo`, `sqlc`,`sqlite` and `templ` while also learning the basics of `docker` and `redis`

1. to build for prod

```
docker compose up --build
```


2. to build for dev
```
docker compose -f docker-compose.dev.yml up --build
```

3. schema definitions

```
users:
    - id `INTEGER PRIMARY KEY`
    - username `TEXT UNIQUE NOT NULL`
    - email `TEXT UNIQUE NOT NULL`
    - hash_password `TEXT NOT NULL`

links:
    - id `INTEGER PRIMARY KEY`
    - short_id `TEXT UNIQUE NOT NULL`
    - orig_url `TEXT NOT NULL`
    - expiry `DATETIME NOT NULL`
    - user_id `INTEGER NOT NULL REFERENCES accounts(id)`


linkdata: 
    - id `INTEGER PRIMARY KEY`
    - link_id `INTEGER NOT NULL REFERENCES links(id)`
    - access_time `DATETIME DEFAULT CURRENT_TIMESTAMP`
    - country `TEXT`
    - ip_address `TEXT`
    - user_agent `TEXT`
```