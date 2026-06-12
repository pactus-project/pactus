#!/usr/bin/env python3

import sys
import plyvel

BATCH_SIZE = 10000

src_db = plyvel.DB(sys.argv[1], create_if_missing=False)
dst_db = plyvel.DB(sys.argv[2], create_if_missing=True)

count = 0
batch = dst_db.write_batch()

for key, value in src_db:
    batch.put(key, value)
    count += 1

    if count % BATCH_SIZE == 0:
        batch.write()
        batch = dst_db.write_batch()
        print(f"Copied {count:,} keys")

batch.write()

src_db.close()
dst_db.close()

print(f"Done. Copied {count:,} keys.")