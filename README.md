# Day Day Movie

RESTful service that provides day movie 

### Support Login

- [x] Google login

- [x] Facebook login


### Docker
1. Build Image

    ```console
    docker build -t patricelee/daydaymovie -f ./deploy/Dockerfile .
    ```

1. Run container

    ```console
    docker run --rm --name daydaymovie --network host -d -p 8080:8080 --env-file=./cmd/movie/.env -v "$(pwd)"/data:/data -v "$(pwd)"/assets:/assets patricelee/daydaymovie
    ```