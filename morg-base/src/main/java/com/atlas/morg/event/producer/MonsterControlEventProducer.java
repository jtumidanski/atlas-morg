package com.atlas.morg.event.producer;

import com.atlas.kafka.KafkaProducerFactory;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterControlEvent;
import com.atlas.morg.rest.event.MonsterControlEventType;
import org.apache.kafka.clients.producer.Producer;
import org.apache.kafka.clients.producer.ProducerRecord;

public class MonsterControlEventProducer {
   private static final Object lock = new Object();

   private static volatile MonsterControlEventProducer instance;

   private final Producer<Long, MonsterControlEvent> producer;

   public static MonsterControlEventProducer getInstance() {
      MonsterControlEventProducer result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new MonsterControlEventProducer();
               instance = result;
            }
         }
      }
      return result;
   }

   private MonsterControlEventProducer() {
      producer = KafkaProducerFactory.createProducer("Monster Registry", System.getenv("BOOTSTRAP_SERVERS"));
   }

   public void sendControl(int worldId, int channelId, int characterId, int uniqueId) {
      String topic = System.getenv(EventConstants.TOPIC_CONTROL_MONSTER_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new MonsterControlEvent(worldId, channelId, characterId, uniqueId, MonsterControlEventType.START)));
   }

   public void clearControl(int worldId, int channelId, int mapId, int uniqueId) {
      String topic = System.getenv(EventConstants.TOPIC_CONTROL_MONSTER_EVENT);
      long key = produceKey(worldId, channelId);
      producer.send(new ProducerRecord<>(topic, key,
            new MonsterControlEvent(worldId, channelId, mapId, uniqueId, MonsterControlEventType.STOP)));
   }

   protected Long produceKey(int worldId, int channelId) {
      return (long) ((worldId * 1000) + channelId);
   }
}
