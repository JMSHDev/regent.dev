from django.core.management.base import BaseCommand, CommandError
from device.models import Credentials


class Command(BaseCommand):
    help = "Create MQTT username and password."

    def add_arguments(self, parser):
        parser.add_argument("name", type=str)
        parser.add_argument("password", type=str)

    def handle(self, *args, **options):
        if Credentials.objects.filter(name=options["name"]).exists():
            raise CommandError("Provided username already exists.")
        else:
            new_credentials = Credentials.create(options["name"], options["password"])
            new_credentials.save()
        self.stdout.write(self.style.SUCCESS(f"User {options['name']} created."))
