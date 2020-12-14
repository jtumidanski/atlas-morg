package com.atlas.morg;

import java.util.ArrayList;
import java.util.Collection;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

import com.atlas.morg.builder.MonsterBuilder;
import com.atlas.morg.model.DamageSummary;
import com.atlas.morg.model.MapKey;
import com.atlas.morg.model.Monster;

public class MonsterRegistry {
   private static final Object lock = new Object();

   private static volatile MonsterRegistry instance;

   private static final Object registryLock = new Object();

   private final Map<Integer, Monster> monsterMap;

   private final Map<MapKey, Set<Integer>> monstersInMapMap;

   private final Map<Integer, Integer> monsterLocks;

   private final AtomicInteger runningUniqueId = new AtomicInteger(1000000001);

   public static MonsterRegistry getInstance() {
      MonsterRegistry result = instance;
      if (result == null) {
         synchronized (lock) {
            result = instance;
            if (result == null) {
               result = new MonsterRegistry();
               instance = result;
            }
         }
      }
      return result;
   }

   private MonsterRegistry() {
      monsterMap = new ConcurrentHashMap<>();
      monstersInMapMap = new ConcurrentHashMap<>();
      monsterLocks = new ConcurrentHashMap<>();
   }

   public Monster createMonster(int worldId, int channelId, int mapId, int monsterId, int x, int y, int fh, int stance, int team,
                                int hp) {
      synchronized (registryLock) {
         List<Integer> existingIds = new ArrayList<>(monsterMap.keySet());
         int currentUniqueId;
         do {
            if ((currentUniqueId = runningUniqueId.incrementAndGet()) >= 2000000000) {
               runningUniqueId.set(currentUniqueId = 1000000001);
            }
         } while (existingIds.contains(currentUniqueId));

         Monster monster = new MonsterBuilder(worldId, channelId, mapId, currentUniqueId, monsterId, x, y, fh, stance, team, hp)
               .build();
         monsterMap.put(monster.uniqueId(), monster);

         MapKey mapKey = new MapKey(worldId, channelId, mapId);
         if (!monstersInMapMap.containsKey(mapKey)) {
            monstersInMapMap.put(mapKey, new HashSet<>());
         }
         monstersInMapMap.get(mapKey).add(monster.uniqueId());

         return monster;
      }
   }

   public Optional<Monster> getMonster(int uniqueId) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            return Optional.of(monsterMap.get(uniqueId));
         }
         return Optional.empty();
      }
   }

   public Set<Monster> getMonstersInMap(int worldId, int channelId, int mapId) {
      MapKey mapKey = new MapKey(worldId, channelId, mapId);
      if (monstersInMapMap.containsKey(mapKey)) {
         return monstersInMapMap.get(mapKey).stream()
               .map(monsterMap::get)
               .collect(Collectors.toUnmodifiableSet());
      }
      return Collections.emptySet();
   }

   public Collection<Monster> getMonsters() {
      return Collections.unmodifiableCollection(monsterMap.values());
   }

   protected Integer getMonsterLock(int uniqueId) {
      monsterLocks.putIfAbsent(uniqueId, uniqueId);
      return monsterLocks.get(uniqueId);
   }

   public void moveMonster(int uniqueId, int endX, int endY, int stance) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            Monster monster = monsterMap.get(uniqueId);
            monsterMap.put(uniqueId, monster.move(endX, endY, stance));
         }
      }
   }

   public void controlMonster(int uniqueId, int characterId) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            Monster monster = monsterMap.get(uniqueId);
            monsterMap.put(uniqueId, monster.control(characterId));
         }
      }
   }

   public void clearControl(int uniqueId) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            Monster monster = monsterMap.get(uniqueId);
            monsterMap.put(uniqueId, monster.clearControl());
         }
      }
   }

   public Optional<DamageSummary> applyDamage(int characterId, long damage, int uniqueId) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            Monster monster = monsterMap.get(uniqueId);
            Monster updatedMonster = monster.damage(characterId, damage);
            monsterMap.put(uniqueId, updatedMonster);
            return Optional.of(new DamageSummary(characterId, uniqueId, damage, monster.hp() - updatedMonster.hp(),
                  updatedMonster.hp() == 0));
         }
         return Optional.empty();
      }
   }

   public void removeMonster(int uniqueId) {
      synchronized (getMonsterLock(uniqueId)) {
         if (monsterMap.containsKey(uniqueId)) {
            Monster monster = monsterMap.get(uniqueId);
            monstersInMapMap.get(new MapKey(monster.worldId(), monster.channelId(), monster.mapId())).remove(monster.uniqueId());
            monsterMap.remove(uniqueId);
         }
         monsterLocks.remove(uniqueId);
      }
   }
}
