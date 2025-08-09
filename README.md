# Task Journal

**Task Journal** is a cross-platform personal task management app designed to help users organize, prioritize, and track their daily activities and deadlines.  
The app is available on the web at [taskjournal.online](https://taskjournal.online), as well as on mobile devices running Android and iOS.

## Features

Task Journal offers the following features:

- Creating, searching and filtering, updating and deleting daily and long-term tasks  
- Custom designed icons that help user association with tasks  
- Synched across devices and platforms  
- Tasks history with searching and filtering features

## Technology Stack

- **Web frontend:** Angular, Tailwind CSS  
- **Mobile app:** Flutter (Dart)  
- **Backend:** Golang, PostgreSQL

---

# TaskJournalBackend

**TaskJournalBackend** is the repository for the backend REST api of Task Journal application.  
The backend api is written in Golang, primarily utilizing standard library and pgx library for database connection.
The latest production version (main branch) is hosted on [taskjournal.online/api](https://taskjournal.online/api).

## Installation and requirements

### Requirements
- Go language (version in project: 1.23.5)
- PostgreSQL server

To check Go version run the following command: `go version`.

### Installation and setup

The app establishes a connection via database connection string url, with the following format: postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName.
- YourUserName represents database user name
- YourPassword represents password for the given database user
- YourHostName represents url where database is hosted (usually localhost)
- :5432 is simply the standard opened port for connection in PostgreSQL
- YourDatabaseName represents the name of the database

Database connection string should be set as an environment variable on the system.
On Linux systems, this can be set in the `.bashrc` file by appending the following line at the end of the file:
- export DATABASE_URL=postgres://YourUserName:YourPassword@YourHostName:5432/YourDatabaseName

**Note:** Make sure the user has the appropriate privilages on the database and tables or else executing SQL querys will raise an error. 

Before running the app, ensure that all required database tables are created by executing the `create_script.sql` script.

### Running the backend application

To run the backend app, use the following command: `go run main.go`. It will automatically verify and install any required dependencies if they are missing.

