# FAF.PAD16.2 -- Lab 1: Web Proxy
> **Performed by:** Vlada Magal, group FAF-203 \
> **Verified by:** asist. univ. Maxim Voloșenco

## Checkpoint 1

### Topic

The topic for the first PAD lab is a free, open source **MOOC (Massive Online Open Course) Education Platform** (alike Coursera, edX).

### Application Suitability

A microservices architecture would fit the topic because of:

* Modularity -- it is a complex application, composing of multiple componenents that have clear boundaries (for ex. User Service, Course Service, Video Service, Auth Service, Email Service); a distributed architecture allows each component to be developed (relatively) in parallel.
* Scalability -- the topic implies a large number of requests  since it is an application for *massive* online courses, so a distributed architecture allows scaling up the components as needed to load the bear more efficiently.
* Language Agnosticism -- each microservice can be written in a language more fitting for the service's purpose. In this case, an Email Service or Gateway could be written in Elixir, for example, and some other microservices in Go.
* Fault isolation/tolerance -- a failure in one microservice must not impact the whole application.

Some projects that focus on MOOC education and employ microservices are: edX, Coursera, Udemy, Khan Academy.

### Service Boundaries

Services:
* User Service -- encapsulates managing the users of the platform. Responsible for storing users, creating them, etc.
* Auth Service -- is responsible for authenticating users & generating tokens, validating tokens.
* Course Service -- responsible for managing courses.
* Gateway -- defines the API for the whole platform. Responsible for "forwarding" the request to the specific microservice and collecting the responses. Must also validate the incoming requests and ensure that protected routes are only accessible to authenticated requests.

A simple architecture diagram:

![architecture_diagram](architecture.drawio.png)
### Tech Stack and Communication Patterns

* User Service
  * Language: Go
  * Frameworks: gRPC, Gorm
  * Communication: synchronous communication via RPC
  * Database: Postgres
* Auth Service
  * Language: Go
  * Frameworks: gRPC, Gorm
  * Communication: synchronous communication via RPC
  * Datastore: Redis
* Course Service
  * Language: Go
  * Frameworks: gRPC, Gorm
  * Communication: synchronous communication via RPC
  * Database: Postgres
* Gateway
  * Language: Elixir
  * Frameworks: gRPC, Phoenix
  * Communication: synchronous communication via RPC (with the microservices) and via a REST API with the client.

### Data Management
#### Data
Each service will have its own database/cache, since the data is modular enough that sharing one data source is not necessarry.

The User Service will store user data in a Postgres DB. 

The Course Service will store course-adjacent data (courses, course categories, course chapters, users_courses cross ref table) in a Postgres DB. 

The Auth Service will store the tokens of logged in users in a Redis Cache.

#### Endpoints

HTTP Gateway endpoints: 

* User `/api/v1/users`
  * GET `/:id` -- get a user by ID
    * Response body:
    ```json
    {
      "id": "60ba873d-721b-436e-a687-06e8c6298623",
      "email": "some.email@gmail.com",
      "username": "madalv",
      "first_name":  "Vlada",
      "last_name": "Magal",
      "created_at": "2023-09-12T17:34:36.261285961Z",
      "updated_at": "2023-09-12T17:34:36.261285961Z"
    }
    ```
  * GET `/current` -- gets authenticated user
    * **Requires bearer token in Authorization header**
    * Response body: same as GET `/{id}`
  * POST `/` -- create user
    * Request body:
    ```json
    {
      "email": "some.email@gmail.com",
      "username": "madalv",
      "first_name":  "Vlada",
      "last_name": "Magal",
      "password": "verysecretpassword123",
    }
    ```
    * Response body: same as GET `/{id}`
  * DELETE `/` -- deletes authenticated user
    * **Requires bearer token in Authorization header**

* Auth `api/v1`
  * POST `/login` -- generates access token
    * Request body:
    ```json
    {
      "email": "some.email@gmail.com",
      "password": "verysecretpassword123"
    }
    ```
    * Response body:
    ```json
    {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
    }
    ```
  * POST `/logout` -- invalidates access token for current user
    * **Requires bearer token in Authorization header**

