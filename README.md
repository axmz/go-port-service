# go-port-service

https://go-port-service.onrender.com

This project is meant to demonstrate my ability to setup a typical golang backend project. It represent a service that allows a maritime company bulk upload and query information about maritime ports.

- Hexagonal architecture
- REST API with CRUD operations. <a href="https://www.postman.com/go-port-service/go-port-service">Test with postman</a>
- Stream processing
- WebAuthn authentication
- Middleware
- Containerization
- Deployment
- Generics
- GraphQL
- Session management
- Tests: unit, integration, e2e
- Go Template engine

# TODO:
- Discoverable credentials
- Protect API with auth middleware
- CSRF middleware
- Ovservability: graphana, prometheus
- REST/GraphQL:
    - Add filtering support (e.g., filter by country).
    - Add pagination.
    - Add mutations: updatePort(id: ID!, input: PortInput!).
