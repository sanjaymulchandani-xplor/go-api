
```markdown
# go-api

## Running the API

To start the server, navigate to the `employee_api` directory and run the main Go file:

```bash
cd employee_api/
go run main.go
```

## Using the API

The server listens on port 8080. You can access the employees endpoint at:

```
http://localhost:8080/employees
```

### Making a POST Request

To add a new employee, you can make a POST request to the `/employees` endpoint. Open another terminal and use the following `curl` command:

```bash
curl -X POST http://localhost:8080/employees \
    -H 'Content-Type: application/json' \
    -d '{"id": "1", "name": "John Doe", "role": "Software Engineer"}'
```

#### Breakdown of the curl Command

- `-X POST`: Specifies that this is a POST request.
- `http://localhost:8080/employees`: The URL where your Go server is running and listening for requests to the `/employees` endpoint.
- `-H 'Content-Type: application/json'`: Sets the header to indicate that the body of the request is JSON.
- `-d '{"id": "1", "name": "John Doe", "role": "Software Engineer"}'`: The data being sent in the request. This JSON object represents the new employee you want to add.

After running this command, the server should process the request, add the new employee to the in-memory data store, and return the created employee data as a JSON response. If successful, you might see a response like this:

```json
{
    "id": "1",
    "name": "John Doe",
    "role": "Software Engineer"
}
```

This confirms that the employee has been added. You can further verify by sending a GET request to `/employees/1` or `/employees` to see all employees, including the one you just added.

