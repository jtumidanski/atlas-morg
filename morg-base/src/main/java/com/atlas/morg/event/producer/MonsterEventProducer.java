package com.atlas.morg.event.producer;

import com.atlas.morg.EventProducerRegistry;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterEvent;
import com.atlas.morg.rest.event.MonsterEventType;

public final class MonsterEventProducer {
   private MonsterEventProducer() {
   }

   public static void sendCreated(int worldId, int channelId, int mapId, int uniqueId, int monsterId) {
      EventProducerRegistry.getInstance()
            .send(MonsterEvent.class, EventConstants.TOPIC_MONSTER_EVENT, worldId, channelId,
                  new MonsterEvent(worldId, channelId, mapId, uniqueId, monsterId, null, MonsterEventType.CREATED));
   }

   public static void sendDestroyed(int worldId, int channelId, int mapId, int uniqueId, int monsterId) {
      EventProducerRegistry.getInstance()
            .send(MonsterEvent.class, EventConstants.TOPIC_MONSTER_EVENT, worldId, channelId,
                  new MonsterEvent(worldId, channelId, mapId, uniqueId, monsterId, null, MonsterEventType.DESTROYED));
   }
}
