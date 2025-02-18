# Project Title

## Description

This project is a full-stack application with a Go backend API and React frontend. The application demonstrates a simple integration between a Go API server and a React web client.

It is created with Cursor IDE to test the capabilities of the IDE.

## Table of Contents

- [Project Title](#project-title)
  - [Description](#description)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Running the Application](#running-the-application)
  - [Configuration](#configuration)
  - [Database Setup](#database-setup)
  - [Usage](#usage)
  - [Contributing](#contributing)
  - [License](#license)
    - [Environment Variables](#environment-variables)
  - [Final Steps](#final-steps)

## Prerequisites

Before you begin, ensure you have met the following requirements:

- You have a Mac laptop with [Homebrew](https://brew.sh/) installed.
- You have [Go](https://golang.org/doc/install) installed (version 1.20 or higher).
- You have [Node.js](https://nodejs.org/) installed (for the frontend).
- You have [PostgreSQL](https://www.postgresql.org/download/macosx/) installed.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/<USERNAME>/<REPO>.git
   cd <REPO>
   ```

2. Install backend dependencies:
   
   ```bash
   
   cd backend
   go mod tidy
   ```

3. Install frontend dependencies:
   
   ```bash

   cd ../frontend
   npm install
   ```

## Running the Application

1. **Set up the database** (see [Database Setup](#database-setup) below).
2. **Run the backend**:
   
   ```bash

   cd backend
   go run main.go
   ```

3. **Run the frontend**:
   Open a new terminal window and run:

   ```bash

   cd frontend
   npm start
   
   ```

4. Open your browser and navigate to `http://localhost:3000` to view the application.

## Configuration

Create a `.env` file in the `backend` directory and configure the following environment variables:

```.env
DATABASE_URL=postgres://<username>:<password>@localhost:5432/<dbname>?sslmode=disable
FRONTEND_URL=http://localhost:3000
JWT_SECRET=your-super-secret-key-here
JWT_EXPIRY_MINUTES=60
```

## Database Setup

1. **Install PostgreSQL**:
   If you haven't installed PostgreSQL, you can do so using Homebrew:
   
   ```bash

   brew install postgresql

   ```

2. **Start PostgreSQL**:
   
   ```bash

   brew services start postgresql
   
   ```

3. **Create a new database**:
   
   Open the PostgreSQL command line:
   
   ```bash

   psql postgres
   ```

   Create a new database:

   ```sql

   CREATE DATABASE <dbname>;

   ```

   Replace `<dbname>` with your desired database name.

4. **Run migrations**:
   
   If you have migration files, you can run them to set up your database schema. Ensure you have the necessary migration tool installed (e.g., `migrate`).

## Usage

Provide instructions on how to use the application, including any specific features or functionalities.

## Contributing

If you would like to contribute to this project, please follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -m 'Add some feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a pull request.

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Environment Variables

- `DATABASE_URL`: Connection string for your PostgreSQL database.
- `FRONTEND_URL`: URL for the frontend application.
- `JWT_SECRET`: Secret key used for signing JWT tokens.
- `JWT_EXPIRY_MINUTES`: Expiration time for JWT tokens in minutes.

## Final Steps

Feel free to customize the content to better fit your project specifics, such as the project title, description, and any additional features or instructions. If you need further modifications or additional sections, let me know!
