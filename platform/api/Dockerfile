FROM python:3.8.7-slim-buster

ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1
ENV DEBIAN_FRONTEND noninteractive

RUN mkdir -p /app
RUN mkdir -p /app/static

WORKDIR /app
COPY api ./api
COPY device ./device
COPY entrypoint.sh .
COPY manage.py .
COPY requirements.txt .

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
