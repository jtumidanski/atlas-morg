package com.atlas.morg.processor;

import java.util.Comparator;
import java.util.Map;
import java.util.Optional;
import java.util.stream.Collectors;

import com.atlas.morg.MonsterRegistry;
import com.atlas.morg.event.producer.MonsterControlEventProducer;
import com.atlas.morg.event.producer.MonsterEventProducer;
import com.atlas.morg.model.Monster;

public final class MonsterProcessor {
   private MonsterProcessor() {
   }

   public static Monster createMonster(int worldId, int channelId, int mapId, int monsterId, int x, int y, int fh, int stance,
                                       int team, int hp) {
      Monster monster = MonsterRegistry.getInstance().createMonster(worldId, channelId, mapId, monsterId, x, y, fh, stance, team,
            hp);

      getControllerCandidate(worldId, channelId, mapId)
            .ifPresent(controllerId -> startControl(monster, controllerId));

      MonsterEventProducer.sendCreated(worldId, channelId, mapId, monster.uniqueId(), monsterId);
      return monster;
   }

   protected static Optional<Integer> getControllerCandidate(int worldId, int channelId, int mapId) {
      Map<Integer, Integer> controllerCounts = MapProcessor.getCharacterIdsInMap(worldId, channelId, mapId)
            .collect(Collectors.toMap(id -> id, id -> 0));
      controllerCounts.putAll(
            MonsterRegistry.getInstance().getMonstersInMap(worldId, channelId, mapId).stream()
                  .filter(monster -> monster.controlCharacterId() != null)
                  .collect(Collectors.groupingBy(Monster::controlCharacterId, Collectors.summingInt(monster -> 1))));

      return controllerCounts.entrySet().stream()
            .min(Comparator.comparingInt(Map.Entry::getValue))
            .map(Map.Entry::getKey);
   }

   public static void findNextController(Monster monster) {
      getControllerCandidate(monster.worldId(), monster.channelId(), monster.mapId())
            .ifPresent(id -> startControl(monster, id));
   }

   public static void startControl(Monster monster, int characterId) {
      if (monster.controlCharacterId() != null) {
         stopControl(monster);
      }
      MonsterRegistry.getInstance().controlMonster(monster.uniqueId(), characterId);
      MonsterControlEventProducer.sendControl(monster.worldId(), monster.channelId(), characterId, monster.uniqueId());
   }

   public static void stopControl(Monster monster) {
      MonsterRegistry.getInstance().clearControl(monster.uniqueId());
      MonsterControlEventProducer.clearControl(monster.worldId(), monster.channelId(),
            monster.controlCharacterId(), monster.uniqueId());
   }

   public static void destroyAll() {
      MonsterRegistry.getInstance().getMonsters().forEach(MonsterProcessor::destroyMonster);
   }

   protected static void destroyMonster(Monster monster) {
      MonsterRegistry.getInstance().removeMonster(monster.uniqueId());
      MonsterEventProducer.sendDestroyed(monster.worldId(), monster.channelId(), monster.mapId(), monster.uniqueId(),
            monster.monsterId());
   }
}
