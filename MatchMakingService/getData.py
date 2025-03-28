# generate_matches.py
import json
import random
from pymongo import MongoClient
from urllib.parse import quote_plus

# Configure MongoDB connection using environment variables
mongo_user = os.getenv("MONGO_USER", "")
mongo_pass = os.getenv("MONGO_PASS", "")
mongo_host = os.getenv("MONGO_HOST", "localhost")
mongo_port = os.getenv("MONGO_PORT", "27017")

MONGO_URI = f"mongodb://{quote_plus(mongo_user)}:{quote_plus(mongo_pass)}@{mongo_host}:{mongo_port}/"
client = MongoClient(MONGO_URI)
db = client["ddb"]

# Fetch founders and investors
founders = list(db.founders.find({}))
investors = list(db.investors.find({}))

matches = []
for founder in founders:
    for investor in investors:
        # Check a simple matching condition: industry is in investor.preferred_industries and funding stage matches.
        if (founder.get("industry") in investor.get("preferred_industries", [])) and (founder.get("funding_stage") == investor.get("preferred_funding_stage")):
            match_percentage = random.randint(70, 100)
        else:
            match_percentage = random.randint(0, 69)
        match = {
            "founder_id": str(founder["user_id"]),
            "investor_id": str(investor["user_id"]),
            "fund_required": founder.get("fund_required", 0),
            "total_invested": investor.get("total_invested", 0),
            "industry": founder.get("industry", ""),
            "funding_stage": founder.get("funding_stage", ""),
            "preferred_funding_stage": investor.get("preferred_funding_stage", ""),
            "risk_tolerance": investor.get("risk_tolerance", ""),
            "match_percentage": match_percentage
        }
        matches.append(match)

# Save the matches to a JSON file
with open("matchmaking_data.json", "r+") as f:
    data = json.load(f)
    data["matches"] = matches
    f.seek(0)
    json.dump(data, f, indent=4)
    f.truncate()
print(
    f"Generated {len(matches)} match samples and saved to matchmaking_data.json")
