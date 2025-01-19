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

create_table() {
    echo -e "\nПроверка PostgreSQL"

    # Базовая проверка подключения
    if ! PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c '\q' 2>/dev/null; then
        echo -e "${RED}✗ PostgreSQL недоступен${NC}"
        return 1
    fi


    echo "Выполняем создание таблиц"
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "
        CREATE TABLE public.prices (
             id             integer     NOT NULL,
             name           text        NULL,
             category       text        NULL,
             price          numeric     NULL,
             create_date    date        NULL,
             CONSTRAINT item_pk PRIMARY KEY (id)
        );" 2>/dev/null; then
        echo -e "${GREEN}✓ PostgreSQL таблица создана${NC}"
        return 0
    else
        echo -e "${RED}✗ Ошибка выполнения запроса${NC}"
        return 1
    fi

    return 1
}

build_app(){
  echo "Запуск сборки приложения"
  pwd
  make build BIN=${APP}

  if [ $? -ne 0 ]; then
    echo -e "${RED}✗ Ошибка при выполнении сборки${NC}"
    return 1
  fi

  echo -e "${GREEN}✓ Сборка завершена успешно${NC}"
}

create_table
build_app