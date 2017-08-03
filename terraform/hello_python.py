def lambda_handler(event, context):
    print(event)
    print(context)
    return { "isBase64Encoded": False, "statusCode": 200, "headers": { }, "body": "hello" }
