#!/bin/bash

LOG_FILE="setup_log.txt"
PG_HOST="localhost"
PG_PORT="5432"

# Function to log messages
log_message() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" >> "$LOG_FILE"
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
  log_message "Installing PostgreSQL..."

  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS (using Homebrew)
    brew install postgresql
  elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux (using apt-get for Debian/Ubuntu)
    sudo apt-get update
    sudo apt-get install postgresql
  else
    log_message "Error: Unsupported operating system."
    exit 1
  fi

  log_message "PostgreSQL installed successfully."
}

# Function to start PostgreSQL service
start_postgres_service() {
  log_message "Starting PostgreSQL service..."

  if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS (using Homebrew)
    brew services start postgresql
  elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
    # Linux (using systemctl for systemd-based systems)
    sudo systemctl start postgresql
    # Alternatively, use 'service postgresql start' for non-systemd systems
  else
    log_message "Error: Unsupported operating system."
    exit 1
  fi

  sleep 5  # Give some time for PostgreSQL to start
  log_message "PostgreSQL service started."
}

create_databases() {
  dbname_dev="distilleriesdb_dev"
  dbname_test="distilleriesdb_test"

  create_single_database() {
    local dbname="$1"

    # Check if the database already exists
    if psql -h "$PG_HOST" -lqt | cut -d \| -f 1 | grep -qw "$dbname"; then
      log_message "Database '$dbname' already exists."
    else
      # Create the database
      createdb -h "$PG_HOST" "$dbname"
      log_message "Database '$dbname' created successfully."
    fi
  }
  create_single_database "$dbname_dev"
  create_single_database "$dbname_test"
}

create_postgres_user() {
  local username="$1"
  local password="$2"
  local dbname="$3"

  # Check if the user already exists
  if psql -h "$PG_HOST" -tAc "SELECT 1 FROM pg_roles WHERE rolname='$username';" | grep -qw 1; then
    log_message "User '$username' already exists. Skipping creation."
  else
    # Create the user
    sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "CREATE USER $username WITH PASSWORD '$password';"
    sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "ALTER USER $username CREATEDB;"

    log_message "User '$username' created successfully."

    # Grant privileges on the database to the user
    sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "GRANT CONNECT ON DATABASE $dbname TO $username;"
    sudo -u postgres psql -h "$PG_HOST" -p "$PG_PORT" -c "GRANT USAGE, CREATE ON SCHEMA public TO $username;"
  fi
}

# Main script

# Check for PostgreSQL
if check_for_postgres; then
  # PostgreSQL is installed, create database
  create_databases
  start_postgres_service
else
  # PostgreSQL is not installed, install it
  install_postgres
  start_postgres_service
fi

# Check for PostgreSQL again
if check_for_postgres; then
  # PostgreSQL is now installed, create superuser if none exists
  create_postgres_user "distilleries_of_scotland_user_dev" "dev_password" "distilleriesdb_dev"
  create_postgres_user "distilleries_of_scotland_user_test" "test_password" "distilleriesdb_test"
else
  log_message "Error: Unable to install PostgreSQL. Please check the setup_log.txt file for details."
fi

echo "Setup complete."
