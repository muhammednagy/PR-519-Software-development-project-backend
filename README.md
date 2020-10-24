#Backend for Software development project course

## Description

**Server**  endpoints:

* `GET` `/api/v1/` -  Home page
* `POST` `/api/v1/login` -  Login endpoint. Requires email and password fields in json
* `POST` `/api/v1/users` -  Create user endpoint. Requires email, nickname and password fields in json
* `PUT` `/api/v1/users/ID` -  Update user endpoint. Requires email, nickname and password fields in json and Authorization header that has the JWT
* `PUT` `/api/v1/users/ID` -  Delete user endpoint. requires Authorization header that has the JWT


## Building
To build: ```make build```  
Running or testing requires the environment variables  
To run: ```make run```  
To clean the binary file: ```make clean```

