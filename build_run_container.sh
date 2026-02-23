
docker build -t oregon-archery .
docker run -p 8080:8080 -e oregon-archery:latest
