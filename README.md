# Run app
```
./start.sh
```
It deploys 2 docker containers: app and mysql

# Test
```
./test.sh
```
It runs both unit test and integration test 
The test report is located inside report folder

# TODO

- [ ] Add interceptor to recover from panic
- [ ] The app needs more logging, and a place to collect and view logs for debug (GCP log console or maybe deploy Kibana stack)
- [ ] I'm skipping WagerList integration tests because it is hard to make the result item deterministic without mocking
- [ ] More integration test for storage component