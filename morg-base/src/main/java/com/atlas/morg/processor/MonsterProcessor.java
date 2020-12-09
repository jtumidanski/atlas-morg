package com.atlas.morg.processor;

import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.model.Monster;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static Monster createMonster(int worldId, int channelId, int mapId, int monsterId) {
      Monster monster = MonsterRegistry.getInstance().createMonster(worldId, channelId, mapId, monsterId);
      // TODO send monster spawn message.
      return monster;
   }
}
