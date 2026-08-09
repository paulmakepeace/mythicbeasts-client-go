# Maintaining this fork

This file exists on `integration` only. `main` mirrors upstream exactly and
carries nothing the upstream repository does not.

## Layout

- `main` is upstream's `main`, byte for byte. Nothing is committed to it.
- One topic branch per upstream PR, **branched off `main`**, holding that
  change and nothing else. A branch based on `integration` drags every other
  open change into its diff and cannot be a PR.
- `integration` is `main` plus every topic branch merged `--no-ff`, plus one
  commit carrying this file and the README fork note. It is the default branch
  and the one consumers pin to.

## Order of work

Branch, push, open the PR, and only then touch `integration`.

Integration is rebuilt from settled branches, so merging one that is still in
flight buys nothing and costs the merge, the conflict resolution and a test
run again for every revision. Worse, once `integration` is pushed, a later
correction to the branch cannot be folded into the merge and has to go out as
a forward commit on published history.

Rebuild `integration` only when a topic branch has actually changed, not on
every edit.

## Rebuilding integration

Re-merge under the real `integration` name, because merge messages are public
history:

1. Reset `integration` to `main`.
2. Merge each topic branch `--no-ff`, in the order they were opened.
3. Re-apply the fork docs commit: this file and the README note, with the note
   listing every open PR and naming the current fork tag.
4. Verify by tree-diff against the previous tip. An empty diff means the
   rebuild reproduced the same content on new history, which is the point of
   doing it.

Two README conflicts are expected and are not a sign of anything wrong: the
dormancy section against the deletion section's renamed heading, and the
status-values sections against the dormancy section, both of which add text at
the same anchor. Keep both sides; general semantics precede the flows that
rely on them.

## Tagging and consumers

Tags are lightweight, matching upstream's style: `v0.X.Y-fork.N`. Bump when
the tagged tree would differ, including when only documentation differs, since
godoc is how a caller reads it.

The README's `replace` directive names the tag, so the bump commit has to be
inside the tag it names, or the tagged tree tells a consumer to use the
previous one.

Consumers pin the tag rather than a path, so each bump is followed by
`go mod edit -replace` plus `go mod tidy` in `dormouse` and `mb-probe`. Do not
bump a consumer before the tag is pushed: it cannot resolve, and it leaves a
working consumer unable to build for the sake of getting ahead.

## Conventions

- Fetch upstream with `--no-tags`. Upstream carries a stray `list` tag that
  tag-following re-imports.
- Pushes and `gh` operations are gated: the operator vets outward-facing text
  (commit messages, PR bodies, README) before it leaves the machine.
- Commit as the repository's configured identity. Never override it.
