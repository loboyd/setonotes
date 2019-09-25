Goal: Build a messaging services which can be interacted through via an API.

`messages` table
----------------
```
CREATE TABLE messages (
    id INTEGER PRIMARY KEY,
    from_user INTEGER NOT NULL,
    to_user INTEGER NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT TIMESTAMP,
    body BYTEA NOT NULL);
```

API functionality
-----------------
- [ ] send a message to a user (probably a POST request)
  - [ ] Sending client makes a request
  - [ ] The server authenticates the client+user
  - [ ] Does the receiver heve to be verified, etc?
  - [ ] Database is updated with message (and notification?)
  - [ ] A success/error message is returned in reponse (eventually)
- [ ] receive message(s) (probably a GET request)
  - [ ] Receiving client makes a request
  - [ ] The server authenticates the client+user
  - [ ] Database is queried for all relevant messages
  - [ ] Messages are returned in response
