# Play With CK8S (PWCK8S)

## Introduction
- Play with CK8S is simar to play with docker and play with kubernetes.
- Allows customer to request a Temporary CK8S Project for a limited time. 
- The project will be deleted after the time is up.

## Features
- UI for customer to request a temporary CK8S project.
- UI for customer to manage the temporary CK8S project.
- Backend to manage the temporary CK8S project.
## Prerequisites


## Technology Stack
- Frontend: HTML, CSS, and JavaScript
- Backend: Go With http package
- Database: MySQL

## Authentication
- UI will use Oauth2 from keycloak to authenticate user.
- Keycloak will pass X509 DN to backend.


## API 

### /api/v1/user
- /api/v1/user/login


### /api/v1/project
- /api/v1/project

- assuming the user is authenticated and the user DN is passed via JWT from keycloak.
POST 
- Create a new project for the user 
- The project name will be ${USER_SID}-${RANDOM_NUMBERS}
- The project will be deleted after the time is up.
- The project will be created in the backend and the backend will return the project ID to the frontend.





    - POST: Create a new project
        - Request Body:
        ```json
        {
            "name": "project name",
            "description": "project description",
            "duration": "project duration",
            "user": "user dn"
        }
        ```
    - GET: Get all projects

- /api/v1/project/{id}

