package com.atlas.morg.model;

import com.atlas.morg.builder.MonsterBuilder;

public record Monster(int worldId, int channelId, int mapId, int uniqueId, int monsterId, Integer controlCharacterId, int x, int y,
                      int fh, int stance, int team) {
   public Monster move(int endX, int endY, int stance) {
      return new MonsterBuilder(this).setX(endX).setY(endY).setStance(stance).build();
   }

   public Monster control(int characterId) {
      return new MonsterBuilder(this).setControlCharacterId(characterId).build();
   }

   public Monster clearControl() {
      return new MonsterBuilder(this).setControlCharacterId(null).build();
   }
}
