import json
import aio_pika
from shared.config import get_settings

settings = get_settings()

_connection = None
_channel = None


async def get_channel():
    global _connection, _channel
    if _connection is None or _connection.is_closed:
        _connection = await aio_pika.connect_robust(settings.MQ_CLIENT)
        _channel = await _connection.channel()
    return _channel


async def publish_message(queue_name: str, message: dict):
    channel = await get_channel()
    await channel.declare_queue(queue_name, durable=True)
    await channel.default_exchange.publish(
        aio_pika.Message(
            body=json.dumps(message).encode(),
            delivery_mode=aio_pika.DeliveryMode.PERSISTENT,
        ),
        routing_key=queue_name,
    )


async def close_connection():
    global _connection, _channel
    if _channel and not _channel.is_closed:
        await _channel.close()
    if _connection and not _connection.is_closed:
        await _connection.close()
    _connection = None
    _channel = None
