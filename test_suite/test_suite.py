#! /usr/bin/env python3

import requests, json, hashlib, base64
from datetime import date
from OpenSSL import crypto
from Crypto import Random
from Crypto.Cipher import AES
from time import strftime, localtime
from colorama import Fore, Back, Style

status_color = {
    '+': Fore.GREEN,
    '-': Fore.RED,
    '*': Fore.YELLOW,
    ':': Fore.CYAN,
    ' ': Fore.WHITE
}

def display(status, data, start='', end='\n'):
    print(f"{start}{status_color[status]}[{status}] {Fore.BLUE}[{date.today()} {strftime('%H:%M:%S', localtime())}] {status_color[status]}{Style.BRIGHT}{data}{Fore.RESET}{Style.RESET_ALL}", end=end)

class Admin:
    adminLogin_path = "/session/admin/login"
    addNewUser_path = "/admin/user/new"
    deleteUser_path = "/admin/user/delete"
    deleteAllUsers_path = "/admin/user/deleteallusers"
    newUserFields = ["roll", "name", "email", "gender", "passHash"]
    deleteUserFields = ["roll", "name", "email", "gender"]
    def __init__(self, id, password, host="127.0.0.1", port=8080):
        self.session = requests.session()
        self.host = host
        self.port = port
        self.url = f"http://{host}:{port}"
        self.id = id
        self.password = password
        self.headers = {}
        self.logIn()
    def logIn(self):
        data = self.session.post(f"{self.url}{Admin.adminLogin_path}", data=json.dumps({"id": self.id, "pass": self.password}))
        response = json.loads(data.text)
        try:
            if response["message"] == "Admin logged in successfully !!":
                self.cookie = data.headers["Set-Cookie"]
                self.headers["Cookie"] = self.cookie
        except:
            if "error" in response.keys():
                display('-', f"Error in Admin LogIn: {Back.YELLOW}{response['error']}{Back.RESET}")
    def addUsers(self, users):
        users = [user for user in users if Admin.checkNewUserFormat(user)]
        data = self.session.post(f"{self.url}{Admin.addNewUser_path}", data=json.dumps({"newuser": users}), headers=self.headers)
        response = json.loads(data.text)
        try:
            if response["message"] == "User created successfully.":
                display('+', f"{Back.MAGENTA}{len(users)}{Back.RESET} Users Added")
                return len(users)
        except:
            if "error" in response.keys():
                display('-', f"Error in Admin User Add: {Back.YELLOW}{response['error']}{Back.RESET}")
                return -1
    def deleteAllUsers(self):
        self.session.get(f"{self.url}{Admin.deleteAllUsers_path}", headers=self.headers)
    @staticmethod
    def checkNewUserFormat(user):
        if type(user) == dict:
            if len(user.keys()) != len(Admin.newUserFields):
                return False
            for key in user.keys():
                if key not in Admin.newUserFields:
                    return False
        return True

if __name__ == "__main__":
    pass