package com.atlas.morg.builder;

import com.atlas.morg.model.Monster;

public class MonsterBuilder {
   private final Integer worldId;

   private final Integer channelId;

   private final Integer mapId;

   private final Integer uniqueId;

   private final Integer monsterId;

   private Integer x;

   private Integer y;

   private Integer fh;

   private Integer stance;

   private Integer team;

   private Integer controlCharacterId;

   public MonsterBuilder(Integer worldId, Integer channelId, Integer mapId, Integer uniqueId, Integer monsterId, Integer x,
                         Integer y, Integer fh, Integer stance, Integer team) {
      this.worldId = worldId;
      this.channelId = channelId;
      this.mapId = mapId;
      this.uniqueId = uniqueId;
      this.monsterId = monsterId;
      this.x = x;
      this.y = y;
      this.fh = fh;
      this.stance = stance;
      this.team = team;
   }

   public MonsterBuilder(Monster other) {
      this.worldId = other.worldId();
      this.channelId = other.channelId();
      this.mapId = other.mapId();
      this.uniqueId = other.uniqueId();
      this.monsterId = other.monsterId();
      this.x = other.x();
      this.y = other.y();
      this.fh = other.fh();
      this.stance = other.stance();
      this.team = other.team();
   }

   public MonsterBuilder setX(Integer x) {
      this.x = x;
      return this;
   }

   public MonsterBuilder setY(Integer y) {
      this.y = y;
      return this;
   }

   public MonsterBuilder setFh(Integer fh) {
      this.fh = fh;
      return this;
   }

   public MonsterBuilder setStance(Integer stance) {
      this.stance = stance;
      return this;
   }

   public MonsterBuilder setTeam(Integer team) {
      this.team = team;
      return this;
   }

   public MonsterBuilder setControlCharacterId(Integer controlCharacterId) {
      this.controlCharacterId = controlCharacterId;
      return this;
   }

   public Monster build() {
      return new Monster(worldId, channelId, mapId, uniqueId, monsterId, controlCharacterId, x, y, fh, stance, team);
   }
}
