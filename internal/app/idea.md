1. make core element updates unblocking.
   - errors are ALL handled by the root
   - element queue updates to the root
   - root runs instanstly if not rendering
   - else it runs after current render
   - BUT that means element reads would not be locked. would that be an issue?
   - back to per-element lock held during update and read. but since updates are orchestrated by root, that should be fine?

2. move appcontext to components