* Course `/api/v1/courses`
  * GET `/` -- gets all courses
    * Response body:
    ```json
    [{
      "id": "60ba873d-721b-436e-a687-06e8c6298623",
      "author": {
        "id": "60ba873d-721b-436e-a687-06e8c6298623",
        "username": "madalv",
        "first_name":  "Vlada",
        "last_name": "Magal",
      },
      "title": "Very Nice Course",
      "description": "This is a very nice course."
      "created_at": "2023-09-12T17:34:36.261285961Z",
      "updated_at": "2023-09-12T17:34:36.261285961Z"
    }, ...]
    ```
  * GET `/user` -- gets courses for authenticated user
    * **Requires bearer token in Authorization header**
    * Response body:
    ```json
    [{
      "id": "60ba873d-721b-436e-a687-06e8c6298623",
      "author": {
        "id": "60ba873d-721b-436e-a687-06e8c6298623",
        "username": "madalv",
        "first_name":  "Vlada",
        "last_name": "Magal",
      },
      "title": "Very Nice Course",
      "description": "This is a very nice course."
      "created_at": "2023-09-12T17:34:36.261285961Z",
      "updated_at": "2023-09-12T17:34:36.261285961Z"
    }, ...]
    ```
  * GET `/:id` -- gets course by ID
    * Response body:
    ```json
    {
      "id": "60ba873d-721b-436e-a687-06e8c6298623",
      "author": {
        "id": "60ba873d-721b-436e-a687-06e8c6298623",
        "username": "madalv",
        "first_name":  "Vlada",
        "last_name": "Magal",
      },
      "title": "Very Nice Course",
      "description": "This is a very nice course."
      "categories": [{
        "id": "70ba873d-721b-436e-a687-06e8c6298623",
        "title": "Software Engineering"
      }, ...],
      "created_at": "2023-09-12T17:34:36.261285961Z",
      "updated_at": "2023-09-12T17:34:36.261285961Z",
      "chapters": [{
        "id": "80ba873d-721b-436e-a687-06e8c6298623",
        "title": "Chapter 1 Title"
      }, ...],
    }
    ```
  * POST `/` -- creates course
    * **Requires bearer token in Authorization header**
    * Request body:
    ```json
    {
      "title": "Very Nice Course",
      "description": "This is a very nice course."
      "categories": [{
      "id": "70ba873d-721b-436e-a687-06e8c6298623",
      "title": "Software Engineering"
    }, ...],
    }
    ```
    * Response body: same as GET `/:id`
  * POST `:id/enroll` -- enrolls authenticated user into specified course 
    * **Requires bearer token in Authorization header**
  * DELETE `:id/enroll` -- deletes enrollment for the specified course
    * **Requires bearer token in Authorization header**

* Chapter `api/v1/chapters`
  * GET `/:id` -- get course chapter
    * **Requires bearer token in Authorization header**
    * Response body:
    ```json
    {
      "id": "60ba873d-721b-436e-a687-06e8c6298623",
      "course_id": "70ba873d-721b-436e-a687-06e8c6298623",
      "title": "Chapter 1 Title",
      "created_at": "2023-09-12T17:34:36.261285961Z",
      "updated_at": "2023-09-12T17:34:36.261285961Z",
      "body": "In this chapter we will learn...",
      "videos": [
        "url_1", "url_2"
      ]
    }
    ```
  * PUT `/:id` -- replaces chapter content
    * **Requires bearer token in Authorization header**
    * Request body:
    ```json
    {
      "title": "Chapter 1 Title",
      "body": "In this chapter we will learn...",
      "videos": [
        "url_1", "url_2"
      ]
    }
    ```
  * DELETE `/:id` -- deletes course chapter
    * **Requires bearer token in Authorization header**

### Deployment and Scaling

The services and databases will containerized using Docker. For container orchestration, scaling pods, and load balancing I'll set up a Kubernetes cluster.