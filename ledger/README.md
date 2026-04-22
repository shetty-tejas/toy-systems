# Ledger

Ledger is a Kafka-inspired Commit Log.

## High-Level Design

1. Ledger gets a message from the user via the 'SegmentWriter' struct.
2. The SegmentWriter internally creates the 'segment' folder and creates a segment sub-folder padded to 16 zeros.
  First folder would be '0000000000000000'. Second would be previous+limit = '0000000000001000' if limit is 1000.
3. Each sub-folder would have two files index.idx which would store the positional offset to the message, and store.log which would store the header with an 8 byte length prefix.
4. Writing basically is implemented as:
  - Fetch the latest segment to write the message in.
  - Fetch the latest size of the segment.
  - Use the size as the offset to the current message.
  - Get the length in bytes to the message, and prefix the length as an 8 byte header to the message, and write it into the store.
  - Use this offset and encode it to an 8 byte message, and append it to index.idx.
5. Reading is implemented as:
  - Get the position to read from.
  - The position is then fetched as an ephemeral offset by multiplying it with 8 for finding the index record.
    Basically out of an index with a 1000 records, if you wanna find 500th record, 500*8 would be where you will find it due to the fixed space for each offset (8 bytes).
  - Once you find the offset (position*8 to (position+1)*8), you decode the binary to uint64.
  - This is the offset where you will find the message in the store.
  - Use this offset to get the start of the header-prefix of the particular message.
  - Header-prefix is an 8byte header, so use the data present in offset to offset+8 space.
  - Decode this header prefix and get the length of the message.
  - Finally find the data present from offset+8 to offset+8+headerSize and that will be the message.

