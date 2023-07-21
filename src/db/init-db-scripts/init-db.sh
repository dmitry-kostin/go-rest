#!/bin/bash
set -e

echo "init"
echo "$POSTGRES_INIT"

IFS=',' read -ra ADDR <<< "$POSTGRES_INIT"
for i in "${ADDR[@]}"; do
    IFS=':' read -ra INFO <<< "$i"
    USER="${INFO[0]}"
    PASSWORD="${INFO[1]}"
    DB="${INFO[2]}"

    echo "init $USER $PASSWORD $DB"

    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
      CREATE DATABASE $DB;
      CREATE USER $USER WITH ENCRYPTED PASSWORD '$PASSWORD';
      GRANT ALL PRIVILEGES ON DATABASE $DB TO $USER;
      \c $DB $POSTGRES_USER
      GRANT ALL ON SCHEMA public TO $USER;
EOSQL
done
