package com.atlas.morg.event.producer;

import com.atlas.kafka.KafkaProducerFactory;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterEvent;
import com.atlas.morg.rest.event.MonsterEventType;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class MonsterEventProducer {
   private static final Object lock = new Object();

   private static volatile MonsterEventProducer instance;

   private final Producer<Long, MonsterEvent> producer;

   public static MonsterEventProducer getInstance() {
      MonsterEventProducer result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new MonsterEventProducer();
               instance = result;
            }
         }
      }
      return result;
   }

   private MonsterEventProducer() {
      producer = KafkaProducerFactory.createProducer("Monster Registry", System.getenv("BOOTSTRAP_SERVERS"));
   }

   public void sendCreated(int worldId, int channelId, int mapId, int uniqueId) {
      String topic = System.getenv(EventConstants.TOPIC_MONSTER_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new MonsterEvent(worldId, channelId, mapId, uniqueId, MonsterEventType.CREATED)));
   }

   public void sendKilled(int worldId, int channelId, int mapId, int uniqueId) {
      String topic = System.getenv(EventConstants.TOPIC_MONSTER_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new MonsterEvent(worldId, channelId, mapId, uniqueId, MonsterEventType.KILLED)));
   }

   public void sendDestroyed(int worldId, int channelId, int mapId, int uniqueId) {
      String topic = System.getenv(EventConstants.TOPIC_MONSTER_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new MonsterEvent(worldId, channelId, mapId, uniqueId, MonsterEventType.DESTROYED)));
   }

   protected Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
