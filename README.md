# SMYLET

Conveniently secure **hostel accommodation** for Nigerian schools, with just a few clicks.


## **Development**

### Installation:

1. Clone the repository.

   ```bash
   git clone
    ```

2. Build your docker image.

   ```bash
   docker-compose up --build
   ```
   
3. Go to your browser and visit `localhost:3000/` to view the application.

## **Branching Strategy**:

1. **Main Branch**:
   - Represents the stable, deployable, production-ready version of the application.
   - No direct pushes to maintain stability. Most merges into `main` comes from reviewed feature branches.

2. **Feature Branching**:
   - Each new feature or bugfix gets its own branch created from `main`.
   - Keeps ongoing development isolated from the stable `main` branch.
   - Once a feature or fix is tested and ready, it's merged back into `main` after a code review.

## **Workflow**:

1. When starting work on a new feature or fix, branch off from the latest `main`.

   ```bash
   git checkout main
   git pull
   git checkout -b feature/your-feature-name
   ```

2. Do your work on this feature branch. Commit changes as necessary.

3. Regularly sync your feature branch with `main` to get the latest changes, reducing potential merge conflicts later.

   ```bash
   git pull origin main
   ```

4. Once your feature is complete, push the branch to the remote repository.

   ```bash
   git push -u origin feature/your-feature-name
   ```
5. Ensure your code passes all automated checks and tests.
   ```bash
   docker-compose -f ./tests/integration/docker-compose.yml  up --build
   ```

6. Create a Pull Request (or Merge Request) from your feature branch to `main`.

7. Have at least one other team member review your code.

8. Ensure your Pull Request passes all automated checks and tests.

