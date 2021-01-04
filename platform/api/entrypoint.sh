#!/bin/sh

python manage.py migrate
python manage.py createsuperuser --noinput
python manage.py collectstatic --noinput --clear

daphne -b 0.0.0.0 -p 8000 api.asgi:application

exec "$@"