version: "3"
services:
  {{Name}}_postgres:
    container_name: {{Name}}_db
    hostname: {{Name}}_db
    image: "postgres:11"
    env_file: .env
    ports:
      - "5432:5432"
# UNCOMMENT ONCE YOU HAVE MIGRATIONS
#  {{Name}}_migrations:
#    container_name: migrations
#    image: migrate/migrate:v4.6.2
#    command: ["-path", "/migrations/", "-database", "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@{{Name}}_postgres:5432/${POSTGRES_DB}?sslmode=disable", "up"]
#    depends_on:
#      - {{Name}}_postgres
#    env_file: .env
#    restart: on-failure
#    links: 
#      - {{Name}}_postgres
#    volumes:
#      - ./migrations:/migrations 
#
  {{Name}}:
    container_name: {{Name}}
    build:
      context: .
      dockerfile: Dockerfile
    env_file: .env
    volumes:
      - ./:/go/src/{{Name}}
    ports:
      - "8080:8080"
    links:
      - {{Name}}_postgres

  {{Name}}_test:
    container_name: {{Name}}_test
    build:
      context: .
      dockerfile: test.Dockerfile
    env_file: .env
    volumes:
      - ./:/go/src/{{Name}}
    ports:
      - "9999:9999"
    links:
      - {{Name}}_postgres

