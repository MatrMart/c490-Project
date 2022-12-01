Prerequisites

To establish the database, 

Ensure MySql Server is using legacy authentication method, not sha256.
Using MySql CLI, run the following commands:

create database SocMed;
use SocMed;
source (unzipped folder)/create-tables.sql;

From MySql Workbench, create a user with the name "DBUser" with password "DBPass" that can access the SocMed database.

build and launch main.go

From a webbrowser, head to localhost:8000/login