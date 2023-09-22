#!/bin/bash

sudo systemctl start postgresql

DB_NAME="Page"
DB_USER="aabiji"
DB_PASSWORD="0000"

sudo -u postgres psql -c "CREATE DATABASE $DB_NAME;"
sudo -u postgres psql -c "CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';"
sudo -u postgres psql -c "ALTER ROLE $DB_USER SET client_encoding TO 'utf8';"
sudo -u postgres psql -c "ALTER ROLE $DB_USER SET default_transaction_isolation TO 'read committed';"
sudo -u postgres psql -c "ALTER ROLE $DB_USER SET timezone TO 'UTC';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON SCHEMA public TO $DB_USER;"

echo "Created Postgresql database $DB_NAME for $DB_USER with password $DB_PASSWORD."
echo "Run 'systemctl status postgresql' to see status of the postgresql database service."
echo "Run 'systemctl stop postgresql' to stop the postgresql database service."
echo "Run 'psql -U $DB_USER -d $DB_NAME' to interact with postgresql database'."
