FROM golang
WORKDIR /JWT
COPY ./main.go /JWT/main.go
COPY package.sh /JWT/package.sh
RUN mkdir /database
COPY ./database/JWT_database.db /database/JWT_database.db
RUN apt-get update
RUN apt-get install sqlite3
RUN /JWT/package.sh
CMD ["go run", "/JWT/main.go"]