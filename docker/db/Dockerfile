FROM postgres:16-bullseye

# Копируем начальные данные (если есть)
COPY ./init /docker-entrypoint-initdb.d/
