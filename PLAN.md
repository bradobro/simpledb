# Planning Tasks for AI Document

- I've split api.go into api.go (all declarative) and tid62.go (imperative).
- I'll be moving the tid62 implementation examples elsewhere later, but not now.
- I've signifantly altered the Go api declarations, but haven't committed them, so `git diff api.go` will show my edits.
- My edits are guided by these Go principles of design:
  - If possible, keep interfaces slim
  - If if applicable, name them for the activity they accomplish, e.g. Stringer, Writer, ReadWriter
  - Try to keep names short, utilizing namespaces rathern than long names.
  - Same short names should have fairly parallel semantics
- don't modify any files, but examine api.go, ask and suggest
  - Do the changes seem to be worthwhile?
  - Are there ambiguities they create
  - Can you suggest clarifications in naming, grouping, semantics, or comments?