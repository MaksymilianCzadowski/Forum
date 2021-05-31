# Build an image
docker build -t go-forum .

# Verify that it exist
docker images

# Run the image
docker run -d -p 8080:8080 -it go-forum