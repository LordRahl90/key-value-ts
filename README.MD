## Key Value Timestamp Store

This is a web API system that stores and reads timestamp sequence in a file <br />


### Start Up

* Clone the application
* Run `go mod tidy`
* Run `make start`

To <b>PUT</b> a sequence: 
```bash
curl --location --request PUT 'localhost:8080' \
--header 'Content-Type: application/json' \
--data '{
    "key": "my-key",
    "value": "My Value",
    "timestamp": 101
}'
```
Response:
```json
{
    "message": "sequence saved successfully"
}
```

To <b>GET</> the value:
```bash
curl --location 'localhost:8080?timestamp=101&key=my-key'
```
Expected Response:
```json
{
    "value": "My Value"
}
```

### Test
* Run `make test` to run both the unit and integration tests