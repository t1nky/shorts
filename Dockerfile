FROM library/postgres
COPY database/init.sql /docker-entrypoint-initdb.d/
ENV POSTGRES_USER shorts_user
ENV POSTGRES_PASSWORD docker