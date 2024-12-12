package msgbus

// NATS implemented subjects generator
// command topic: command.<proto.package.name>.<proto.message.name>[event_md.order_key]
// event topic: event.<proto.package.name>.<proto.message.name>[event_md.order_key]
// reply topic is generated on the spot for commands only with the same order
