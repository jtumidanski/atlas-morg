package com.atlas.morg.model;

import java.util.Comparator;
import java.util.List;
import java.util.Optional;
import java.util.stream.Collectors;

import com.atlas.morg.builder.MonsterBuilder;

public record Monster(int worldId, int channelId, int mapId, int uniqueId, int monsterId, Integer controlCharacterId, int x, int y,
                      int fh, int stance, int team, int hp, List<DamageEntry> damageEntries) {
   public Monster move(int endX, int endY, int stance) {
      return new MonsterBuilder(this).setX(endX).setY(endY).setStance(stance).build();
   }

   public Monster control(int characterId) {
      return new MonsterBuilder(this).setControlCharacterId(characterId).build();
   }

   public Monster clearControl() {
      return new MonsterBuilder(this).setControlCharacterId(null).build();
   }

   public boolean alive() {
      return hp > 0;
   }

   public Monster damage(int characterId, long damage) {
      long actualDamage = hp - Math.max(hp - damage, 0);

      return new MonsterBuilder(this)
            .setHp(Long.valueOf(hp - actualDamage).intValue())
            .addDamageEntry(new DamageEntry(characterId, actualDamage))
            .build();
   }

   public List<DamageEntry> damageSummary() {
      return damageEntries().stream()
            .collect(Collectors.groupingBy(DamageEntry::characterId, Collectors.summingLong(DamageEntry::damage)))
            .entrySet().stream()
            .map(entry -> new DamageEntry(entry.getKey(), entry.getValue()))
            .collect(Collectors.toList());
   }

   public Optional<DamageEntry> damageLeader() {
      return damageSummary().stream().max(Comparator.comparing(DamageEntry::damage));
   }
}
