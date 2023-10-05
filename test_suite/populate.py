#! /usr/bin/env python3

import json
from test_suite import Admin

if __name__ == "__main__":
    with open("info.json", 'r') as file:
        info = json.load(file)
    admin = Admin(info["id"], info["pass"])
    with open("student_data.json", 'r') as file:
        student_data = json.load(file)
    admin.deleteAllUsers()
    admin.addUsers(student_data)