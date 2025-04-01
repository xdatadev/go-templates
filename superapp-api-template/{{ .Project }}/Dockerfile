# syntax=docker/dockerfile:1.4

FROM golang:1.23-alpine as builder

# Instalação de dependências de build
RUN apk add --no-cache git
RUN apk add --no-cache aws-cli


ARG GIT_TOKEN

WORKDIR /app

# Copiando os arquivos de dependência primeiro
COPY {{.Scaffold.Resource}}/go.mod {{.Scaffold.Resource}}/go.sum ./

ENV GOPRIVATE=github.com/{{.GitHubUser}}/*
ENV GONOPROXY=github.com/{{.GitHubUser}}
RUN git config --global url."https://${GIT_TOKEN}:x-oauth-basic@github.com/{{.Scaffold.GitHubUser}}".insteadOf "https://github.com/{{.Scaffold.GitHubUser}}"

RUN --mount=type=ssh \
    go mod download


COPY assistant/. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {{.Scaffold.Resource}} 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-w -s" -o {{.Scaffold.Resource}} .

FROM alpine:3.19
COPY --from=builder /app/{{.Scaffold.Resource}} /app/{{.Scaffold.Resource}}


ENV TZ=America/Sao_Paulo
ENV GIN_MODE=release
# Adiciona certificados CA e timezone
RUN apk --no-cache add ca-certificates tzdata
RUN apk add --no-cache openssh

# Cria um usuário não-root
RUN addgroup -S app && adduser -S app -G app

# Define o diretório de trabalho
WORKDIR /app

# Copia o binário do builder
COPY --from=builder /app/{{.Scaffold.Resource}} .

# Muda a propriedade dos arquivos
RUN chown -R app:app /app
#RUN chmod 600 /app/jumpbox_kp_dev.pem
#RUN ssh -4 -i jumpbox_kp_dev.pem -o StrictHostKeyChecking=no -f -N -L 5434:superappdb.cktkjggmdpez.us-east-1.rds.amazonaws.com:5432 ubuntu@54.160.32.218

# Muda para o usuário não-root
USER app

# Expõe a porta (ajuste conforme sua necessidade)
EXPOSE 8080

# Comando para executar a aplicação
#CMD [ "sh", "-c", "ssh -4 -i jumpbox_kp_dev.pem -o StrictHostKeyChecking=no -f -N -L 5434:superappdb.cktkjggmdpez.us-east-1.rds.amazonaws.com:5432 ubuntu@54.160.32.218"]

ENTRYPOINT [ "/app/{{.Scaffold.Resource}}" ]

