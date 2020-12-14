package com.atlas.morg.rest.builder;

import java.util.List;

import com.app.common.builder.RecordBuilder;
import com.atlas.morg.rest.attribute.DamageEntry;
import com.atlas.morg.rest.attribute.MonsterAttributes;

import builder.AttributeResultBuilder;

public class MonsterAttributesBuilder extends RecordBuilder<MonsterAttributes, MonsterAttributesBuilder>
      implements AttributeResultBuilder {
   private Integer worldId;

   private Integer channelId;

   private Integer mapId;

   private Integer monsterId;

   private Integer controlCharacterId;

   private Integer x;

   private Integer y;

   private Integer fh;

   private Integer stance;

   private Integer team;

   private Integer hp;

   private List<DamageEntry> damageEntries;

   @Override
   public MonsterAttributes construct() {
      return new MonsterAttributes(worldId, channelId, mapId, monsterId, controlCharacterId, x, y, fh, stance, team, hp,
            damageEntries);
   }

   @Override
   public MonsterAttributesBuilder getThis() {
      return this;
   }

   public MonsterAttributesBuilder setWorldId(Integer worldId) {
      this.worldId = worldId;
      return getThis();
   }

   public MonsterAttributesBuilder setChannelId(Integer channelId) {
      this.channelId = channelId;
      return getThis();
   }

   public MonsterAttributesBuilder setMapId(Integer mapId) {
      this.mapId = mapId;
      return getThis();
   }

   public MonsterAttributesBuilder setMonsterId(Integer monsterId) {
      this.monsterId = monsterId;
      return getThis();
   }

   public MonsterAttributesBuilder setControlCharacterId(Integer controlCharacterId) {
      this.controlCharacterId = controlCharacterId;
      return getThis();
   }

   public MonsterAttributesBuilder setX(Integer x) {
      this.x = x;
      return getThis();
   }

   public MonsterAttributesBuilder setY(Integer y) {
      this.y = y;
      return getThis();
   }

   public MonsterAttributesBuilder setFh(Integer fh) {
      this.fh = fh;
      return getThis();
   }

   public MonsterAttributesBuilder setStance(Integer stance) {
      this.stance = stance;
      return getThis();
   }

   public MonsterAttributesBuilder setTeam(Integer team) {
      this.team = team;
      return getThis();
   }

   public MonsterAttributesBuilder setHp(Integer hp) {
      this.hp = hp;
      return getThis();
   }

   public MonsterAttributesBuilder setDamageEntries(List<DamageEntry> damageEntries) {
      this.damageEntries = damageEntries;
      return getThis();
   }
}
