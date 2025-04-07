# Pub/Sub Output Writer

## 🔧 Configuration Format (YAML or JSON)
```yaml
output:
  type: pubsub
  config:
    topic: "projects/my-project/topics/ducto-output"
    attributes_from_fields: ["source", "event_type"]
    ordering_key_field: "correlation_id"
    enable_message_ordering: true
    timeout: 5s
```

## 🧩 Supported Fields

| Field	                    | Type              | Required | Description                                                      |
|---------------------------|-------------------|----------|------------------------------------------------------------------|
| `topic`                   | string            | 	✅       | Full topic path (e.g., `projects/PROJECT_ID/topics/TOPIC_NAME`)  |
| `attributes_from_fields`  | []string          | ❌        | 	Optional list of input fields to extract and send as attributes |
| `ordering_key_field`      | string            | ❌        | 	Extract value from a field and use as `orderingKey`             |
| `enable_message_ordering` | bool              | 	❌       | 	Enables message ordering (must be enabled on the topic)         |
| `timeout`                 | duration (string) | ❌        | Timeout for publishing each message (default: `5s`)              |

## 🧪 Behavior
- ✅ Marshals input map[string]interface{} as JSON payload. 
- ✅ Extracts Pub/Sub attributes from specified fields in the payload. 
- ✅ Adds ordering key if configured and value is present. 
- ✅ Publishes message using GCP Pub/Sub Go client. 
- ✅ Honors timeout using context with deadline. 
- ✅ Collects and returns errors on failure. 
- ❌ Does not retry (yet).

## 💡 Example Transformed Output
Input JSON:
```json5
{
  "source": "api",
  "event_type": "user.created",
  "correlation_id": "abc123",
  "data": {
    "name": "Alice"
  }
}
```

Results in a Pub/Sub message:
- **Data:** Full JSON as bytes
- **Ordering Key:** abc123
- **Attributes:**
```json5
{
  "source": "api",
  "event_type": "user.created"
}
```

## 🔒 IAM Requirements
The service account must have:
- `pubsub.topics.publish` permission on the topic

## 🧪 Unit Testing Plan
- ✅ Uses pubsubtest or manual testPublisher interface 
- ✅ Validate attributes extraction 
- ✅ Validate timeout behavior 
- ✅ Validate correct message structure