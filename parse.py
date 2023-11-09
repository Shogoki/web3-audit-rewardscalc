import json
import sys
POT_SIZE = int(sys.argv[1]) if len(sys.argv) == 2 else 30500
def get_dup_id(body):
    lines = body.split("\n")
    for i in range(1, len(lines)):
        if "Duplicate of https" in lines[-i]:
            return lines[-i][lines[-i].rfind("/")+1:].strip()
        elif "Duplicate of #" in lines[-i]:
            return lines[-i].replace("Duplicate of #", "").strip()

    return None 

with open("output.json", "r") as f:
    output = json.load(f)


open_issues = output["data"]["repository"]["openIssues"]["nodes"]
# add dups field to open_issues 
for iss in open_issues:
    iss["dups"] = []
dups = output["data"]["repository"]["duplicates"]["nodes"]
for d in dups:
    _id = get_dup_id(d["body"])
    if(_id is None):
        # check if we have a resolved escalation here
        if any([l["name"] == "Escalation Resolved" for l in d["labels"]["nodes"]]):
            print("no dup anymore")
            continue
        else:
            raise Exception(d)
    # add dup to main issue 
    for iss in open_issues:
        if iss["number"] == int(_id):
            iss["dups"].append(d)


## Now we have mains and dups sorted
# enrich issues
total_shares =  0
for iss in open_issues:
    # print("Issue ", iss["number"], "-", iss["title"]);
    sev = [l["name"] for l in iss["labels"]["nodes"] if l["name"] in ["High", "Medium"]].pop()
    sev_multiplier = 5 if sev == "High" else 1
    total_found = 1 + len(iss["dups"])
    shares = sev_multiplier * (0.9 ** (total_found -1)) / total_found
    total_shares += shares * total_found
    iss["shares"] = shares
    iss["sev"] = sev

# final results
sanity = 0
for iss in open_issues:
    print("Issue ", iss["number"],"DUPS:", len(iss["dups"]),"SEVERITY:", iss["sev"])
    print("Share p. Watson:", iss["shares"])
    share_percent = iss["shares"]/total_shares
    print("Percent of Total Shares", share_percent *100 )
    print("Payout:", POT_SIZE * share_percent, "USDC")
    sanity += POT_SIZE * share_percent * (len(iss["dups"])+1)
    print("")
print("SANITY:", sanity, "POT", POT_SIZE)
