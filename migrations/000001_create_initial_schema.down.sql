-- migrations/000001_create_initial_schema.down.sql

DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS task_priority;
DROP TYPE IF EXISTS task_status;