package com.atlas.morg.processor;

import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.event.producer.MonsterEventProducer;
import com.atlas.morg.model.Monster;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static Monster createMonster(int worldId, int channelId, int mapId, int monsterId, int x, int y, int fh, int stance,
                                       int team) {
      Monster monster = MonsterRegistry.getInstance().createMonster(worldId, channelId, mapId, monsterId, x, y, fh, stance, team);
      MonsterEventProducer.getInstance().sendCreated(worldId, channelId, mapId, monster.uniqueId());
      return monster;
   }
}
