#!/bin/bash
set -euo pipefail

AWS_ENDPOINT="--endpoint-url=http://localhost:4566"

echo "🚀 Configurando recursos AWS en LocalStack..."

# Verificar que LocalStack esté ejecutándose
echo "📦 Verificando LocalStack..."
if ! docker ps | grep -q localstack; then
    echo "❌ LocalStack no está ejecutándose. Iniciando..."
    docker-compose up -d
    echo "⏳ Esperando que LocalStack esté listo..."
    sleep 10
else
    echo "✅ LocalStack está ejecutándose"
fi

# Crear tabla DynamoDB de eventos solo si no existe
echo "🗄️ Configurando tabla DynamoDB de eventos..."
table_exists=$(aws $AWS_ENDPOINT dynamodb list-tables 2>/dev/null | grep 'events' || true)
if [ -z "$table_exists" ]; then
  echo "📝 Creando tabla DynamoDB 'events'..."
  aws $AWS_ENDPOINT dynamodb create-table \
    --table-name events \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
  echo "✅ Tabla DynamoDB 'events' creada exitosamente"
else
  echo "✅ La tabla DynamoDB 'events' ya existe."
fi

# Crear tabla DynamoDB de categorías solo si no existe
echo "🗄️ Configurando tabla DynamoDB de categorías..."
table_exists=$(aws $AWS_ENDPOINT dynamodb list-tables 2>/dev/null | grep 'categories' || true)
if [ -z "$table_exists" ]; then
  echo "📝 Creando tabla DynamoDB 'categories'..."
  aws $AWS_ENDPOINT dynamodb create-table \
    --table-name categories \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
  echo "✅ Tabla DynamoDB 'categories' creada exitosamente"
else
  echo "✅ La tabla DynamoDB 'categories' ya existe."
fi



# Crear cola SQS solo si no existe
echo "📬 Configurando cola SQS..."
queue_exists=$(aws $AWS_ENDPOINT sqs list-queues 2>/dev/null | grep 'event-queue' || true)
if [ -z "$queue_exists" ]; then
  echo "📝 Creando cola SQS 'event-queue'..."
  aws $AWS_ENDPOINT sqs create-queue --queue-name event-queue
  echo "✅ Cola SQS 'event-queue' creada exitosamente"
else
  echo "✅ La cola SQS 'event-queue' ya existe."
fi

# Verificar configuración
echo "🔍 Verificando configuración..."
echo "📊 Tablas DynamoDB:"
aws $AWS_ENDPOINT dynamodb list-tables 2>/dev/null || echo "❌ Error listando tablas DynamoDB"

echo "📊 Colas SQS:"
aws $AWS_ENDPOINT sqs list-queues 2>/dev/null || echo "❌ Error listando colas SQS"

echo "🎉 Configuración completada!"
echo "💡 Ahora puedes ejecutar: go run cmd/main.go" 