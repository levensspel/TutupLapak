System Design: https://excalidraw.com/#json=3EjwymYsAaerdhbhCYlOx,5egnqq3IPfBnw3s666XBFw
Development Roadmap: https://excalidraw.com/#json=3EjwymYsAaerdhbhCYlOx,5egnqq3IPfBnw3s666XBFw

Some things to nicely know about

**Independent**
- There must not be code dependency at all among services
- Service A must not know about Service B detail implementation, including tech stack and DB it uses
- Service A must only know API that is exposed by service B

**Distributed Caching**
- An VM instance is utilized especially for distributed caching
- It can be accessed by all services