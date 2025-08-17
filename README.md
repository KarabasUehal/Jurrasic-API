"# Jurrasic-API" 

How to start Jurrasic-API from Docker:

Select file "docker-compose.yml" and open the terminal, or just open CMD and there open application directory. Then use command:

docker-compose up --build

After this application will start and you can check if PostgresQL tables maked. Use command:

docker-compose exec postgres psql -U postgres -d mydb -c "\dt" - to see if tables exists
docker-compose exec postgres psql -U postgres -d mydb -c "SELECT * FROM Dinosaurus;" - to see the content

To check if Redis started, use command:

docker-compose exec redis redis-cli PING

Answer must be a "PONG"

To see the logs you can use:

docker-compose logs app

When you will sure that everything is working you can make requests to localhost:3000 like this:

Use GET request "/dino" - to get list of all dinosaurus (It can be used even from simple browser)
Use GET request "/dino/id" - to det some dinosaur by id (It can be used even form simple browser)

Use POST request "/dino/add" - to add new dinosaur with some json-data (from curl, Postman or similar apps)
Use PUT request "/dino/update/id" - to update some dinosaur (from curl, Postman or similar apps)
Use DELETE request "/dino/delete/id" - to delete some dinosaur (from curl, Postman or similar apps)

To stop container use command:

docker-compose down

After stopping and removing container you can find saved data in docker volume. Follow path:

\\wsl.localhost\docker-desktop\mnt\docker-desktop-disk\data\docker\volumes

Thanks for reading!