#!/bin/bash

# Obtener el token Dojo de la variable de entorno
# Ejecutar el comando para 
dagger call test-all --source=.
# Ejecutar el comando Dagger para generar el reporte Trivy
dagger call vulnerability --source=. export --path=./reports/trivy-report-test.json -vv

# Ejecutar el comando Go para importar el reporte a Defect Dojo
MSYS_NO_PATHCONV=1 docker run --network=host --mount type=bind,source="$(pwd)"/reports,target=/reports filstimatias/dojoapi:48 http://192.168.49.2:30080/api/v2/import-scan/ 9a8105ed28e6f29d2af9e1d56706a8da010643e7 /reports/defectdojo-api-trivy.config /reports/trivy-report-test.json

# Eliminar el archivo de reporte si existe
if [ -f "reports/trivy-report-test.json" ]; then
  rm reports/trivy-report-test.json
fi


# Guia de comandos
# docker run -d -it --network=host ubuntu bash
# apt-get update && apt-get install -y inetutils-ping && apt-get install -y curl
# kubectl -n test port-forward  svc/my-app01-service 8080:8080

# minikube start --listen-address=0.0.0.0

# kubectl -n defectdojo port-forward svc/defectdojo-django-np 8080:80
