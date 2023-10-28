
# SMYLET

SMYLET makes securing hostel accommodation for Nigerian schools convenient, offering a seamless experience with just a few clicks.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

Ensure you have the following installed:
- [Docker](https://docs.docker.com/get-docker/)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

### Installation

1. Clone the repository:

    ```bash
    git clone <repository-url>
    ```

2. Navigate to the project directory:

    ```bash
    cd path-to-project
    ```

3. Build your Docker images and start the services:

    ```bash
    docker-compose up --build
    ```

4. The application will be running at `http://localhost:8000/`.

## Development Workflow

### Branching Strategy

- **Main Branch**: 
  - The production-ready branch.
  - Direct pushes are disabled to maintain stability. Merging into `main` is done through Pull Requests.

- **Feature Branches**: 
  - For each feature or bug fix, a new branch is created from `main`.
  - Feature branches are merged back into `main` after code review.

### Steps

1. Create a new branch from the latest `main`:

    ```bash
    git checkout main
    git pull
    git checkout -b feature/your-feature-name
    ```

2. Work on your feature and commit your changes.

3. Sync your branch with the latest changes from `main`:

    ```bash
    git pull origin main
    ```

4. Push your branch to the remote repository:

    ```bash
    git push -u origin feature/your-feature-name
    ```

5. Run automated checks and tests:

    ```bash
    docker-compose -f ./tests/integration/docker-compose.yml up --build
    ```

6. Format your code:

    ```bash
    make go-format
    ```

7. Create a Pull Request from your branch to `main`.

8. Generate new migration files if there are changes to database models:

    ```bash
    make create-migrate MIGRATION_TARGET=short-description-of-migration
    ```

9. Have your code reviewed by at least one other team member.

10. Ensure your Pull Request passes all checks and tests.

### Running the App

- **Development**:

    ```bash
    docker-compose up --build
    ```

- **Production**:

    ```bash
    docker build -t smylet-backend:latest -f docker/Dockerfile.prod .
    docker run -p 8000:8000 smylet-backend:latest
    ```

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct.