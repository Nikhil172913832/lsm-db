1. Every successful Put() is durable in the WAL before it returns.
2. Every successful Put() exists in exactly one memtable.
3. A WAL and its memtable always correspond.
4. A memtable being flushed is immutable.
5. Readers never observe partially rotated state.
6. Recovery reconstructs the exact state before crash.