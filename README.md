# Booktastic

Booktastic is a project born with the idea of keeping track of the books I have acquired, that I have read, that I am reading, that I have to acquire and some details that I will add.

Something particular about this repository, is that it will not have code as such, but the particularities regarding the design of the system. I will be leaving the links to the related code, since I will use a microservices oriented architecture, therefore the general architecture and general sequence diagrams will be placed here and all diagrams related to the microservices.


A **microservice** represents a single self-contained functionality in a larger ecosystem, which is here represented in the following diagram (the overall system):

The tech stack I'll be using:

- **Go** - for microservices
- **Traefik** - as API Gateway
- **Postgres** - for book service storage
- **Auth0** - for authentication
- **Mongo** - for notifications data
- **RabbitMQ** - for Messages

1. #### System Architecture

![](https://github.com/Edmartt/booktastic/blob/main/assets/system%20architecture.jpg)


2. #### General flow

![](https://github.com/Edmartt/booktastic/blob/main/assets/sequence%20flow.jpg)
