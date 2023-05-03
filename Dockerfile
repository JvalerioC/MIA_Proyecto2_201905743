# Utilizar una imagen base de Go
FROM golang:1.20

# Establecer el directorio de trabajo en /app
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copiar los archivos de la aplicación al contenedor
COPY . .

# Compilar la aplicación
RUN go build -o app .
RUN apt-get update && apt-get install -y graphviz

#se expone el puerto
EXPOSE 3000

# Establecer el comando predeterminado que se ejecutará cuando se inicie el contenedor
CMD ["./app"]