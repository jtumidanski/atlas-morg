package com.atlas.morg;

import java.util.ArrayList;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

import com.atlas.morg.builder.MonsterBuilder;
import com.atlas.morg.model.MapKey;
import com.atlas.morg.model.Monster;

public class MonsterRegistry {
   private static final Object lock = new Object();

   private static volatile MonsterRegistry instance;

   private static final Object registryLock = new Object();

   private final Map<Integer, Monster> monsterMap;

   private final Map<MapKey, Set<Monster>> monstersInMapMap;

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

   public Monster createMonster(int worldId, int channelId, int mapId, int monsterId, int x, int y, int fh, int stance, int team) {
      synchronized (registryLock) {
         List<Integer> existingIds = new ArrayList<>(monsterMap.keySet());
         int currentUniqueId;
         do {
            if ((currentUniqueId = runningUniqueId.incrementAndGet()) >= 2000000000) {
               runningUniqueId.set(currentUniqueId = 1000000001);
            }
         } while (existingIds.contains(currentUniqueId));

         Monster monster = new MonsterBuilder(worldId, channelId, mapId, currentUniqueId, monsterId, x, y, fh, stance, team)
               .build();
         monsterMap.put(monster.uniqueId(), monster);

         MapKey mapKey = new MapKey(worldId, channelId, mapId);
         if (!monstersInMapMap.containsKey(mapKey)) {
            monstersInMapMap.put(mapKey, new HashSet<>());
         }
         monstersInMapMap.get(mapKey).add(monster);

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
         return Collections.unmodifiableSet(monstersInMapMap.get(mapKey));
      }
      return Collections.emptySet();
   }

   protected Integer getMonsterLock(int uniqueId) {
      monsterLocks.putIfAbsent(uniqueId, uniqueId);
      return monsterLocks.get(uniqueId);
   }
}
