from django.core.management.base import BaseCommand, CommandError
from device.models import MqttAuth, MqttAcl
from django.conf import settings


class Command(BaseCommand):
    def handle(self, *args, **options):
        name = settings.CUSTOMER_ID
        passwd = settings.CUSTOMER_PASSWORD
        if MqttAuth.objects.filter(username=name).exists() or MqttAuth.objects.filter(username=name).exists():
            raise CommandError("Provided username already exists.")
        else:
            auth = MqttAuth.create(name, passwd, True)
            aclin = MqttAcl(allow=1, username=name, access=3, topic=f"devices/in/{name}/#")
            aclout = MqttAcl(allow=1, username=name, access=3, topic=f"devices/out/{name}/#")
            auth.save()
            aclin.save()
            aclout.save()
        self.stdout.write(self.style.SUCCESS(f"User {name} created."))
