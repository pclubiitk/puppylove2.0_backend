import json

data = {
    "newuser": []
}

# Generating user data and appending to the 'newuser' list
for i in range(210000, 219999 + 1):
    user = {
        "roll": str(i),
        "name": "User" + str(i - 210000 + 1),
        "email": str(i),
        "gender": "1",
        "passHash": "aaaa"
    }
    data["newuser"].append(user)

# Writing the data to a JSON file
with open("user_data.json", "w") as json_file:
    json.dump(data, json_file, indent=4)

print("JSON file generated successfully!")
