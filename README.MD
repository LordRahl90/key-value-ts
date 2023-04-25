## Key Value Timestamp Store

This is a web API system that stores and reads timestamp sequence in a file <br />


### Start Up

* Clone the application
* Run `go mod tidy`
* Run `make start`

To <b>PUT</b> a sequence: 
```bash
curl -X PUT http://localhost:8080 -H 'Content-Type: application/json' -d '{"key": "mykey", "value": "myvalue", "timestamp" : 1673524092123456}'
```
Response [Code: 201]:
```json
{
    "message": "sequence saved successfully"
}
```

To <b>GET</b> the value:
```bash
curl -X GET http://localhost:8080 -H 'Content-Type: application/json' -d '{"key":"mykey", "timestamp": 1673524092123456}'
```
Response [Code: 200]:
```json
{
    "value": "myvalue"
}
```

### Test
* Run `make test` to run both the unit and integration tests