package com.atlas.morg.event.producer;

import com.atlas.morg.EventProducerRegistry;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterKilledEvent;

public final class MonsterKilledEventProducer {
   private MonsterKilledEventProducer() {
   }

   public static void sendKilled(int worldId, int channelId, int mapId, int uniqueId, int monsterId, int x, int y, int killerId) {
      EventProducerRegistry.getInstance()
            .send(MonsterKilledEvent.class,
                  EventConstants.TOPIC_MONSTER_KILLED_EVENT,
                  worldId, channelId,
                  new MonsterKilledEvent(worldId, channelId, mapId, uniqueId, monsterId, x, y, killerId));
   }
}