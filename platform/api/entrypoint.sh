#!/bin/sh

python manage.py migrate
python manage.py createsuperuser --noinput
python manage.py collectstatic --noinput --clear

exec "$@"