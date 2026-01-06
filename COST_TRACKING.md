# Cost Tracking Implementation for Gas Town

## Overview

Implemented a bridge between Claude Code's token tracking and Gas Town's cost system to capture per-session token usage and costs.

## Problem

- **Claude Code** tracks tokens internally (`~/.claude/stats-cache.json`) but doesn't calculate costs
- **Gas Town** has `gt costs` infrastructure but couldn't access Claude's token data
- The `gt costs record` command expected cost data from tmux but couldn't find it

## Solution

### 1. Cost Calculation Script

**Location**: `/Users/sws/gt/scripts/track-session-cost.sh`

**Features**:
- Calculates costs from token usage using model-specific pricing
- Supports Claude Opus 4.5 and Sonnet 4.5 pricing
- Tracks per-session deltas using snapshots
- Creates `session_ended` events with cost data

**Usage**:
```bash
# At session start - capture baseline
track-session-cost.sh start SESSION_NAME

# At session end - calculate delta and record cost
track-session-cost.sh end SESSION_NAME [WORK_ITEM_ID]
```

### 2. Hook Integration

**Modified Files**:
- `/Users/sws/gt/gastown/polecats/capable/.claude/settings.json`

**SessionStart Hook**:
```bash
/Users/sws/gt/scripts/track-session-cost.sh start $(tmux display-message -p '#S' 2>/dev/null || echo 'unknown')
```

**Stop Hook**:
```bash
/Users/sws/gt/scripts/track-session-cost.sh end $(tmux display-message -p '#S' 2>/dev/null || echo 'unknown') $(gt hook --json 2>/dev/null | jq -r '.bead_id // empty' || echo '')
```

### 3. Model Pricing

| Model | Input | Output | Cache Read | Cache Create |
|-------|-------|--------|-----------|-------------|
| **Opus 4.5** | $15/MTok | $75/MTok | $1.50/MTok | $18.75/MTok |
| **Sonnet 4.5** | $3/MTok | $15/MTok | $0.30/MTok | $3.75/MTok |

## Architecture

```
┌─────────────────┐
│  Claude Code    │
│  API Calls      │
└────────┬────────┘
         │ (tracks tokens)
         ↓
┌─────────────────┐
│ stats-cache.json│ <── SessionStart: snapshot
│  - inputTokens  │
│  - outputTokens │
│  - cacheRead    │
│  - cacheCreate  │
└────────┬────────┘
         │
         │ Stop: calculate delta
         ↓
┌─────────────────┐
│ Cost Calculator │ (track-session-cost.sh)
│  - Read tokens  │
│  - Apply pricing│
│  - Calculate $  │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│ .events.jsonl   │ (session_ended events)
│  - cost_usd     │
│  - tokens       │
│  - work_item    │
└────────┬────────┘
         │
         ↓
┌─────────────────┐
│   gt costs      │ (displays costs)
│  --today        │
│  --by-role      │
│  --by-rig       │
└─────────────────┘
```

## Event Format

```json
{
  "ts": "2026-01-06T19:59:00Z",
  "source": "calculate-session-cost",
  "type": "session_ended",
  "rig": "gastown",
  "role": "polecat",
  "worker": "capable",
  "session": "gt-gastown-capable",
  "work_item": "ga-ncr",
  "cost_usd": 5.42,
  "tokens": {
    "input": 1234,
    "output": 5678,
    "cache_read": 123456,
    "cache_create": 7890
  },
  "model": "claude-sonnet-4-5-20250929"
}
```

## Limitations

1. **Per-Session Granularity**: Tracks costs per Claude Code session, not per-execution node within a session
2. **Snapshot-Based**: Relies on delta tracking between session start/end, not real-time
3. **Stats Cache Dependency**: Depends on Claude Code updating `stats-cache.json` at session end

## Future Enhancements

To achieve **per-node/edge cost tracking** (as specified in ga-ncr), consider:

1. **OpenTelemetry Integration**:
   ```bash
   export CLAUDE_CODE_ENABLE_TELEMETRY=1
   export OTEL_EXPORTER_OTLP_ENDPOINT=http://collector:4317
   ```
   - Captures `claude_code.api_request` events with per-request tokens/cost
   - Real-time streaming of metrics

2. **Agent SDK Integration**:
   - Use `onMessage` callbacks to capture per-message token usage
   - Implement node-level cost attribution in Witness/Refinery coordination

3. **Task Attribution**:
   - The Stop hook already captures `work_item` from `gt hook`
   - Can be extended to track sub-tasks within a session

## Testing

```bash
# Manual test
SESSION=gt-gastown-capable
/Users/sws/gt/scripts/track-session-cost.sh start $SESSION
# ... work in session ...
/Users/sws/gt/scripts/track-session-cost.sh end $SESSION ga-test

# Check results
gt costs --today
tail /Users/sws/gt/.events.jsonl | jq 'select(.type == "session_ended")'
```

## Deployment

To deploy to all polecats:
1. Copy script to `/Users/sws/gt/scripts/track-session-cost.sh` (done)
2. Update each polecat's `.claude/settings.json` with hooks (partially done)
3. Update global template at `/Users/sws/gt/.claude/settings.json` for new polecats

## Related Beads

- **ga-ncr**: Capture token/cost metrics per node/edge (this implementation)
- **ga-6ib**: Track execution patterns (loops/fanout/queueing)
- **ga-cf1**: Implement budget enforcement (blocked on ga-ncr)
- **ga-d2b**: Build cost summary & diagram system (blocked on ga-ncr + ga-6ib)
