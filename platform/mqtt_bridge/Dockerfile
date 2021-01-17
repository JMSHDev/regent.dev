FROM python:3.8.7-slim-buster

WORKDIR /app
COPY main.py .
COPY entrypoint.sh .
COPY requirements.txt .

RUN pip install -r requirements.txt

ENTRYPOINT ["/app/entrypoint.sh"]