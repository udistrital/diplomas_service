# diplomas_mid

API MID para la gestion de diplomas digitales.

## Flujo implementado

`POST /v1/documento_digital/:id/firmar`

Orquesta el primer flujo de firma de diploma:

1. Valida el rol activo del firmante en `administrativa_amazon_api`.
2. Envia el PDF original a `firma_electronica_mid` por `POST /api/v2/firma_electronica`.
3. Recibe `file` firmado y `uuid_documento`.
4. Guarda el PDF firmado en S3 con la ruta `diplomas/{uuid_documento}/diploma.pdf`.
5. Actualiza `diplomas_crud` con `PUT /v1/documento_digital/:id/uuid`.

`GET /v1/firmantes/rol-activo?documento_identidad=10000000`

Valida si un documento pertenece a un secretario academico, decano, secretaria general o rector activo.

`POST /v1/firmantes/firma`

Valida el rol activo del firmante, sube la imagen de firma a `DOCUMENTOS_CRUD` y registra el enlace en `diplomas_crud.firma_firmante`.
La imagen recibida se normaliza antes de almacenar: se quita fondo claro, se convierte el trazo a negro, se recortan bordes vacios y se guarda como PNG transparente.

## Variables

```shell
DIPLOMAS_MID_HTTP_PORT=8090
FIRMA_ELECTRONICA_MID_URL=http://localhost:8080/api
DIPLOMAS_CRUD_URL=http://localhost:8081
ADMINISTRATIVA_AMAZON_API_URL=${ADMINISTRATIVA_AMAZON_API_URL}
DOCUMENTOS_CRUD_GESTOR_DOCUMENTAL_URL=${DOCUMENTOS_CRUD_GESTOR_DOCUMENTAL_URL}
AWS_REGION=us-east-1
AWS_ENDPOINT_URL=http://localhost:4566
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
DIPLOMAS_S3_BUCKET=diplomas-digital-local
```

Si corre en Docker y consume servicios del host, usar `host.docker.internal` en las URLs.

## Ejecucion

```shell
go mod tidy
go run .
```

Docker:

```shell
docker compose up -d --build
```
