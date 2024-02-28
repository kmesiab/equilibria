-- Create the 'equilibria' database
CREATE DATABASE IF NOT EXISTS ${DATABASE_NAME};
USE ${DATABASE_NAME};

-- Create the 'equilibria' user and grant privileges using the root user
CREATE USER '${MYSQL_ROOT_USER}'@'%' IDENTIFIED WITH 'mysql_native_password' BY '${MYSQL_ROOT_PASSWORD}';

-- Grant privileges using the root user
GRANT ALL PRIVILEGES ON ${DATABASE_NAME}.* TO '${MYSQL_ROOT_USER}'@'%';

-- Flush privileges to apply the changes
FLUSH PRIVILEGES;

CREATE USER '${DATABASE_USER}'@'%' IDENTIFIED BY '${DATABASE_PASSWORD}';
GRANT SELECT, INSERT, UPDATE, DELETE ON ${DATABASE_NAME}.* TO '${DATABASE_USER}'@'%';
FLUSH PRIVILEGES;
