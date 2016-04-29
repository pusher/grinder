Grinder
=======

API
---

### Authentication

Authentication is done by adding an `Authorization` header to every request. The header contains a token unique to each user, which is returned when a user is claimed.


### Unclaimed Users

    GET /claim HTTP/1.1

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8

    [
      {"id":"71bebef5-a66b-4633-811f-cf7e391e9d85","name":"Paul Bayer-Eynck"},
      {"id":"e2ca893e-ffdf-444b-b582-e2b2f61acc5a","name":"Laurie Wang"}
    ]


### Claim User

    POST /claim/:user_uuid HTTP/1.1

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=itf-8

    {
      "id": "e2ca893e-ffdf-444b-b582-e2b2f61acc5a",
      "name": "Laurie Wang",
      "token": "<user token>"
    }

  - token: user authentication token


### Get Current User

    GET /user HTTP/1.1
    Authorization: Token <user token>

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8

    {
      "id": "e2ca893e-ffdf-444b-b582-e2b2f61acc5a",
      "name": "Laurie Wang",
      "available": true
    }

  - available: user is available for the next round


### Toggle Availability

    POST /user/available HTTP/1.1
    Authorization: Token <user token>

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8

    {
      "id": "e2ca893e-ffdf-444b-b582-e2b2f61acc5a",
      "name": "Laurie Wang",
      "available": true
    }


### List Matches

    GET /user/match HTTP/1.1
    Authorization: Token <user token>

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8

    [
      {"id":"71bebef5-a66b-4633-811f-cf7e391e9d85","name":"Paul Bayer-Eynck", "match": false},
      {"id":"e2ca893e-ffdf-444b-b582-e2b2f61acc5a","name":"Laurie Wang", "match": false}
    ]

  - match: if the current user has swiped to match with user


### Match User

    POST /user/match/:user_id
    Authorization: Token <user token>

    HTTP/1.1 200 OK
    Content-Type: application/json; charset=utf-8


### Internal: Reset

Reset user availability and match state after round.

    GET /admin/reset HTTP/1.1

    HTTP/1.1 200 OK


### Internal: Match

Perform matching between available users

    GET /admin/match HTTP/1.1

    HTTP/1.1 200 OK

