FROM golang
WORKDIR /app
#COPY ./main.go /main.go
COPY app/package.sh /app/package.sh
#RUN mkdir /database
#COPY ./database/JWT_database.db /database/JWT_database.db
#RUN mkdir /keys
#COPY ./keys/id_rsa /keys/id_rsa
RUN apt-get update
RUN apt-get install sqlite3
RUN sh /app/package.sh
#CMD [ "go", "run", "/main.go" ]