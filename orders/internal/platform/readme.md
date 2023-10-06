Folder with platform-level logic. This directory contains all the platform-level logic that will build up the actual project, like setting up the database and storing migrations.

### dev commands

```
migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up
```