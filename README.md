# Snippetbox

Based on [Let's Go!](https://lets-go.alexedwards.net/)

## Prerequisites

1. Install Go and MySQL

   ```sh
   brew install go
   brew install mysql
   ```

   If you haven't already, start MySQL

   ```sh
   brew services start mysql
   ```

2. Generate TLS Certificate

   Already included in this repo, but in case you want to generate a new one.

   ```sh
   make cert
   ```

3. Prepare Database

   Login as root, use empty string as password

   ```sh
   mysql -u root -p
   Enter password:
   mysql>
   ```

   Create database

   ```sh
   CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   USE snippetbox;

   CREATE TABLE snippets (
     id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
     title VARCHAR(100) NOT NULL,
     content TEXT NOT NULL,
     created DATETIME NOT NULL,
     expires DATETIME NOT NULL
   );

   CREATE INDEX idx_snippets_created ON snippets(created);

   CREATE TABLE users (
     id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
     name VARCHAR(255) NOT NULL,
     email VARCHAR(255) NOT NULL,
     hashed_password CHAR(60) NOT NULL,
     created DATETIME NOT NULL,
     active BOOLEAN NOT NULL DEFAULT TRUE
   );

   ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
   ```

   Insert some fake data

   ```sh
   INSERT INTO snippets (title, content, created, expires) VALUES (
     'An old silent pond',
     'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
      UTC_TIMESTAMP(),
      DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
   );

   INSERT INTO snippets (title, content, created, expires) VALUES (
     'Over the wintry forest',
     'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
     UTC_TIMESTAMP(),
     DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
   );

   INSERT INTO snippets (title, content, created, expires) VALUES (
     'First autumn morning',
     'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
     UTC_TIMESTAMP(),
     DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
   );
   ```

   Create user for database

   ```sh
   CREATE USER 'web'@'localhost';
   GRANT SELECT, INSERT, UPDATE ON snippetbox.* TO 'web'@'localhost';
   ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';
   ```

   Create test database for testing purposes

   ```sh
   CREATE DATABASE test_snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```

   Create test user for test database

   ```sh
   CREATE USER 'test_web'@'localhost';
   GRANT CREATE, DROP, ALTER, INDEX, SELECT, INSERT, UPDATE ON test_snippetbox.* TO 'test_web'@'localhost';
   ALTER USER 'test_web'@'localhost' IDENTIFIED BY 'pass';
   ```

## Development

To develop, run

```
make run
```

## Testing

To test, run

```sh
make test
```

## Notes

New things that I learned:
* Output INFO log to stdout and ERROR log to stderr instead of a file, but during the launch of the app, redirect the stream to files.
* Always log on 500 error, always log
* The `Range` HTTP header can be used for partial download
* The `Allow` HTTP header used to tell clients what methods are allowed
* Template can be cached, and in Go, it can also be embedded
* Use SQL prepared statement to prevent SQL injection attacks
* Always set security headers
* Only allow specific group/user (the app) to read TLS cert files
* On network calls, think about timeouts and retries with backoff
* For page that require authentication, disable cache by setting "Cache-Control" "no-store"
