#!/bin/bash

LOG_FILE="setup_log.txt"
PG_HOST="localhost"        # Set your PostgreSQL host here
PG_PORT="5432"             # Set your PostgreSQL port here
PG_DATABASE="distilleries_db"
PG_SUPERUSER="distilleries_of_scotland_user"
PG_PASSWORD="ultra_secret_password"

export PGHOST="$PG_HOST"
export PGPORT="$PG_PORT"
export PGDATABASE="$PG_DATABASE"
export PGSUPERUSER="$PG_SUPERUSER"
export PGPASSWORD="$PG_PASSWORD"

# Function to log messages
log_message() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" > "$LOG_FILE"
}

check_for_postgres() {
  # Check if PostgreSQL is installed
  if command -v psql &> /dev/null; then
    log_message "PostgreSQL is installed."
    return 0  # PostgreSQL is installed
  else
    log_message "PostgreSQL is not installed."
    return 1  # PostgreSQL is not installed
  fi
}

install_postgres() {
  # Install PostgreSQL
  log_message "Installing PostgreSQL..."
  brew install postgresql
  log_message "PostgreSQL installed successfully."
}

start_postgres_service() {
  # Start PostgreSQL service
  log_message "Starting PostgreSQL service..."
  brew services start postgresql
  sleep 5  # Give some time for PostgreSQL to start
  log_message "PostgreSQL service started."
}

create_database() {
  dbname="distilleries_db"

  # Check if the database already exists
  if psql -h "$PG_HOST" -lqt | cut -d \| -f 1 | grep -qw "$dbname"; then
    log_message "Database '$dbname' already exists."
  else
    # Create the database
    createdb -h "$PG_HOST" "$dbname"
    log_message "Database '$dbname' created successfully."
  fi
}

create_postgres_superuser() {
  # Check if a superuser already exists
  if psql -h "$PG_HOST" -p "$PG_PORT" -tAc "SELECT COUNT(*) FROM pg_user WHERE usename = '$PG_SUPERUSER' AND usecreatedb AND usesuper;" | grep -qw 0; then
    # Check if PostgreSQL is running
    if brew services list | grep -q "postgresql"; then
      # Create superuser and set it as a member of 'distilleries_db'
      sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "CREATE USER $PG_SUPERUSER WITH PASSWORD '$PG_PASSWORD' SUPERUSER;"
      sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "ALTER USER $PG_SUPERUSER WITH CREATEDB;"
      sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "GRANT ALL PRIVILEGES ON DATABASE $PG_DATABASE TO $PG_SUPERUSER;"

      log_message "Superuser '$PG_SUPERUSER' created successfully."
    else
      log_message "Error: PostgreSQL service is not running. Please start the service and run the script again."
    fi
  else
    log_message "Superuser '$PG_SUPERUSER' already exists. Skipping creation."
  fi
}

# Main script

# Check for PostgreSQL
if check_for_postgres; then
  # PostgreSQL is installed, create database
  create_database
  start_postgres_service
else
  # PostgreSQL is not installed, install it
  install_postgres
  start_postgres_service
fi

# Check for PostgreSQL again
if check_for_postgres; then
  # PostgreSQL is now installed, create superuser if none exists
  create_postgres_superuser
else
  log_message "Error: Unable to install PostgreSQL. Please check the setup_log.txt file for details."
fi

export PGHOST="$PG_HOST"
export PGPORT="$PG_PORT"
export PGDATABASE="$PG_DATABASE"
export PGSUPERUSER="$PG_SUPERUSER"
export PGPASSWORD="$PG_PASSWORD"
# Connection String
export CONN_STR="user=$PG_SUPERUSER dbname=$PG_DATABASE host=$PG_HOST port=$PG_PORT password=$PG_PASSWORD sslmode=disable"
echo "PGHOST is set to: $PG_HOST"
echo $CONN_STR
echo "Setup complete."
