# Generated by Django 3.1.4 on 2021-03-21 22:50

from django.db import migrations, models
import django.db.models.deletion


class Migration(migrations.Migration):

    initial = True

    dependencies = []

    operations = [
        migrations.CreateModel(
            name="Device",
            fields=[
                ("id", models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name="ID")),
                ("name", models.CharField(max_length=50, unique=True)),
                ("customer", models.CharField(max_length=50)),
                ("agent_status", models.CharField(default="offline", max_length=10)),
                ("program_status", models.CharField(default="down", max_length=10)),
                ("last_updated", models.DateTimeField(auto_now=True)),
            ],
        ),
        migrations.CreateModel(
            name="Telemetry",
            fields=[
                ("id", models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name="ID")),
                ("created_on", models.DateTimeField(auto_now_add=True)),
                ("state", models.JSONField()),
                (
                    "device",
                    models.ForeignKey(
                        on_delete=django.db.models.deletion.CASCADE,
                        related_name="telemetry",
                        related_query_name="telemetry",
                        to="device.device",
                    ),
                ),
            ],
        ),
        migrations.CreateModel(
            name="MqttAuth",
            fields=[
                ("id", models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name="ID")),
                ("username", models.CharField(max_length=100, unique=True)),
                ("password", models.CharField(max_length=100)),
                ("salt", models.CharField(max_length=10)),
                ("activated", models.BooleanField(default=False)),
                (
                    "device",
                    models.ForeignKey(
                        null=True,
                        on_delete=django.db.models.deletion.CASCADE,
                        related_name="auth",
                        related_query_name="auth",
                        to="device.device",
                    ),
                ),
            ],
        ),
        migrations.CreateModel(
            name="MqttAcl",
            fields=[
                ("id", models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name="ID")),
                ("allow", models.SmallIntegerField()),
                ("ipaddr", models.CharField(max_length=60, null=True)),
                ("username", models.CharField(max_length=100, null=True)),
                ("clientid", models.CharField(max_length=100, null=True)),
                ("access", models.SmallIntegerField()),
                ("topic", models.CharField(max_length=100)),
                (
                    "device",
                    models.ForeignKey(
                        null=True,
                        on_delete=django.db.models.deletion.CASCADE,
                        related_name="acl",
                        related_query_name="acl",
                        to="device.device",
                    ),
                ),
            ],
        ),
    ]
