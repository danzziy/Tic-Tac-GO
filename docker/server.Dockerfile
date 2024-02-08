FROM scratch
WORKDIR /app

COPY ./tic ./tic

EXPOSE 8081

ENTRYPOINT ["./tic"]
