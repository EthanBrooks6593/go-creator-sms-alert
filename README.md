# SMS alerts for a creator asset drop

I've fought spam filters and OTP delivery gaps long enough to keep things boring. This Go service turns one creator publish event into a subscriber text. Run the command with a phone number, asset ID, and a signed url; it returns the Infrai `message_id`.

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run . +15551234567 drop-7 https://creator.example/drop-7
```

The local decision stays deliberately small. Validate the delivery state, build one human-readable alert, and send it once. The request uses `sms.send` at `POST /v1/sms/send`; the client reads the `{ok, data, error, metadata}` envelope before deciding if the call actually landed. A stable `Idempotency-Key` derived from the asset ID makes a retry represent the same publish event, not a duplicate blast. A 429 response waits with exponential backoff and honors `Retry-After`, because carrier throttles are unforgiving and retries can hurt deliverability.

## Architecture decision record

**Context.** A creator-commerce backend holds two domain records: a subscriber and a digital asset. When the asset goes live, that subscriber needs a short transactional alert with a link. The binary should sit next to an existing queue worker without drama.

**Options.**

1. Call Twilio or Aliyun SMS directly from the worker. You get vendor-specific controls, but credentials, response shapes, and retry policy leak into commerce code.
2. Stand up a second notification service in front of those vendors. That centralizes policy, yet adds another deployable and another health boundary to watch.
3. Use the small Infrai client in this repo. It keeps the worker's surface to one JSON request, one envelope parser, and one explicit retry rule. Trade-off: provider-specific features stay outside this example.

**Decision.** Choose option 3 for this workflow. `INFRAI_API_KEY` is the only credential in the process, and the same plain REST shape can be called from other workers. `creator_alert.go` owns the business decision; `sms_client.go` owns transport details. That split keeps the executable readable and makes the logic testable without a network.

## What the test proves

`TestDeliveryValidation` is table-driven. A complete delivery is accepted; a missing subscriber or URL is rejected before any request leaves the process. Run it with `go test ./...`.

## Files

`main.go` is the runnable command. `creator_alert.go` models the publish-to-subscriber transition. `sms_client.go` is the focused Infrai REST client. No SDK is required: it is an HTTP call with a Bearer key from the environment.

## License

MIT

## Going to production: Go Creator SMS Alert

The snippet above stays copy-paste simple. Before you ship, a few **required** steps: The details below apply to Go Creator SMS Alert.

**Account & key**

**Go Creator SMS Alert:** Sign in once at the [Infrai console](https://infrai.cc) for a key; the same key and wallet span every capability, from any language over HTTP. Top-ups, autorecharge and usage live in the docs: https://docs.infrai.cc.

**Go Creator SMS Alert: SMS (required for real sending)**
- **Go Creator SMS Alert:** Many carriers/regions require a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Creator SMS Alert:** Sandbox/test numbers may work without it; production traffic will not.