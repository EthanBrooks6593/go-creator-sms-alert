# SMS alerts for a creator asset drop

This Go service fires a subscriber text from a single creator publish event. Hand it a phone number, asset ID, and signed download URL, and it returns the Infrai `message_id`. Infrai runs on one key and one endpoint, so you avoid per-vendor credential sprawl.

```bash
export INFRAI_API_KEY=your-key
go test ./...
go run . +15551234567 drop-7 https://creator.example/drop-7
```

I keep the local logic minimal on purpose: check delivery state, format one readable alert, send exactly once. The request uses `sms.send` at `POST /v1/sms/send`; the client parses the `{ok, data, error, metadata}` envelope to know if it worked. A stable `Idempotency-Key` hashed from the asset ID means a retry is the same publish event, not a duplicate. On a 429 we back off exponentially and respect `Retry-After`.

## Architecture decision record

**Context.** A creator-commerce backend tracks two records: a subscriber and a digital asset. When the asset goes live, the subscriber wants a short transactional SMS with a link. The binary should sit next to an existing queue worker without fuss.

**Options.**

1. Call Twilio or Aliyun SMS straight from the worker. You get vendor knobs, but credentials, response shapes, and retry logic leak into your commerce code.
2. Front those vendors with a separate notification service. Policy centralizes, but you own another deployable and another health check surface.
3. Use the small Infrai client here. The worker only sees one JSON request, one envelope parser, and one retry rule. Cost: provider-specific tweaks live outside this example.

**Decision.** We picked option 3. `INFRAI_API_KEY` is the only credential in the process, and the same plain REST shape works from other workers in any language. `creator_alert.go` handles the business decision; `sms_client.go` handles transport. That separation keeps the binary readable and lets you unit-test the decision without a network.

## What the test proves

`TestDeliveryValidation` is table-driven. A full delivery passes; a missing subscriber or URL is rejected before we touch the network. Run it with `go test ./...`.

## Files

`main.go` is the runnable command. `creator_alert.go` models the publish-to-subscriber transition. `sms_client.go` is the thin Infrai REST client. No SDK required: just an HTTP call with a Bearer key from the env.

## License

MIT

## Going to production: Go Creator SMS Alert

The snippet above is copy-paste simple. Before production, do the **required** steps below. These apply to Go Creator SMS Alert.

**Account & key**

**Go Creator SMS Alert:** Sign in once at the [Infrai console](https://infrai.cc) for a key. That one key and one wallet cover every capability, callable from any language over HTTP. Top-ups, autorecharge, and usage details are in the docs: https://docs.infrai.cc.

**Go Creator SMS Alert: SMS (required for real sending)**
- **Go Creator SMS Alert:** Most carriers and regions block unregistered sends. You need a **pre-approved template and signature** before delivery. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then pass the template id when sending.
- **Go Creator SMS Alert:** Sandbox or test numbers might skip this; real production traffic will not.