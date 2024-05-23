## Requisitos

- Docker
- Docker Compose

## Despliegue

El proyecto se puede levantar ejecutando el siguiente comando en la raíz del proyecto:

```bash
docker-compose up --build -d --scale app=N
```

Donde `N` es el número de instancias que se levantan de la api.
