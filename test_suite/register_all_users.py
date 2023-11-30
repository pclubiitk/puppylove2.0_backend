#! /usr/bin/env python3

import psycopg2, json, sys
from test_suite import User

if __name__ == "__main__":
    with open("info.json", 'r') as file:
        info = json.load(file)
    password = sys.argv[1]
    if len(sys.argv) == 3:
        database = sys.argv[2]
    else:
        database = "puppylove"
    if len(sys.argv) == 4:
        host = sys.argv[3]
    else:
        host = "127.0.0.1"
    if len(sys.argv) == 5:
        port = sys.argv[4]
    else:
        port = "5432"
    connection = psycopg2.connect(database=database, host=host, port=port, user=info["id"], password=info["pass"])
    cursor = connection.cursor()
    cursor.execute("select id, auth_c from users")
    data = cursor.fetchall()
    idAuthCode = {}
    for entity in data:
        idAuthCode[entity[0]] = entity[1]
    for index, entity in enumerate(idAuthCode.items()):
        id, auth_c = entity
        user = User(id, password)
        user.loginFirst(auth_c)
        print(f"\r{index+1}/{len(idAuthCode)} ({(index+1)/len(idAuthCode):.2f})", end='')