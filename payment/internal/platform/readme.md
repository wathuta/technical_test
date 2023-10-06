### Description
- Folder with platform-level logic. This directory contains all the platform-level logic that will build up the actual project, like setting up the database and storing migrations.
- All the interaction with external services are stored in this folder.

### dev commands

```
migrate -path $(MIGRATIONS_FOLDER) -database "$(DATABASE_URL)" up
```