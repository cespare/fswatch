# fswatch

This is a small Go library that builds on top of github.com/fsnotify/fsnotify to
add three features that I almost always want when using this library:

- Recursive directory watches (and watch removals -- see
  https://github.com/fsnotify/fsnotify/issues/41).
- Ignore CHMOD-only changes; this works around a macOS issue that Spotlight
  likes to touch files again immediately after they've been touched (and I
  normally don't care about attribute changes; only additions and deletions).
- Coalesce rapid sequences of changes over a configurable threshold: this helps
  in various scenarios, such as editors writing out via temp files.
