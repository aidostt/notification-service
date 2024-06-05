# notification-service


### To run the container you should write
docker build -t notification-service .
docker run -p 6060:6060 --env-file .env -ti notification-service