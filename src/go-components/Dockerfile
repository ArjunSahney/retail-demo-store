FROM public.ecr.aws/s5z3t2n9/golang:1.15-alpine AS build
ARG GOPROXY_OVERRIDE=https://proxy.golang.org
WORKDIR /src/
RUN apk add --no-cache git bash
RUN go get -u github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
COPY src/go-components-service/go.mod /src/
COPY src/go-components-service/main.go /src/
COPY src/go-components-service/handlers /src/handlers
COPY src/go-components-service/data/*.* /src/data/
COPY src/go-components-service/ . /src/

RUN echo "Setting GOPROXY to $GOPROXY_OVERRIDE"
RUN go env -w GOPROXY=$GOPROXY_OVERRIDE
RUN CGO_ENABLED=0 go build -o /bin/go-components-service
RUN apk add ca-certificates
FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/go-components-service /bin/go-components-service
COPY --from=build /src/data/*.* /bin/data/

EXPOSE 80
ENTRYPOINT ["/bin/go-components-service"]