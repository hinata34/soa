### Register Request
```
curl -v -X POST "http://localhost:8080/register" --header "Content-Type: application/json" --data '{"login": "hinata34", "email": "hinata34@yandex.ru", "password": "lolkek"}'
```

### Login Request
```
curl -v -X POST "http://localhost:8080/login" --header "Content-Type: application/json" --data '{"login": "hinata34", "email": "hinata34@yandex.ru", "password": "lolkek"}'
```

### Get Profile Request
Change {jwt_token} to token
```
curl -v -X GET "http://localhost:8080/profile" --header "Cookie: jwt=={jwt_token}"
```

### Update Profile Request
Change {jwt_token} to token
```
curl -v -X POST "http://localhost:8080/profile" --header "Content-Type: application/json" --header "Cookie: jwt={jwt_token}" --data '{"name": "Ivan", "surname": "Kochkarev", "mobile_number": "+79161123"}'
```