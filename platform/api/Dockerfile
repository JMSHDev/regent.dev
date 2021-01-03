# pull official base image
FROM python:3.8.7-slim-buster

# set work directory
WORKDIR /usr/src/api

# set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

# Required for twisted
ENV DEBIAN_FRONTEND noninteractive
RUN apt update && apt install -y gcc

# install dependencies
RUN pip install --upgrade pip
COPY ./requirements.txt .
RUN pip install -r requirements.txt

# copy entrypoint.sh
COPY ./entrypoint.sh .

# copy project
COPY . .

# run entrypoint.sh
ENTRYPOINT ["/usr/src/api/entrypoint.sh"]