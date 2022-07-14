## Run app

```
./start.sh
```

It deploys 2 docker containers: app at 8080 and mysql at 5432

## Test

```
./test.sh
```

It runs both unit test and integration test.
Mysql took really long to init in my machine somehow, so i add retry inside the code
The test report is located inside report folder

## Project structure

### Services
- Handle business logic
- Handle validation
- Define interface for storage layer to implement, it accepts the fact that underlying storage is a sql dbms that support transaction

### Transport
- Handle requests from outside world (we only have HTTP for now)
- Convert serialized obj into model and call service layer

### Storage
- Implements the interface exposed inside service layer 

### Common
- Acts as a wrapper for third party lib, for now we only have Mysql

### Sqlc
- Let us write raw query and use tool to generate golang code for those queries

## TODO

- [ ] Proper authentication
- [ ] Add interceptor to recover from panic.
- [ ] The app needs more logging, and a place to collect and view logs for debug (GCP log console or maybe deploy Kibana stack).
- [ ] I'm skipping WagerList integration tests because it is hard to make the result item deterministic without mocking.
- [ ] More integration test for storage component.