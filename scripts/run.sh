#!/bin/bash

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Конфигурация
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="project-sem-1"
DB_USER="validator"
DB_PASSWORD="val1dat0r"
APP=prices_backend

run_app(){
  echo "Попытка запуска приложения"
  cd ..
  ./${APP} -d ${DB_HOST}:${DB_PORT}/${DB_NAME} -p ${DB_PASSWORD} -u ${DB_USER}
}

run_app