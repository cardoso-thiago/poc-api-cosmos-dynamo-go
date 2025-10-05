import os
import uuid
import random
import time

CLOUD_PROVIDER = os.getenv("CLOUD_PROVIDER", "azure")

relationships = ["family", "friend", "neighbor"]
sub_rel = [["mother", "father", "brother"], None, None]

uuids = []
print(f"Starting data seed for {CLOUD_PROVIDER}...")
start = time.time()

if CLOUD_PROVIDER == "aws":
    import boto3
    from botocore.exceptions import ClientError
    endpoint = os.getenv("DYNAMO_ENDPOINT", "http://localhost:8000")
    region = os.getenv("DYNAMO_REGION", "us-east-1")
    table_name = os.getenv("DYNAMO_TABLE", "RelTable")

    dynamodb = boto3.resource(
        'dynamodb',
        endpoint_url=endpoint,
        region_name=region,
        aws_access_key_id='dummy',
        aws_secret_access_key='dummy',
    )
    table = dynamodb.Table(table_name)

    # Cria a tabela se não existir
    try:
        table.load()
    except ClientError as e:
        if e.response['Error']['Code'] == 'ResourceNotFoundException':
            dynamodb.create_table(
                TableName=table_name,
                KeySchema=[{'AttributeName': 'id', 'KeyType': 'HASH'}],
                AttributeDefinitions=[{'AttributeName': 'id', 'AttributeType': 'S'}],
                ProvisionedThroughput={'ReadCapacityUnits': 5, 'WriteCapacityUnits': 5}
            ).wait_until_exists()
            table = dynamodb.Table(table_name)
        else:
            raise

    for i in range(50):
        rid = str(uuid.uuid4())
        uuids.append(rid)
        idx = random.randint(0, len(relationships) - 1)
        item = {
            "id": rid,
            "relationship": relationships[idx]
        }
        if sub_rel[idx]:
            item["relationship_type"] = [{"relationship": r} for r in sub_rel[idx]]
        table.put_item(Item=item)

else:
    from azure.cosmos import CosmosClient, PartitionKey
    endpoint = os.getenv("COSMOS_ENDPOINT")
    key = os.getenv("COSMOS_KEY")
    database_name = os.getenv("COSMOS_DATABASE", "RelDB")
    container_name = os.getenv("COSMOS_CONTAINER", "RelContainer")

    client = CosmosClient(endpoint, key, connection_verify=False)
    db = client.create_database_if_not_exists(id=database_name)
    container = db.create_container_if_not_exists(
        id=container_name,
        partition_key=PartitionKey(path="/id"),
        offer_throughput=400
    )
    for i in range(50):
        rid = str(uuid.uuid4())
        uuids.append(rid)
        idx = random.randint(0, len(relationships) - 1)
        item = {
            "id": rid,
            "relationship": relationships[idx]
        }
        if sub_rel[idx]:
            item["relationship_type"] = [{"relationship": r} for r in sub_rel[idx]]
        container.upsert_item(item)

elapsed = time.time() - start
print(f"Seed finished in {elapsed:.2f}s")

try:
    with open("uuids.txt", "w") as f:
        for u in uuids:
            f.write(u + "\n")
except Exception as e:
    print(f"Erro ao salvar uuids.txt: {e}")
    exit(1)

print("\n==================== UUIDs GERADOS ====================")
for u in uuids:
    print(u)
print("================= FIM DOS UUIDs GERADOS ===============\n")