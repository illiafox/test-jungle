# API for user images uploading and management

## Deployment

### docker-compose

Service starts after containers are up

```shell
make start # docker-compose up -d
```

| Service     | Port  | Username   | Password   |
|-------------|-------|------------|------------|
| PostgresSQL | 5432  | dev_user   | dev_pass   |
| Minio       | 9000  | minio_user | minio_pass |
| Jaeger      | 16686 |            |            |

## Application endpoints

#### API endpoint: `http://127.0.0.0:8080/`

- #### `/upload-picture` - upload picture
- #### `/images` get user images
- #### `/login` user login

#### Metrics endpoint `http://127.0.0.0:8082/`

- #### `/metrics` - [Prometheus](https://github.com/prometheus/client_golang) metrics
- #### `/debug/pprof/` - [pprof](https://pkg.go.dev/runtime/pprof)

All traces are available in Jaeger at `http://127.0.0.0:16686`

--- 

### Examples

Already added users:

| Username | Password     |
|----------|--------------|
| Illia    | somepassword |
| Michael  | 12341234     |

#### 1. Login
```bash
curl --location 'http://0.0.0.0:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"Illia",
    "password":"somepassword"
}'
```
Output:
```json
{"access_token":"token"}
```
 To place your token:
```bash
JWT_ACCESS_TOKEN=$(curl --location 'http://0.0.0.0:8080/login' \
--header 'Content-Type: application/json' \
--data '{
    "username":"Illia",
    "password":"somepassword"
}' | jq -r .access_token)
```


#### 2. Upload image

```bash
curl --location '0.0.0.0:8080/upload-picture' \
--header "Authorization: Bearer ${JWT_ACCESS_TOKEN}" \
--form 'image=@"./testdata/picture.png"'
```
Output:
```json
{"url":"http://0.0.0.0:9000/images/0e0e93ed-e209-4d31-9e79-53d4e7ed3983picture.png"}
```
You can visit this endpoint to view your picture

#### 3. Get pictures list
```bash
curl --location '0.0.0.0:8080/images' \
--header "Authorization: Bearer ${JWT_ACCESS_TOKEN}"
```
Output
```json
{
    "images": [
        {
            "name": "picture.png",
            "content_type": "image/png",
            "url": "http://0.0.0.0:9000/images/0e0e93ed-e209-4d31-9e79-53d4e7ed3983picture.png",
            "size": 937406,
            "created_at": "2023-03-21T20:16:28.253334Z"
        }
    ]
}
```