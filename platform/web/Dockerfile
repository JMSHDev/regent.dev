###########
# BUILDER #
###########

FROM node:14.15.3-buster as builder

WORKDIR /frontend
COPY ./frontend/public ./public
COPY ./frontend/src ./src
COPY ./frontend/.browserslistrc .
COPY ./frontend/.eslintrc.js .
COPY ./frontend/babel.config.js .
COPY ./frontend/package.json .
COPY ./frontend/tsconfig.json .
COPY ./frontend/vue.config.js .
RUN apt-get update -y
RUN apt-get install dos2unix -y
RUN dos2unix *

RUN npm install
RUN npm run-script build

#########
# FINAL #
#########

FROM python:3.8.7-slim-buster

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1
ENV DEBIAN_FRONTEND noninteractive

RUN mkdir -p /app
RUN mkdir -p /app/static
RUN mkdir -p /app/dist

WORKDIR /app
COPY entrypoint.sh .
COPY api/api ./api
COPY api/device ./device
COPY api/manage.py .
COPY api/requirements.txt .
RUN apt-get update -y
RUN apt-get install dos2unix -y
RUN dos2unix *

COPY --from=builder /frontend/dist ./dist

RUN apt update && \
apt install -y gcc && \
pip install --upgrade pip && \
pip install -r requirements.txt && \
apt purge -y gcc && \
apt autoremove -y

RUN addgroup --system app && adduser --system app && usermod -G app app
RUN chown -R app:app /app
USER app

ENTRYPOINT ["/app/entrypoint.sh"]
