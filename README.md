# fswatch

This is a small Go library that builds on top of github.com/fsnotify/fsnotify to
add some features that I almost always want when using this library:

- Recursive directory watches (and watch removals -- see
  https://github.com/fsnotify/fsnotify/issues/41).
- Ignore CHMOD-only changes; this works around a macOS issue that Spotlight
  likes to touch files again immediately after they've been touched (and I
  normally don't care about attribute changes; only additions and deletions).
- Coalesce rapid sequences of changes over a configurable threshold: this helps
  in various scenarios, such as editors writing out via temp files.

## Warning

This is a library I wrote purely for my own needs and I'm pretty sure it has
some bugs. Or maybe fsnotify has some bugs. The whole file system notification
thing has all kinds of fun flakiness, in my experience.

## To do

Another thing I often want is to be able to push down some filter function to
avoid walking certain directories (or files). But since that can be accomplished
more efficiently by making better use of the underlying APIs (i.e., writing a
better fsnotify).
