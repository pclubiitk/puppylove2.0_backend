import requests
import json
s = requests.Session()

adminid = input("Enter id ")
password = input("enter pass ")

url = "http://localhost:8080"

headers = {
    'Content-Type': 'application/json',
}

res = s.post(f'{url}/session/admin/login',data=json.dumps({
    "id":adminid,
    "pass":password
}))
cookie = res.headers.get("Set-Cookie").split(";")[0]
headers_with_cookie = {
    'Content-Type': 'application/json',
    'Cookie': cookie
}

method = input("to end the permit enter 1\nto publish the results enter 2\n")

if method == "1":
    res = s.get(f'{url}/admin/TogglePermit',headers=headers_with_cookie)
    print(res.json)
elif method == "2":
    res = s.get(f'{url}/admin/publish',headers=headers_with_cookie)
    print(res.json)

