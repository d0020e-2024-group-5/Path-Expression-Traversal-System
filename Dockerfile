FROM golang:1.23

# have a dir 
WORKDIR /usr/app/

# when modules is used, downlaod stuff to image
COPY ./src/go.mod ./src/go.mod
COPY ./src/go.sum ./src/go.sum
RUN go -C ./src mod download && go -C ./src mod  verify

# Copy everything from this dir to our image at /usr/app
COPY . .

# build our executable
RUN go build -C ./src -o ../main .

# run the executable /usr/app/main
CMD [ "./main" ]
