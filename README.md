# redisgow

## Building My Own Redis in Go

### Project Overview
This project is a personal attempt to build a Redis-like key-value store using Go. Huge thanks to [Build Redis from Scratch](https://www.build-redis-from-scratch.dev/en/introduction) — without that resource, starting would have been much harder.

### Key Steps in Development
1. **Creating a TCP Listener**
   - I implemented a TCP listener that uses the same port as Redis.

2. **Implementing Serialization/Deserialization**
   - I followed Redis protocol specification for serializing and deserializing data. You can learn more here: [Redis Protocol Spec](https://redis.io/docs/latest/develop/reference/protocol-spec/).

3. **Implementing Basic Commands**
   - I started with basic Redis commands, beginning with `GET` and `SET`.

4. **Adding Persistence with AOF**
   - I implemented Append-Only File (AOF) persistence. For more details on other Redis persistence options, check out [Redis Persistence](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/).

5. **Adding TTL Support**
   - I extended the `SET` and `GET` commands to support TTL (Time-to-Live) and also added a separate `TTL` command.
   - I followed Redis handling expired keys as described here: [Redis Expiry Command](https://redis.io/docs/latest/commands/expire/).
   - Initially, I used a separate map to store expiration times. Expired keys are passively deleted when accessed.
   - **Note**: I did not follow Redis AOF TTL handling exactly. Instead of adding a `DELETE` command in the AOF, I stored a timestamp for TTLs. While this approach works, there's room for improvement to simplify the logic.

6. **Implementing Active Expiry Mechanism**
   - I added a mechanism to actively delete expired keys. A ticker triggers every second, and it randomly checks 50% of the keys in the TTL map.
   - If more than 75% of the checked keys are not expired in one cycle, the process stops. The next sweep will start after 1 second.

### Final Thoughts
This project is far from perfect, and there are many areas for improvement. However, the key takeaway for me was understanding how certain features work under the hood — like Redis key deletion mechanism, which initially seemed like magic.

---

### Enjoy Learning!
