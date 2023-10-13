# User microservice
written in GO using GRPC
## Prepare database
1. Create postgres db (docker on any other way)
2. Run this sql migration (can be found in app/storage/postgres/migrations)   ```CREATE   TABLE users (    
   id SERIAL PRIMARY KEY,   
   username VARCHAR(255) NOT NULL,   
   email VARCHAR(255) NOT NULL,   
   password VARCHAR(255) NOT NULL,   
   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP   
   );  ```
## How to start?
1. Clone repo ```git clone https://github.com/zumosik/grpc-user-auth-service-go``` after```cd your-golang-grpc-app```
2. You need to change all postgres data in configs/config.yml (if you want also change port and timeout)
3. Build the Docker image: ```docker build -t grpc-user-service```
4. Run container ```docker run -p 8081:8081 grpc-user-service```


    
