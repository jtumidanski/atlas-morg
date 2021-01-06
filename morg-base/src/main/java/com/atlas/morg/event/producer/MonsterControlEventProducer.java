package com.atlas.morg.event.producer;

import com.atlas.morg.EventProducerRegistry;
import com.atlas.morg.rest.constant.EventConstants;
import com.atlas.morg.rest.event.MonsterControlEvent;
import com.atlas.morg.rest.event.MonsterControlEventType;

public final class MonsterControlEventProducer {
   private MonsterControlEventProducer() {
   }

   public static void sendControl(int worldId, int channelId, int characterId, int uniqueId) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_CONTROL_MONSTER_EVENT, uniqueId,
            new MonsterControlEvent(worldId, channelId, characterId, uniqueId, MonsterControlEventType.START));
   }

   public static void clearControl(int worldId, int channelId, int mapId, int uniqueId) {
      EventProducerRegistry.getInstance().send(EventConstants.TOPIC_CONTROL_MONSTER_EVENT, uniqueId,
            new MonsterControlEvent(worldId, channelId, mapId, uniqueId, MonsterControlEventType.STOP));
   }
}
