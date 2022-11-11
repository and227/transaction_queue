# Simple transactions service

1. build executable
```
go build -o build/transaction_service_main user_service/main.go
go build -o build/transaction_service_main transaction_service/main.go
```
2. up containers
```bash
docker-compose up --build
```
3. Go to db container
```bash
docker exec transaction_queue_db_1 bash
psql -U tx_user -h localhost -p 5432 -d tx_db
```
4. apply migrations.sql
5. create test users (defauld balance for user = 1000)
```
http://localhost:8080/users

{
    "name": "Jonn"
}
```
6. create transactions
```
http://localhost:8080/transactions
{
    "user_id": 11,
    "amount": 2000,
    "type": "withdraw"
}
```