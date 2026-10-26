# diplomas_service

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

Docker (preparar primero `main` y `runtime/` según sección CI/CD):

```shell
docker compose up -d --build
```

## CI/CD con Drone

Configuración basada en [lineamientos OAS](https://github.com/udistrital/lineamientos_oas/tree/master/generacion_de_apis/cicd).
El módulo Go y las variables `DIPLOMAS_MID_*` conservan sus nombres por compatibilidad.

El pipeline `.drone.yml` verifica formato, ejecuta `go vet`, pruebas con cobertura y compila
un binario Linux amd64 estático. Sonar analiza código y cobertura en eventos push;
las solicitudes de cambio ejecutan validación y compilación sin publicar ni usar secretos.
Se admiten `feat/*`, `feature/*`, `develop`, `release/*`, `hotfix/*` y `master`.

| Rama (push) | Imagen ECR | Servicio ECS | Clúster |
| --- | --- | --- | --- |
| `release/*` | `diplomas_service:release` y SHA corto | `diplomas_service_test` | `test` |
| `master` | `diplomas_service:latest` y SHA corto | `diplomas_service_prod` | `oas` |

El nombre real se obtiene de `DRONE_REPO_NAME`. Los pasos son secuenciales: un fallo
impide publicar o desplegar. ECS fuerza actualización y espera estabilidad del servicio.
`Dockerfile` empaqueta el binario compilado por Drone sobre distroless Debian 12,
siguiendo la estructura de [syllabus_mid](https://github.com/udistrital/syllabus_mid).
El paso build instala Chromium en Debian 12 y ejecuta `scripts/prepare-chromium-runtime.sh`
para copiar navegador, bibliotecas y configuración de fuentes a `runtime/`.
La imagen incluye estos archivos y `assets/` para generar diplomas sin instalar paquetes
en la imagen final ni requerir shell.

### Configuración externa requerida

1. Activar repositorio en Drone con ejecutor Docker Linux amd64.
2. Registrar secretos `SONAR_HOST`, `SONAR_TOKEN`, `AWS_ACCESS_KEY_ID`,
   `AWS_SECRET_ACCESS_KEY` y `AWS_REGION`. No habilitarlos para PR externos.
3. Crear/configurar proyecto SonarQube `diplomas_service`.
4. Proveer permisos AWS para publicar/crear repositorio ECR, actualizar servicios ECS
   y consultar su estado. Repositorio ECR usa nombre del repositorio GitHub.
5. Crear previamente clústeres, servicios y definiciones de tarea ECS indicados arriba.
   La tarea test debe consumir etiqueta `release`; producción, `latest`.
   Configurar arquitectura X86_64, puerto y health check adecuados al servicio.
6. Inyectar variables de `conf/app.conf` mediante definición de tarea o secretos ECS.
   Definir `DIPLOMAS_MID_HTTP_PORT=8090`, `DIPLOMAS_MID_RUNMODE=prod`, URLs de APIs,
   `AWS_REGION` y `DIPLOMAS_S3_BUCKET`. `AWS_ENDPOINT_URL` se reserva para emuladores
   o endpoints alternativos. En ECS preferir rol de tarea con acceso S3 a claves estáticas.
   Los secretos de publicación Drone no se transfieren automáticamente al contenedor.

Este pipeline actualiza servicios existentes; no crea infraestructura ni registra nuevas
versiones de definición de tarea. Cambios de variables, recursos o puerto requieren
actualizar dicha definición en AWS. La estabilidad ECS depende del health check configurado.

### Verificación local

```shell
go vet ./...
go test -coverprofile=coverage.out ./...
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .
# Ejecutar preparación en Debian 12 con Chromium instalado:
sh scripts/prepare-chromium-runtime.sh
docker build -t diplomas_service:ci .
```

Las pruebas que requieren `PDFA_TEST_INPUT` se omiten salvo configuración explícita.
Antes de `docker compose up -d --build`, preparar `main` y `runtime/` como en Drone.
La compilación y preparación deben usar Linux amd64 y Debian 12 para mantener compatibilidad
con las bibliotecas de la imagen final. `runtime/` contiene artefactos generados y no se versiona.
