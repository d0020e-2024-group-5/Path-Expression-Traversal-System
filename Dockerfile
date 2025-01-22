FROM golang:1.23

# have a dir 
WORKDIR /usr/app

# when modules is used, downlaod stuff to image
# COPY go.mod go.sum
# RUN go mod download && go mod verify

# Copy everything from this dir to our image at /usr/app
COPY . .

# build our executable
RUN go build -o ./main src/main.go

# run the executable /usr/app/main
CMD [ "./main" ]
