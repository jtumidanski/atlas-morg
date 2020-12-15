package com.atlas.morg.rest;

import java.util.stream.Collectors;

import com.atlas.morg.model.Monster;
import com.atlas.morg.rest.attribute.MonsterAttributes;
import com.atlas.morg.rest.builder.MonsterAttributesBuilder;

import builder.ResultObjectBuilder;

public final class ResultObjectFactory {
   private ResultObjectFactory() {
   }

   public static ResultObjectBuilder create(Monster monster) {
      return new ResultObjectBuilder(MonsterAttributes.class, monster.uniqueId())
            .setAttribute(new MonsterAttributesBuilder()
                  .setWorldId(monster.worldId())
                  .setChannelId(monster.channelId())
                  .setMapId(monster.mapId())
                  .setMonsterId(monster.monsterId())
                  .setControlCharacterId(monster.controlCharacterId())
                  .setX(monster.x())
                  .setY(monster.y())
                  .setFh(monster.fh())
                  .setStance(monster.stance())
                  .setTeam(monster.team())
                  .setHp(monster.hp())
                  .setDamageEntries(monster.damageSummary().stream()
                        .map(entry -> new com.atlas.morg.rest.attribute.DamageEntry(entry.characterId(), entry.damage()))
                        .collect(Collectors.toList())
                  )
            );
   }
}
