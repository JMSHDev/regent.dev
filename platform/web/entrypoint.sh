#!/bin/sh

python manage.py migrate
python manage.py createsuperuser --noinput
python manage.py collectstatic --noinput --clear
python manage.py drf_create_token -r "$DJANGO_SUPERUSER_USERNAME"
python manage.py create_customer_credentials
cp -r /app/dist/* /app/static

daphne -b 0.0.0.0 -p 8000 api.asgi:application
