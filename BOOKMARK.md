# Development Bookmark

## Status 2 Jan 2026

### What Is

- I have a project split into 4 repos. Cheaper on AI tokens, but even slower to develop.
- In general, it tries to get too broad: too many runtimes, too many id types, too many anticipated storage backends with too many different capabilities in terms of concurrency, atomicity, ordering, etc., aimed to support too wide a range of projects, too early.

I'd advise:
- Make the repos private again.
- Put the Go lib on hold; `workflow` already has a backend.
- Evaluate it for `posture` carefully--don't make it the main project.
  - Consider starting with S3 if that's the direction you want to end up.
  - Focus on the TS version if you decide to use it
  - Keep an eye on the spec, but simplify it. Especially, stick with just random ids, safe creation, last update wins behavior. Don't reinvent Postgres.

### Vision and Aspirations

This project began as a simple idea with complicated aspirations such as:

- Support several of my small-dataset projects with a data store that:
  - understandable with UNIX cli tools ("`less` to list")
  - git friendly
  - interoperated with cli tools, Go, and TypeScript
  - could be moved to object storage (basic S3 API) without too much overhead
- Experiment using Claude Sonnet 4.5 to add and refine features.


### Discovery

I discovered that typical project dangers remain, "even" with AI. Their texture and genesis take on a slightly different flavor and speed, but their nature is familiar, and their cautions the same:

- early splitting, generalization, and optimization are a warning
- if accompanied by too much excitement, step away and shrink scope;
- put that flood of "absolutely brilliant ideas" somewhere like an `ideas/` folder instead of code...l
- and don't ask AI to comment or riff on those ideas until:
  - 72 hours have passed
  - You've spent code time in another project
  - You've taken your anti-tunnel vision meds (a walk, a dance, an hour with nature, play that guitar, dance, call your mother, hold someones hand).

### What Happened So Far

- I started with a simple idea and a monorepo and found I could move fast on a simple idea with AI.
- I quickly had a Go library, then I imagined projects I'd done before and id styles I wish I *and added them*.
- I quickly added a TypeScript version and found I could support multiple runtimes (Node, Deno, Bun) with AI help.
- I started planning to use this in two projects and design for BOTH of them at the same time, though they had different requirements ("hey, my arms can still reach around that").

- I found that each design tweak became much more expensive in time, AI Tokens, documentation, and testing.
- AI helped, but the breadth made me less attentive to testing. I trusted AI with the testing more, only to find that it did not always test what it said/thought it was testing.
- I split the project into:
  - An abstract spec
  - A Go implementation
  - A Bun, Deno, Node implementation
  - A cross-runtime tester.

That made some things go faster, but others slower.
