// METADATA
//
// Statusline Demo Assets
//
// Purpose
//   Provide a self-contained demo harness for the statusline binary so builds
//   can be previewed without hooking into Claude’s runtime plumbing.
//
// Files
//   session_example.jsonc  → Annotated SessionContext (human-readable, comments)
//   session_example.json   → Comment-free JSON for piping into the binary
//
// Workflow
//
//   # 1. Build (or reuse) the executable
//   go build -o statusline statusline.go
//
//   # 2. Run the demo using the comment-free JSON payload
//   cat demo/session_example.json | ./statusline
//
//   # 3. (Optional) Adjust inputs
//   #    - Edit session_example.jsonc for clarity, then strip comments:
//   python3 - <<'PY'
//   import json, re
//   from pathlib import Path
//   src = Path("demo/session_example.jsonc").read_text().splitlines()
//   clean = []
//   for line in src:
//       clean.append(re.sub(r"//.*", "", line))
//   Path("demo/session_example.json").write_text("\n".join(clean))
//   json.load(open("demo/session_example.json"))
//   PY
//
// Notes
//   - Temporal/system sections still hit live hooks; expect current machine data.
//   - Changing session_id lets you test reminder cadence deterministically.
//   - Long CWD in the sample showcases the new path-compaction logic.
//   - Keep the .json file in sync with .jsonc after edits so demos stay turnkey.
