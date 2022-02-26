FROM postgres
ENV POSTGRES_PASSWORD docker
ENV POSTGRES_DB docker
COPY ./migrations/*.sql /docker-entrypoint-initdb.d/