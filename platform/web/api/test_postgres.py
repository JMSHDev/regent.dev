import os
from time import sleep
import psycopg2

MAX_WAIT = 300

DBNAME = os.environ["SQL_DATABASE"]
USER = os.environ["SQL_USER"]
PASSWORD = os.environ["SQL_PASSWORD"]
HOST = os.environ["SQL_HOST"]


def postgres_test():
    try:
        conn = psycopg2.connect(
            f"dbname='{DBNAME}' user='{USER}' host='{HOST}' password='{PASSWORD}' connect_timeout=1"
        )
        conn.close()
        return True
    except Exception as exp:
        print(str(exp))
        return False


if __name__ == "__main__":
    for i in range(0, MAX_WAIT):
        print(f"Testing Postgres connection, attempt {i+1}.")
        if postgres_test():
            print("Connection successful.")
            exit(0)
        sleep(1)
    print(f"Could not connect to Postgres in {MAX_WAIT} seconds.")
    exit(1)
