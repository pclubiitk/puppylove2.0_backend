import json
import requests

adminid = input("Enter admin id ")
password = input("Enter admin pass ")

f = open("./data.json")
data = json.load(f)
headers = {
    'Content-Type': 'application/json',
}

url = "http://127.0.0.1:8080/session/admin/login"
payload = json.dumps(data)
authdata = json.dumps({
    "id":adminid,
    "pass":password
})
session = requests.Session()
response = session.post(url,headers=headers,data=authdata)
cookie = response.headers.get("Set-Cookie").split(";")[0]
headers_with_cookie = {
    'Content-Type': 'application/json',
    'Cookie': cookie
}
addUrl = "http://127.0.0.1:8080/admin/user/new"

response = session.post(addUrl,headers=headers_with_cookie,data=payload)
print(response.json())
