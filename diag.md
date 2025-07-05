```mermaid
graph TD
  nginx_deployment["nginx<br>(Deployment)"]
  nginx_svc_service["nginx-svc<br>(Service)"]
  nginx_svc_service --> nginx_deployment
```