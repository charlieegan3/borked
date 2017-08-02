def lambda_handler(event, context):
    print(event)
    print(context)
    return {'key1': 1, 'key2': 2}
