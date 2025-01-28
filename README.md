# Terraform HTTP Backend

Implementation of Terraforms HTTP backend, https://developer.hashicorp.com/terraform/language/backend/http.

* Docker: `docker run -p 9944:9944 -v ./terraform-http-backend:/data sbnarra/terraform-http-backend`
  * [Compose](./docker-compose.yml)
* Helm: `helm repo add terraform-http-backend https://github.com/sbnarra/terraform-http-backend/chart && helm install terraform-http-backend/terraform-http-backend`
* Go: `git clone git@github.com:sbnarra/terraform-http-backend.git && cd terraform-http-backend && go run ./cmd/server`

## Example Usage

```hcl
terraform {
  backend "http" {
    address = "http://localhost:9944/states/<unique/path/to/deploymen>"
    lock_address = "http://localhost:9944/locks/<unique/path/to/deployment>"
    unlock_address = "http://localhost:9944/locks/<unique/path/to/deployment>"
    username = "<username>"  # Optional: If AUTH_USERNAME is defined
    password = "<password>"  # Optional: If AUTH_PASSWORD is defined
  }
}
```

## Configuration

Configuration is set using environment variables...

| Env | Desc | Default |
| - | - | - |
| DATA_DIR | Directory to store states/locks | /data |
| PORT | Listener port | 9944 |
| AUTH_USERNAME | Basic authentication username | |
| AUTH_PASSWORD | Basic authentication password | |
