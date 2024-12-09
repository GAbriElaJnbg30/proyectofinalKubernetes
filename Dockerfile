FROM golang:1.18

# Instalar el cliente MySQL en el contenedor
RUN apt-get update && apt-get install -y default-mysql-client

# Establecer el directorio de trabajo dentro del contenedor
WORKDIR /api-rest

# Copiar los archivos de Go y los archivos necesarios
COPY go.mod ./ 
COPY go.sum ./ 
COPY .env ./ 
RUN go mod download

COPY *.go ./ 
COPY . ./ 

# Construir la aplicación
RUN go build -o /api

# Exponer el puerto de la aplicación
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["/api"]
