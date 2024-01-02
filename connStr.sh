#!/bin/bash

LOG_FILE="setup_log.txt"
export PG_HOST="localhost"  # Set your PostgreSQL host here

# Function to log messages
log_message() {
  rm -f "$LOG_FILE"
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
  superuser="distilleries_of_scotland_user"
  password="ultra_secret_password"

  # Check if a superuser already exists
  if psql -h "$PG_HOST" -tAc "SELECT COUNT(*) FROM pg_user WHERE usename = '$superuser' AND usecreatedb AND usesuper;" | grep -qw 0; then
    # Check if PostgreSQL is running
    if brew services list | grep -q "postgresql"; then
      # Create superuser and set it as a member of 'distilleries_db'
      sudo -u postgres psql -h "$PG_HOST" -c "CREATE USER $superuser WITH PASSWORD '$password' SUPERUSER;"
      sudo -u postgres psql -h "$PG_HOST" -c "ALTER USER $superuser WITH CREATEDB;"
      sudo -u postgres psql -h "$PG_HOST" -c "GRANT ALL PRIVILEGES ON DATABASE distilleries_db TO $superuser;"

      log_message "Superuser '$superuser' created successfully."
    else
      log_message "Error: PostgreSQL service is not running. Please start the service and run the script again."
    fi
  else
    log_message "Superuser '$superuser' already exists. Skipping creation."
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

echo "PGHOST is set to: $PG_HOST"
echo "Setup complete."
