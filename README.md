# Ticket Events API

API para gestionar eventos y categorías, integrando AWS SQS, S3 y DynamoDB mediante LocalStack.

## Tech Stack

### LocalStack
* DynamoDB: Almacenamiento de eventos y categorías
* SQS (Simple Queue Service): Procesamiento asíncrono para notificaciones

### Golang
* Gin: Framework web
* AWS SDK v2: Integración con servicios AWS

## Levantar LocalStack con Docker Compose

```bash
docker-compose up -d
```

Verificar que LocalStack está corriendo:

```bash
docker-compose logs localstack
```

## Crear recursos en LocalStack

Ejecutar el script aws-config.sh para crear la cola SQS y las tablas DynamoDB necesarias:

```bash
bash aws-config.sh
```

## Levantar la API de Go

Instalar dependencias:

```bash
go mod tidy
```

Cargar la data de prueba
```bash
go run scripts/fake-data.go
```

Ejecutar la API:

```bash
go run cmd/main.go
```

## Verificar en LocalStack

### Ver mensajes en SQS:
```bash
aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/event-queue
```



### Ver eventos en DynamoDB:
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name events
```

### Ver categorías en DynamoDB:
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name categories
```

## Documentación

- **Swagger**: Disponible en `docs/swagger.yaml`
- **Postman Collection**: Disponible en `docs/ticket-events.postman_collection.json`

## Estructura del Proyecto

```
ticket-events/
├── cmd/
│   └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── awsconfig/           # Configuración de AWS
│   ├── db/                  # Cliente de DynamoDB
│   ├── handler/             # Handlers HTTP
│   ├── model/               # Modelos de datos
│   ├── queue/               # Cliente de SQS
│   └── service/             # Servicios de negocio
``` 