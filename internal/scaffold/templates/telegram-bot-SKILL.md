---
name: telegram-bot
description: >
  Telegram bot integration for receiving and replying to messages. Load this skill when you
  receive a `[go-clappie] Telegram` prefixed message or when the user asks about Telegram.
---

# Telegram Bot

Bidirectional Telegram messaging via the telegram-bridge. Messages arrive as notifications in
`notifications/dirty/`, and you reply by writing files to `notifications/outbox/`.

**All paths are relative to the project root** (the directory where Claude was opened).

## Receiving Messages

When a Telegram message arrives, you'll see a notification like:

```
[go-clappie] Telegram [chat:124160036] → Bruno (@brunojuliao): What time is it?
```

The full message (with metadata) is also saved in `notifications/dirty/telegram-<id>.txt`,
relative to the project root.

## Replying to Messages

To reply, write a `.txt` file to the `notifications/outbox/` directory in the project root:

```
Your reply text here. Plain text only — no Markdown.

---
[meta]
chat_id: 124160036
```

The bridge picks up the file, sends it via the Telegram API, and deletes it.

### Example

User sends: "What time is it?"

Write the file relative to the project root — do NOT use absolute paths or `~/`:

```bash
cat > notifications/outbox/reply-$(date +%s).txt << 'EOF'
It's 17:15 on March 5th.

---
[meta]
chat_id: 124160036
EOF
```

### Threading (optional)

To reply to a specific message, add `reply_to` with the Telegram message ID:

```
Got it!

---
[meta]
chat_id: 124160036
reply_to: 12345
```

## Behavior Guidelines

- **Ack fast before doing work.** The user sees chat messages, not a terminal. Send a quick reply first, then do the work, then send the result.
- **Multiple short messages, not walls of text.** Break long responses into separate outbox files.
- **Plain text only.** Telegram renders plain text — no Markdown formatting.
- **Be conversational.** This is a chat, not a report. Keep it natural.
- **Stay open for conversation.** Don't sign off unless the user says goodbye.

## File Naming

Use any unique `.txt` filename in the outbox directory. Suggestions:
- `reply-<timestamp>.txt`
- `tg-<chat_id>-<timestamp>.txt`

## Batch Messages

When 3+ messages arrive at once, you'll see a consolidated notification:

```
[go-clappie] Telegram [chat:124160036] → 5 new messages from Bruno (see notifications/dirty/)
```

Read the individual files in `notifications/dirty/` (in the project root) to see all messages,
then reply as needed.

## Setup

This skill is installed automatically by `go-clappie init`. The telegram-bridge binary must be
running alongside Claude for bidirectional messaging to work.
